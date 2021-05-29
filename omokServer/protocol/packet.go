package protocol

import (
	"encoding/binary"
	"reflect"

	. "scommon"
)

const (
	PACKET_TYPE_NORMAL   = 0
	PACKET_TYPE_COMPRESS = 1
	PACKET_TYPE_SECURE   = 2
)

const (
	MAX_USER_ID_BYTE_LENGTH      = 16
	MAX_USER_PW_BYTE_LENGTH      = 16
	MAX_CHAT_MESSAGE_BYTE_LENGTH = 126
)

type Header struct {
	TotalSize  int16
	ID         int16
	PacketType int8 // 비트 필드로 데이터 설정. 0 이면 Normal, 1번 비트 On(압축), 2번 비트 On(암호화)
}

var _clientSessionHeaderSize int16
var _ServerSessionHeaderSize int16

func Init_packet() {
	_clientSessionHeaderSize = protocolInitHeaderSize()
	_ServerSessionHeaderSize = protocolInitHeaderSize()
}

func protocolInitHeaderSize() int16 {
	var packetHeader Header
	headerSize := Sizeof(reflect.TypeOf(packetHeader))
	return (int16)(headerSize)
}

// Header의 PacketID만 읽는다
func PeekPacketID(rawData []byte) int16 {
	packetID := binary.LittleEndian.Uint16(rawData[2:])
	return int16(packetID)
}

// 보디데이터의 참조만 가져간다
func PeekPacketBody(rawData []byte) (bodySize int16, refBody []byte) {
	headerSize := _clientSessionHeaderSize
	totalSize := int16(binary.LittleEndian.Uint16(rawData))
	bodySize = totalSize - headerSize

	if bodySize > 0 {
		refBody = rawData[headerSize:]
	}

	return bodySize, refBody
}

func DecodingPacketHeader(header *Header, data []byte) {
	reader := MakeBReader(data, true)
	header.TotalSize, _ = reader.ReadS16()
	header.ID, _ = reader.ReadS16()
	header.PacketType, _ = reader.ReadS8()
}

func EncodingPacketHeader(writer *RawBinaryData, totalSize int16, pktId int16, packetType int8) {
	writer.WriteS16(totalSize)
	writer.WriteS16(pktId)
	writer.WriteS8(packetType)
}


///<<< 패킷 인코딩/디코딩

// [[[ 로그인 ]]] PACKET_ID_LOGIN_REQ
type LoginReqPacket struct {
	UserID []byte
	AUthCode uint64
}

func (loginReq LoginReqPacket) EncodingPacket() ([]byte, int16) {
	totalSize := _clientSessionHeaderSize + MAX_USER_ID_BYTE_LENGTH + 8
	sendBuf := make([]byte, totalSize)

	writer := MakeBWriter(sendBuf, true)
	EncodingPacketHeader(&writer, totalSize, PACKET_ID_LOGIN_REQ, 0)
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
	totalSize := _clientSessionHeaderSize + 2

	writer := MakeBWriter(sendBuf, true)
	EncodingPacketHeader(&writer, totalSize, PACKET_ID_LOGIN_RES, 0)
	writer.WriteS16(loginRes.Result)
	return sendBuf[:totalSize], totalSize
}



// [[[ 허트비트 ]]] PACKET_ID_HEARTBEAT_REQ
type HeartBeatReqPacket struct {
	Expected int64
}

func (req HeartBeatReqPacket) EncodingPacket(sendBuf []byte) ([]byte, int16) {
	totalSize := _clientSessionHeaderSize + 8

	writer := MakeBWriter(sendBuf, true)
	EncodingPacketHeader(&writer, totalSize, PACKET_ID_HEARTBEAT_REQ, 0)
	writer.WriteS64(req.Expected)
	return sendBuf[:totalSize], totalSize
}


type HeartBeatResPacket struct {
	Expected int64
}

func (res *HeartBeatResPacket) Decoding(bodyData []byte) bool {
	bodySize := 8
	if len(bodyData) != bodySize {
		return false
	}

	reader := MakeBReader(bodyData, true)
	res.Expected, _ = reader.ReadS64()
	return true
}



// [[[  ]]]   PACKET_ID_ERROR_NTF
type ErrorNtfPacket struct {
	ErrorCode int16
}

func (response ErrorNtfPacket) EncodingPacket(errorCode int16) ([]byte, int16) {
	totalSize := _clientSessionHeaderSize + 2
	sendBuf := make([]byte, totalSize)

	writer := MakeBWriter(sendBuf, true)
	EncodingPacketHeader(&writer, totalSize, PACKET_ID_ERROR_NTF, 0)
	writer.WriteS16(errorCode)
	return sendBuf, totalSize
}

func (response *ErrorNtfPacket) Decoding(bodyData []byte) bool {
	if len(bodyData) != 2 {
		return false
	}

	reader := MakeBReader(bodyData, true)
	response.ErrorCode, _ = reader.ReadS16()
	return true
}


