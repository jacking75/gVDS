package omokServer

type sendRingBuffer struct {
	data          []byte
	allocSize     int // data의 실제 할당 크기
	writeCursor   int
}

// 링버퍼를 한바퀴 돌 때 앞에 데이터를 다 사용했는지 체크 하지 않는다. 즉 낙관적이다
// 그래서 버퍼의 크기(maxSize)는 꽤 넉넉해야 한다.
func newSendBuffer(maxSize int) *sendRingBuffer {
	if maxSize <= 0 {
		return nil
	}

	b := &sendRingBuffer{
		allocSize:     maxSize,
		data:          make([]byte, maxSize),
		writeCursor:   0,
	}
	return b
}

func (b *sendRingBuffer) reset() {
	b.writeCursor = 0
}

func (b *sendRingBuffer) getBuffer(requiredSize int) []byte {
	remain := b.allocSize - b.writeCursor
	if remain < requiredSize {
		b.writeCursor = 0
	}

	return b.data[b.writeCursor:]
}

func (b *sendRingBuffer) aheadWriteCursor(size int) {
	b.writeCursor += size
}