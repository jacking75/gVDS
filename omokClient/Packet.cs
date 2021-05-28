using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace omokClient
{
    public enum PacketID
    {
        #region PubProxyServer
        PACKET_ID_BATTLE_START = 1101, 
        PACKET_ID_BATTLE_GAME_PLAY = 1102, 
        PACKET_ID_BATTLE_GAME_END = 1103,
        #endregion


        #region SubProxyServer
        CLIENT_BEGIN = 2101,

        REQUEST_BATTLE_WATCHING = 2116,
        RESPONSE_BATTLE_WATCHING = 2117,

        NOTIFY_BATTLE_PLAY_DATA = 2121,

        REQUEST_BATTLE_WATCHING_END = 2131,
        RESPONSE_BATTLE_WATCHING_END = 2132,

        CLIENT_END = 2200,
        #endregion
    }


    #region PubProxyServer
    public class BattleStartNtfPacket
    {
        public UInt64 CompetitionCode;
        public UInt64 BattleCode;

        public byte[] ToBytes(PacketID packetID)
        {
            List<byte> dataSource = new List<byte>();
            dataSource.AddRange(BitConverter.GetBytes((UInt16)packetID));
            dataSource.AddRange(BitConverter.GetBytes(CompetitionCode));
            dataSource.AddRange(BitConverter.GetBytes(BattleCode));
            return dataSource.ToArray();
        }
    }
    #endregion



    #region SubProxyServer
    class PacketDef
    {
        public const UInt16 PACKET_HEADER_SIZE = 5;
        public const int MAX_USER_ID_BYTE_LENGTH = 32;
    }

    public class NoneBodyPacket
    {
        public static byte[] ToBytes(PacketID packetID)
        {
            List<byte> dataSource = new List<byte>();
            dataSource.AddRange(BitConverter.GetBytes((Int16)PacketDef.PACKET_HEADER_SIZE));
            dataSource.AddRange(BitConverter.GetBytes((Int16)packetID));
            dataSource.AddRange(BitConverter.GetBytes((SByte)0));
            return dataSource.ToArray();
        }
    }


    public class BattleWatchingReqPacket
    {
        public UInt64 CompetitionCode;
        public UInt64 BattleCode;
        public string UserID;
        public UInt64 AuthCode;
        public UInt32 StartFrameIndex;

        public byte[] ToBytes(PacketID packetID)
        {
            var userID = new byte[PacketDef.MAX_USER_ID_BYTE_LENGTH];
            Encoding.UTF8.GetBytes(UserID).CopyTo(userID, 0);

            const int packetSize = PacketDef.PACKET_HEADER_SIZE + 8 + 8 + PacketDef.MAX_USER_ID_BYTE_LENGTH + 8 + 4;

            List<byte> dataSource = new List<byte>();
            dataSource.AddRange(BitConverter.GetBytes((Int16)packetSize));
            dataSource.AddRange(BitConverter.GetBytes((Int16)packetID));
            dataSource.AddRange(BitConverter.GetBytes((SByte)0));
            dataSource.AddRange(BitConverter.GetBytes(CompetitionCode));
            dataSource.AddRange(BitConverter.GetBytes(BattleCode));
            dataSource.AddRange(userID);
            dataSource.AddRange(BitConverter.GetBytes(AuthCode));
            dataSource.AddRange(BitConverter.GetBytes(StartFrameIndex));
            return dataSource.ToArray();
        }
    }
  
     public class BattleWatchingResPacket
     {
         public Int16 Result;

         public bool FromBytes(byte[] bodyData)
         {
            Result = BitConverter.ToInt16(bodyData, 0);
            return true;
         }
     }



    public class BattleWatchingEndReqPacket
    {
    }

    public class BattleWatchingEndResPacket
    {
        public Int16 Result;

        public bool FromBytes(byte[] bodyData)
        {
            Result = BitConverter.ToInt16(bodyData, 0);
            return true;
        }
    }
    #endregion
}
