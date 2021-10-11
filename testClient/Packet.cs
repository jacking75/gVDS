using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace omokClient
{
    public enum PacketID
   {

        PACKET_ID_DEV_ECHO = 92,        

        PACKET_ID_LOGIN_REQ = 701, 
        PACKET_ID_LOGIN_RES = 702,
    }


    class PacketDef
    {
        public const UInt16 PACKET_HEADER_SIZE = 5;
        public const int MAX_USER_ID_BYTE_LENGTH = 16;
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


    public class LoginReqPacket
    {
        public string UserID;
        public UInt64 AuthCode;
        
        public byte[] ToBytes(PacketID packetID)
        {
            var userID = new byte[PacketDef.MAX_USER_ID_BYTE_LENGTH];
            Encoding.UTF8.GetBytes(UserID).CopyTo(userID, 0);

            const int packetSize = PacketDef.PACKET_HEADER_SIZE + PacketDef.MAX_USER_ID_BYTE_LENGTH + 8;

            List<byte> dataSource = new List<byte>();
            dataSource.AddRange(BitConverter.GetBytes((Int16)packetSize));
            dataSource.AddRange(BitConverter.GetBytes((Int16)packetID));
            dataSource.AddRange(BitConverter.GetBytes((SByte)0));
            dataSource.AddRange(userID);
            dataSource.AddRange(BitConverter.GetBytes(AuthCode));
            return dataSource.ToArray();
        }
    }
  
     public class LoginResPacket
    {
         public Int16 Result;

         public bool FromBytes(byte[] bodyData)
         {
            Result = BitConverter.ToInt16(bodyData, 0);
            return true;
         }
     }


}
