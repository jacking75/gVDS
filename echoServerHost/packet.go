package main

import (
	"encoding/binary"
	. "scommon"
)

const (
	PACKET_ID_DEV_ECHO_REQ = 192
	PACKET_ID_DEV_ECHO_RES = 193
)

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
	headerSize := int16(5)
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


