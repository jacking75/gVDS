package scommon

import "errors"

var ErrRingBufferSizePositive = errors.New("ringBuffer: maxSize must be positive")
var ErrRingBufferSizeGreaterMaxPacketSize = errors.New("ringBuffer: maxSize must be greater than or equal to maxPacketSize")
var ErrRingBufferPacketSizePositive = errors.New("ringBuffer: maxPacketSize must be positive")
var ErrRingBufferWritePacketTooLarge = errors.New("ringBuffer: packet Data is too large")
var ErrRingBufferWritePacketRemainSize = errors.New("ringBuffer: packet Data is less Remain")
var ErrRingBufferWritePacketSpaceAllocDisable = errors.New("ringBuffer: WritePacketSpaceAllocDisable")
var ErrRingBufferWritePacketSpaceAllocTooLarge = errors.New("ringBuffer: write packet space is too large")
var ErrRingBufferWritePacketSpaceAllocRemainSize = errors.New("ringBuffer: write packet space is less Remain")
var ErrRingBufferDisableSpendPacket = errors.New("ringBuffer: disable spend packet")
var ErrRingBufferSpendPacketExpectedCursor = errors.New("ringBuffer: spend packet expectedCursor is too large")
var ErrRingBufferDisableGetView = errors.New("ringBuffer: disable get view")
var ErrRingBufferGetViewExpectedCursor = errors.New("ringBuffer: get view expectedCursor is too large")