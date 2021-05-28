package smallNet

type ringBuffer struct {
	_data        []byte
	_allocSize   int // data의 실제 할당 크기
	_maxSize     int // data의 최대 크기

	_maxPacketSize int

	_writeCursor int
	_readCursor  int
}

type ringbufferErr int

const (
	err_ringbuffer_none = 0
	err_ringbuffer_less_zero = 1
	err_ringbuffer_too_small = 2
)

// 링버퍼를 한바퀴 돌 때 앞에 데이터를 다 사용했는지 체크 하지 않는다. 즉 낙관적이다
// 그래서 버퍼의 크기(maxSize)는 꽤 넉넉해야 한다.
func newRingBuffer(maxSize int, maxPacketSize int) (*ringBuffer, ringbufferErr) {
	if maxSize <= 0 {
		return nil, err_ringbuffer_less_zero
	}

	if maxSize < maxPacketSize {
		return nil, err_ringbuffer_too_small
	}

	b := &ringBuffer{
		_maxSize:       maxSize,
		_maxPacketSize: maxPacketSize,
		_allocSize:     maxSize + maxPacketSize,
		_data:          make([]byte, maxSize+maxPacketSize),
		_writeCursor: 0,
		_readCursor: 0,
	}
	return b, err_ringbuffer_none
}

func (b *ringBuffer) reset() {
	b._writeCursor = 0
	b._readCursor = 0
}

func (b *ringBuffer) getWriteBuffer(requiredSize int) []byte {
	remain := b._allocSize - b._writeCursor
	if remain < requiredSize {
		readsize := b._writeCursor - b._readCursor
		if readsize > 0 {
			copy(b._data[0:], b._data[b._readCursor:b._writeCursor])
			b._writeCursor = readsize
			b._readCursor = 0
		} else {
			b._writeCursor = 0
			b._readCursor = 0
		}
	}

	return b._data[b._writeCursor:]
}

func (b *ringBuffer) aheadWriteCursor(size int) {
	b._writeCursor += size
}

func (b *ringBuffer) aheadReadCursor(size int) {
	b._readCursor += size
}

func (b *ringBuffer) aheadWRCursor(size int) {
	b._writeCursor += size
	b._readCursor += size
}

func (b *ringBuffer) readAbleBuffer() ([]byte, int) {
	readAblesize := b._writeCursor - b._readCursor
	return b._data[b._readCursor:b._writeCursor], readAblesize
}