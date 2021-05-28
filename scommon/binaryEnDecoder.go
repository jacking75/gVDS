// 출처: https://github.com/gonet2/agent
package scommon

import (
	"encoding/binary"
	"errors"
	"reflect"
)

// 타입의 크기를 계산한다
func Sizeof(t reflect.Type) int {
	switch t.Kind() {
	case reflect.Array:
		//fmt.Println("reflect.Array")
		if s := Sizeof(t.Elem()); s >= 0 {
			return s * t.Len()
		}

	case reflect.Struct:
		//fmt.Println("reflect.Struct")
		sum := 0
		for i, n := 0, t.NumField(); i < n; i++ {
			s := Sizeof(t.Field(i).Type)
			if s < 0 {
				return -1
			}
			sum += s
		}
		return sum

	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128:
		//fmt.Println("reflect.int")
		return int(t.Size())
	case reflect.Slice:
		//fmt.Println("reflect.Slice:", sizeof(t.Elem()))
		return 0
	}

	return -1

}

type RawBinaryData struct {
	pos   int
	data  []byte
	order binary.ByteOrder
}

func MakeBReader(buffer []byte, isLittleEndian bool) RawBinaryData {
	if isLittleEndian {
		return RawBinaryData{data: buffer, order: binary.LittleEndian}
	}
	return RawBinaryData{data: buffer, order: binary.BigEndian}
}

func MakeBWriter(buffer []byte, isLittleEndian bool) RawBinaryData {
	if isLittleEndian {
		return RawBinaryData{data: buffer, order: binary.LittleEndian}
	}
	return RawBinaryData{data: buffer, order: binary.BigEndian}
}

func (p *RawBinaryData) Data() []byte {
	return p.data
}

func (p *RawBinaryData) Length() int {
	return len(p.data)
}


//=============================================== Readers
func (p *RawBinaryData) ReadBool() (ret bool, err error) {
	b, _err := p.ReadByte()

	if b != byte(1) {
		return false, _err
	}

	return true, _err
}

func (p *RawBinaryData) ReadS8() (ret int8, err error) {
	_ret, _err := p.ReadByte()
	ret = int8(_ret)
	err = _err
	return
}

func (p *RawBinaryData) ReadU16() (ret uint16, err error) {
	if p.pos+2 > len(p.data) {
		err = errors.New("read uint16 failed")
		return
	}
	buf := p.data[p.pos : p.pos+2]
	ret = p.order.Uint16(buf)
	p.pos += 2
	return
}

func (p *RawBinaryData) ReadS16() (ret int16, err error) {
	_ret, _err := p.ReadU16()
	ret = int16(_ret)
	err = _err
	return
}

func (p *RawBinaryData) ReadU32() (ret uint32, err error) {
	if p.pos+4 > len(p.data) {
		err = errors.New("read uint32 failed")
		return
	}

	buf := p.data[p.pos : p.pos+4]
	ret = p.order.Uint32(buf)
	p.pos += 4
	return
}

func (p *RawBinaryData) ReadS32() (ret int32, err error) {
	_ret, _err := p.ReadU32()
	ret = int32(_ret)
	err = _err
	return
}

func (p *RawBinaryData) ReadU64() (ret uint64, err error) {
	if p.pos+8 > len(p.data) {
		err = errors.New("read uint64 failed")
		return
	}

	buf := p.data[p.pos : p.pos+8]
	ret = p.order.Uint64(buf)
	p.pos += 8
	return
}

func (p *RawBinaryData) ReadS64() (ret int64, err error) {
	_ret, _err := p.ReadU64()
	ret = int64(_ret)
	err = _err
	return
}

func (p *RawBinaryData) ReadByte() (ret byte, err error) {
	if p.pos >= len(p.data) {
		err = errors.New("read byte failed")
		return
	}

	ret = p.data[p.pos]
	p.pos++
	return
}

func (p *RawBinaryData) ReadBytes(readSize int) (refSlice []byte) {
	refSlice = p.data[p.pos : p.pos+readSize]
	p.pos += readSize
	return
}

func (p *RawBinaryData) ReadString() (ret string, err error) {
	if p.pos+2 > len(p.data) {
		err = errors.New("read string header failed")
		return
	}

	size, _ := p.ReadU16()
	if p.pos+int(size) > len(p.data) {
		err = errors.New("read string Data failed")
		return
	}

	bytes := p.data[p.pos : p.pos+int(size)]
	p.pos += int(size)
	ret = string(bytes)
	return
}



//================================================ Writers
func (p *RawBinaryData) WriteS8(v int8) {
	p.data[p.pos] = (byte)(v)
	p.pos++
}

func (p *RawBinaryData) WriteU16(v uint16) {
	p.order.PutUint16(p.data[p.pos:], v)
	p.pos += 2
}

func (p *RawBinaryData) WriteS16(v int16) {
	p.WriteU16(uint16(v))
}

func (p *RawBinaryData) WriteBytes(v []byte) {
	copy(p.data[p.pos:], v)
	p.pos += len(v)
}

func (p *RawBinaryData) WriteU32(v uint32) {
	p.order.PutUint32(p.data[p.pos:], v)
	p.pos += 4
}

func (p *RawBinaryData) WriteS32(v int32) {
	p.WriteU32(uint32(v))
}

func (p *RawBinaryData) WriteU64(v uint64) {
	p.order.PutUint64(p.data[p.pos:], v)
	p.pos += 8
}

func (p *RawBinaryData) WriteS64(v int64) {
	p.WriteU64(uint64(v))
}

func (p *RawBinaryData) WriteString(v string) {
	copyLen := copy(p.data[p.pos:], v)
	p.pos += copyLen
}