/// [ 방 입장 ]
type RoomEnterReqPacket struct {
	RoomNumber int32
}

func (request RoomEnterReqPacket) EncodingPacket() ([]byte, int16) {
	totalSize := _clientSessionHeaderSize + (4)
	sendBuf := make([]byte, totalSize)

	writer := MakeBWriter(sendBuf, true)
	EncodingPacketHeader(&writer, totalSize, PACKET_ID_ROOM_ENTER_REQ, 0)
	writer.WriteS32(request.RoomNumber)
	return sendBuf, totalSize
}

func (request *RoomEnterReqPacket) Decoding(bodyData []byte) bool {
	if len(bodyData) != (4) {
		return false
	}

	reader := MakeBReader(bodyData, true)
	request.RoomNumber, _ = reader.ReadS32()
	return true
}


type RoomEnterResPacket struct {
	Result int16
	RoomNumber int32
	RoomUserUniqueId uint64
}

func (response RoomEnterResPacket) EncodingPacket() ([]byte, int16) {
	totalSize := _clientSessionHeaderSize + 2 + 4 + 8
	sendBuf := make([]byte, totalSize)

	writer := MakeBWriter(sendBuf, true)
	EncodingPacketHeader(&writer, totalSize, PACKET_ID_ROOM_ENTER_RES,  0)
	writer.WriteS16(response.Result)
	writer.WriteS32(response.RoomNumber)
	writer.WriteU64(response.RoomUserUniqueId)
	return sendBuf, totalSize
}

func (response *RoomEnterResPacket) Decoding(bodyData []byte) bool {
	if len(bodyData) != (2+4+8) {
		return false
	}

	reader := MakeBReader(bodyData, true)
	response.Result, _ = reader.ReadS16()
	response.RoomNumber, _ = reader.ReadS32()
	response.RoomUserUniqueId, _ = reader.ReadU64()
	return true
}


///  방에 있는 있는 유저 리스트 통보(채널에 들어오는 유저에게)
type RoomUserInfoPktData struct {
	UniqueId int64
	IDLen int8
	ID    []byte
}

type RoomUserListNtfPacket struct {
	UserCount int8
	UserList  []byte //RoomUserInfoPktData
}

func (notify RoomUserListNtfPacket) EncodingPacket(userInfoListSize int16) ([]byte, int16) {
	bodySize := 1 + userInfoListSize
	totalSize := _clientSessionHeaderSize + bodySize
	sendBuf := make([]byte, totalSize)

	writer := MakeBWriter(sendBuf, true)
	EncodingPacketHeader(&writer, totalSize, PACKET_ID_ROOM_USER_LIST_NTF, 0)
	writer.WriteS8(notify.UserCount)
	writer.WriteBytes(notify.UserList)
	return sendBuf, totalSize
}

func (notify RoomUserListNtfPacket) Decoding(bodyData []byte) bool {
	reader := MakeBReader(bodyData, true)
	notify.UserCount, _ = reader.ReadS8()
	notify.UserList = reader.ReadBytes(len(bodyData) - 1)
	return true
}


// 채널에 있는 유저들에게 새 유저의 정보를 알려준다
type RoomNewUserNtfPacket struct {
	User []byte //RoomUserInfoPktData
}

func (notify RoomNewUserNtfPacket) EncodingPacket(userInfoSize int16) ([]byte, int16) {
	totalSize := _clientSessionHeaderSize + userInfoSize
	sendBuf := make([]byte, totalSize)

	writer := MakeBWriter(sendBuf, true)
	EncodingPacketHeader(&writer, totalSize, PACKET_ID_ROOM_NEW_USER_NTF, 0)
	writer.WriteBytes(notify.User)
	return sendBuf, totalSize
}



//<<< 방에서 나가기
type RoomLeaveResPacket struct {
	Result int16
}

func (response RoomLeaveResPacket) EncodingPacket() ([]byte, int16) {
	totalSize := _clientSessionHeaderSize + 2
	sendBuf := make([]byte, totalSize)

	writer := MakeBWriter(sendBuf, true)
	EncodingPacketHeader(&writer, totalSize, PACKET_ID_ROOM_LEAVE_RES, 0)
	return sendBuf, totalSize
}

func (response *RoomLeaveResPacket) Decoding(bodyData []byte) bool {
	reader := MakeBReader(bodyData, true)
	response.Result, _ = reader.ReadS16()
	return true
}


type RoomLeaveUserNtfPacket struct {
	UserUniqueId uint64
}

func (notify RoomLeaveUserNtfPacket) EncodingPacket() ([]byte, int16) {
	totalSize := _clientSessionHeaderSize + 8
	sendBuf := make([]byte, totalSize)

	writer := MakeBWriter(sendBuf, true)
	EncodingPacketHeader(&writer, totalSize, PACKET_ID_ROOM_LEAVE_USER_NTF, 0)
	writer.WriteU64(notify.UserUniqueId)
	return sendBuf, totalSize
}

