package main

import (
	"encoding/binary"
	. "scommon"
)

const (
	ERROR_CODE_NONE = 1

	ERROR_CODE_LOGIN_USER_ALREADY  = 31


)


const (
	PACKET_ID_DEV_ECHO_REQ = 192
	PACKET_ID_DEV_ECHO_RES = 193

	PACKET_ID_LOGIN_REQ = 701
	PACKET_ID_LOGIN_RES = 702
)

const PACKET_HEADER_SIZE = 5

type Header struct {
	TotalSize  int16
	ID         int16
	PacketType int8 // 비트 필드로 데이터 설정. 0 이면 Normal, 1번 비트 On(압축), 2번 비트 On(암호화)
}

// Header의 PacketID만 읽는다
func peekPacketID(rawData []byte) int16 {
	packetID := binary.LittleEndian.Uint16(rawData[2:])
	return int16(packetID)
}

// 보디 데이터의 참조만 가져간다
func peekPacketBody(rawData []byte) (bodySize int16, refBody []byte) {
	headerSize := int16(PACKET_HEADER_SIZE)
	totalSize := int16(binary.LittleEndian.Uint16(rawData))
	bodySize = totalSize - headerSize

	if bodySize > 0 {
		refBody = rawData[headerSize:]
	}

	return bodySize, refBody
}

func decodingPacketHeader(header *Header, data []byte) {
	reader := MakeBReader(data, true)
	header.TotalSize, _ = reader.ReadS16()
	header.ID, _ = reader.ReadS16()
	header.PacketType, _ = reader.ReadS8()
}

func encodingPacketHeader(writer *RawBinaryData, totalSize int16, pktId int16, packetType int8) {
	writer.WriteS16(totalSize)
	writer.WriteS16(pktId)
	writer.WriteS8(packetType)
}



// [[[ 로그인 ]]] PACKET_ID_LOGIN_REQ
const MAX_USER_ID_BYTE_LENGTH      = 16

type LoginReqPacket struct {
	UserID []byte
	AUthCode uint64
}

func (loginReq LoginReqPacket) EncodingPacket() ([]byte, int16) {
	totalSize := int16(PACKET_HEADER_SIZE + MAX_USER_ID_BYTE_LENGTH + 8)
	sendBuf := make([]byte, totalSize)

	writer := MakeBWriter(sendBuf, true)
	encodingPacketHeader(&writer, totalSize, PACKET_ID_LOGIN_REQ, 0)
	writer.WriteBytes(loginReq.UserID[:])
	writer.WriteU64(loginReq.AUthCode)
	return sendBuf, totalSize
}

func (loginReq *LoginReqPacket) Decoding(bodyData []byte) bool {
	bodySize := MAX_USER_ID_BYTE_LENGTH + 8
	if len(bodyData) != bodySize {
		return false
	}

	reader := MakeBReader(bodyData, true)
	loginReq.UserID = reader.ReadBytes(MAX_USER_ID_BYTE_LENGTH)
	loginReq.AUthCode, _ = reader.ReadU64()
	return true
}


type LoginResPacket struct {
	Result int16
}

func (loginRes LoginResPacket) EncodingPacket(sendBuf []byte) ([]byte, int16) {
	totalSize := int16(PACKET_HEADER_SIZE + 2)

	writer := MakeBWriter(sendBuf, true)
	encodingPacketHeader(&writer, totalSize, PACKET_ID_LOGIN_RES, 0)
	writer.WriteS16(loginRes.Result)
	return sendBuf[:totalSize], totalSize
}
