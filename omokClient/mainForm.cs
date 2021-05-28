﻿using System;
using System.Collections.Generic;
using System.ComponentModel;
using System.Data;
using System.Drawing;
using System.Linq;
using System.Net.WebSockets;
using System.Text;
using System.Threading;
using System.Threading.Tasks;
using System.Windows.Forms;

namespace omokClient
{
    public partial class mainForm : Form
    {
        bool IsBackGroundProcessRunning = false;

        System.Windows.Forms.Timer dispatcherUITimer;

        ClientSimpleTcp Network = new ClientSimpleTcp();
        bool IsNetworkThreadRunning = false;
        System.Threading.Thread NetworkReadThread = null;
        PacketBufferManager PacketBuffer = new PacketBufferManager();
        System.Collections.Concurrent.ConcurrentQueue<(Int16, byte[])> ReceivePacketQueue = new();

        public mainForm()
        {
            InitializeComponent();
        }

        void BackGroundProcess(object sender, EventArgs e)
        {
            ProcessLog();

            try
            {
                if( ReceivePacketQueue.TryDequeue(out ValueTuple<Int16, byte[]> value))
                {
                    PacketProcess(value.Item1, value.Item2);
                } 
            }
            catch (Exception ex)
            {
                UILogger.Write($"ReadPacketQueueProcess. error:{ex.Message}", LOG_LEVEL.ERROR);
            }
        }

        private void ProcessLog()
        {
            // 너무 이 작업만 할 수 없으므로 일정 작업 이상을 하면 일단 패스한다.
            int logWorkCount = 0;

            while (IsBackGroundProcessRunning)
            {
                System.Threading.Thread.Sleep(1);

                string msg;

                if (UILogger.GetLog(out msg))
                {
                    ++logWorkCount;

                    if (listBoxLog.Items.Count > 512)
                    {
                        listBoxLog.Items.Clear();
                    }

                    listBoxLog.Items.Add(msg);
                    listBoxLog.SelectedIndex = listBoxLog.Items.Count - 1;
                }
                else
                {
                    break;
                }

                if (logWorkCount > 8)
                {
                    break;
                }
            }
        }

        void PacketProcess(Int16 packetID, byte[] packetBody)
        { 
            if(packetID == (Int16)PacketID.RESPONSE_BATTLE_WATCHING)
            {
                var res = new BattleWatchingResPacket();
                res.FromBytes(packetBody);

                UILogger.Write($"[PacketProcess - RESPONSE_BATTLE_WATCHING] Result: {res.Result}");
            }
            else if (packetID == (Int16)PacketID.RESPONSE_BATTLE_WATCHING_END)
            {
                var res = new BattleWatchingEndResPacket();
                res.FromBytes(packetBody);

                UILogger.Write($"[PacketProcess - RESPONSE_BATTLE_WATCHING_END] Result: {res.Result}");
            }
            else
            {
                UILogger.Write($"Unknown Packet ID: {packetID}");
            }
        }

        
        void NetworkReadProcess()
        {
            while (IsNetworkThreadRunning)
            {
                if (Network.IsConnected() == false)
                {
                    Thread.Sleep(1);
                    continue;
                }

                var recvData = Network.Receive();

                if (recvData != null)
                {
                    PacketBuffer.Write(recvData.Item2, 0, recvData.Item1);

                    while (true)
                    {
                        var data = PacketBuffer.Read();
                        if (data.Count < 1)
                        {
                            break;
                        }

                        var packetHeader = new byte[PacketDef.PACKET_HEADER_SIZE];
                        Buffer.BlockCopy(data.Array, data.Offset, packetHeader, 0, PacketDef.PACKET_HEADER_SIZE);
                        var packetID = BitConverter.ToInt16(packetHeader, 2);

                        var bodySzie = data.Count - PacketDef.PACKET_HEADER_SIZE;
                        var packetBody = new byte[bodySzie];
                        Buffer.BlockCopy(data.Array, data.Offset + PacketDef.PACKET_HEADER_SIZE, packetBody, 0, bodySzie);

                        ReceivePacketQueue.Enqueue((packetID, packetBody));
                    }
                    //UILogger.Write($"받은 데이터: {recvData.Item2}", LOG_LEVEL.INFO);
                }
                else
                {
                    Network.Close();
                    SetDisconnectd();
                    UILogger.Write("서버와 접속 종료 !!!", LOG_LEVEL.INFO);
                }
            }
        }

