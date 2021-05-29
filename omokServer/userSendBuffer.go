package omokServer

type sendRingBuffer struct {
	_data        []byte
	_allocSize   int // data의 실제 할당 크기
	_writeCursor int
}

// 링버퍼를 한바퀴 돌 때 앞에 데이터를 다 사용했는지 체크 하지 않는다. 즉 낙관적이다
// 그래서 버퍼의 크기(maxSize)는 꽤 넉넉해야 한다.
func newSendBuffer(maxSize int) *sendRingBuffer {
	if maxSize <= 0 {
		return nil
	}

	b := &sendRingBuffer{
		_allocSize:   maxSize,
		_data:        make([]byte, maxSize),
		_writeCursor: 0,
	}
	return b
}

func (b *sendRingBuffer) reset() {
	b._writeCursor = 0
}

func (b *sendRingBuffer) getBuffer(requiredSize int) []byte {
	remain := b._allocSize - b._writeCursor
	if remain < requiredSize {
		b._writeCursor = 0
	}

	return b._data[b._writeCursor:]
}

func (b *sendRingBuffer) aheadWriteCursor(size int) {
	b._writeCursor += size
}