func (notify RoomLeaveUserNtfPacket) Decoding(bodyData []byte) bool {
	if len(bodyData) != 8 {
		return false
	}

	reader := MakeBReader(bodyData, true)
	notify.UserUniqueId, _ = reader.ReadU64()
	return true
}




/// [ 방 채팅 ]]
type RoomChatReqPacket struct {
	MsgLength int16
	Msgs      []byte
}

func (request RoomChatReqPacket) EncodingPacket() ([]byte, int16) {
	totalSize := _clientSessionHeaderSize + 2 + int16(request.MsgLength)
	sendBuf := make([]byte, totalSize)
	writer := MakeBWriter(sendBuf, true)
	EncodingPacketHeader(&writer, totalSize, PACKET_ID_ROOM_CHAT_REQ, 0)

	writer.WriteS16(request.MsgLength)
	writer.WriteBytes(request.Msgs)
	return sendBuf, totalSize
}

func (request *RoomChatReqPacket) Decoding(bodyData []byte) bool {
	bodyLength := len(bodyData)
	if bodyLength < 2 {
		return false
	}

	reader := MakeBReader(bodyData, true)
	request.MsgLength, _ = reader.ReadS16()

	if bodyLength != int((2 + request.MsgLength)) {
		return false
	}

	request.Msgs = bodyData[2:]
	return true
}


type RoomChatResPacket struct {
	Result int16
}

func (response RoomChatResPacket) EncodingPacket() ([]byte, int16) {
	totalSize := _clientSessionHeaderSize + 2
	sendBuf := make([]byte, totalSize)

	writer := MakeBWriter(sendBuf, true)
	EncodingPacketHeader(&writer, totalSize, PACKET_ID_ROOM_CHAT_RES, 0)
	return sendBuf, totalSize
}

func (response *RoomChatResPacket) Decoding(bodyData []byte) bool {
	reader := MakeBReader(bodyData, true)
	response.Result, _ = reader.ReadS16()
	return true
}


type RoomChatNtfPacket struct {
	RoomUserUniqueId uint64
	MsgLen        int16
	Msg           []byte
}

func (response RoomChatNtfPacket) EncodingPacket() ([]byte, int16) {
	totalSize := _clientSessionHeaderSize + 8 + int16(2) + response.MsgLen
	sendBuf := make([]byte, totalSize)
	writer := MakeBWriter(sendBuf, true)
	EncodingPacketHeader(&writer, totalSize, PACKET_ID_ROOM_CHAT_NOTIFY, 0)

	writer.WriteU64(response.RoomUserUniqueId)
	writer.WriteS16(response.MsgLen)
	writer.WriteBytes(response.Msg)
	return sendBuf, totalSize
}

func (response *RoomChatNtfPacket) Decoding(bodyData []byte) bool {
	reader := MakeBReader(bodyData, true)
	response.RoomUserUniqueId, _ = reader.ReadU64()
	response.MsgLen, _ = reader.ReadS16()
	response.Msg = reader.ReadBytes(int(response.MsgLen))
	return true
}



///<<< Room Relay
type RoomRelayReqPacket struct {
	Data      []byte
}

func (request RoomRelayReqPacket) EncodingPacket(size int16) ([]byte, int16) {
	totalSize := _clientSessionHeaderSize + int16(len(request.Data))
	sendBuf := make([]byte, totalSize)
	writer := MakeBWriter(sendBuf, true)
	EncodingPacketHeader(&writer, totalSize, PACKET_ID_ROOM_RELAY_REQ, 0)

	writer.WriteBytes(request.Data)
	return sendBuf, totalSize
}

func (request *RoomRelayReqPacket) Decoding(bodyData []byte) bool {
	reader := MakeBReader(bodyData, true)
	request.Data = reader.ReadBytes(len(bodyData))
	return true
}


type RoomRelayNtfPacket struct {
	RoomUserUniqueId uint64
	Data      []byte
}

func (notify RoomRelayNtfPacket) EncodingPacket(size int16) ([]byte, int16) {
	totalSize := _clientSessionHeaderSize + 8 + int16(len(notify.Data))
	sendBuf := make([]byte, totalSize)
	writer := MakeBWriter(sendBuf, true)
	EncodingPacketHeader(&writer, totalSize, PACKET_ID_ROOM_RELAY_NTF, 0)

	writer.WriteU64(notify.RoomUserUniqueId)
	writer.WriteBytes(notify.Data)
	return sendBuf, totalSize
}

func (notify *RoomRelayNtfPacket) Decoding(bodyData []byte) bool {
	reader := MakeBReader(bodyData, true)
	notify.RoomUserUniqueId, _ = reader.ReadU64()
	notify.Data = reader.ReadBytes(len(bodyData) - 8)
	return true
}