        public void SendPacketToSubProxyServer(byte[] packetData)
        {
            if (Network.IsConnected() == false)
            {
                UILogger.Write("서버 연결이 되어 있지 않습니다", LOG_LEVEL.ERROR);
                return;
            }

            Network.Send(packetData);
        }

        public void SetDisconnectd()
        {
            if (button6.Enabled == false)
            {
                button6.Enabled = true;
                button5.Enabled = false;
            }
                                    
            labelStatus.Text = "서버 접속이 끊어짐";
        }





        private void TestTool_Load(object sender, EventArgs e)
        {
            IsBackGroundProcessRunning = true;
            dispatcherUITimer = new System.Windows.Forms.Timer();
            dispatcherUITimer.Tick += new EventHandler(BackGroundProcess);
            dispatcherUITimer.Interval = 100;
            dispatcherUITimer.Start();


            PacketBuffer.Init((8096 * 10), PacketDef.PACKET_HEADER_SIZE, 1024);

            IsNetworkThreadRunning = true;
            NetworkReadThread = new System.Threading.Thread(this.NetworkReadProcess);
            NetworkReadThread.Start();
        }

        private void TestTool_FormClosing(object sender, FormClosingEventArgs e)
        {
            IsNetworkThreadRunning = false;
            IsBackGroundProcessRunning = false;

            Network.Close();

            NetworkReadThread.Join();
        }


       
       
        // SubProxy Server에 연결
        private void button6_Click(object sender, EventArgs e)
        {
            string address = textBoxIP.Text;

            int port = Convert.ToInt32(textBoxPort.Text);

            if (Network.Connect(address, port))
            {
                labelStatus.Text = string.Format("{0}. 서버에 접속 중", DateTime.Now);
                button6.Enabled = false;
                button5.Enabled = true;

                UILogger.Write("서버에 접속 중", LOG_LEVEL.INFO);
            }
            else
            {
                labelStatus.Text = string.Format("{0}. 서버에 접속 실패", DateTime.Now);
            }
        }

        // SubProxy Server에 끊기
        private void button5_Click(object sender, EventArgs e)
        {
            UILogger.Write("서버에 접속 끊기 요청", LOG_LEVEL.INFO);

            SetDisconnectd();
            Network.Close();
        }

        // 배틀 관전 요청
        private void button7_Click(object sender, EventArgs e)
        {
            var request = new BattleWatchingReqPacket
            {
                CompetitionCode = UInt64.Parse(textBox6.Text),
                BattleCode = UInt64.Parse(textBox5.Text),
                UserID = textBox4.Text,
                AuthCode = UInt64.Parse(textBox7.Text),
                StartFrameIndex = UInt32.Parse(textBox8.Text),
            };

            var packetData = request.ToBytes(PacketID.REQUEST_BATTLE_WATCHING);
            SendPacketToSubProxyServer(packetData);

            UILogger.Write($"[관전 요청]  CompetitionCode:{request.CompetitionCode}, BattleCode:{request.BattleCode}", LOG_LEVEL.INFO);
        }

        // 관전 종료 요청
        private void button8_Click(object sender, EventArgs e)
        {
            var packetData = NoneBodyPacket.ToBytes(PacketID.REQUEST_BATTLE_WATCHING_END);
            SendPacketToSubProxyServer(packetData);

            UILogger.Write("[관전 종료 요청]", LOG_LEVEL.INFO);
        }
    }
}
