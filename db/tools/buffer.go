package tools


import (
	"bytes"
	"fmt"
	"io"
)

type Buffer struct {
	date     []byte
	lastRead int
}

func (b *Buffer) Next(i int) (data []byte) {
	if i > b.Len()-b.lastRead {
		b.lastRead = 0
		return []byte{}
	}

	data = b.date[b.lastRead : b.lastRead+i+1]
	b.lastRead += i

	return data
}

func NewBuffer(b []byte) *Buffer {
	if b != nil {
		return &Buffer{date: bytes.Clone(b)}
	}
	return &Buffer{date: make([]byte, 0)}
}

func (b *Buffer) Reset() {
	if len(b.date) != 0 {
		b.date = b.date[:0]
	}
}

func (b *Buffer) Bytes() []byte {
	return b.date
}

func (b *Buffer) Len() int {
	return len(b.date)
}
func (b *Buffer) WriteFromIndex(d []byte, i int) {
	if len(d) > len(b.date) || i+1 > len(b.date) {
		return
	}
	for z := 0; z < len(d); z++ {
		b.date[i] = d[z]
		i++
	}
}

func (b *Buffer) Ints(a any) {
	switch a.(type) {
	case int64:
		g := []byte{
			byte(0xff & (a.(int64))),
			byte(0xff & (a.(int64) >> 8)),
			byte(0xff & (a.(int64) >> 16)),
			byte(0xff & (a.(int64) >> 24)),
			byte(0xff & (a.(int64) >> 32)),
			byte(0xff & (a.(int64) >> 40)),
			byte(0xff & (a.(int64) >> 48)),
			byte(0xff & (a.(int64) >> 56)),
		}
		b.Write(g)
		break
	case int8:
		b.WriteBytes(byte(a.(int8)))
		break
	case int32:
		g := []byte{
			byte(0xff & (a.(int32))),
			byte(0xff & (a.(int32) >> 8)),
			byte(0xff & (a.(int32) >> 16)),
			byte(0xff & (a.(int32) >> 24)),
		}
		b.Write(g)
		break
	default:
		return
	}

}

func Ints(a any) []byte {
	switch a.(type) {
	case int64:
		g := []byte{
			byte(0xff & (a.(int64))),
			byte(0xff & (a.(int64) >> 8)),
			byte(0xff & (a.(int64) >> 16)),
			byte(0xff & (a.(int64) >> 24)),
			byte(0xff & (a.(int64) >> 32)),
			byte(0xff & (a.(int64) >> 40)),
			byte(0xff & (a.(int64) >> 48)),
			byte(0xff & (a.(int64) >> 56)),
		}
		return g
	case int32:
		g := []byte{
			byte(0xff & (a.(int32))),
			byte(0xff & (a.(int32) >> 8)),
			byte(0xff & (a.(int32) >> 16)),
			byte(0xff & (a.(int32) >> 24)),
		}
		return g
	default:
		return nil
	}

}

func (b *Buffer) WriteToFile(i io.Writer) {
	i.Write(b.date)
}

func (b *Buffer) ReadFile(i io.Reader) {
	tm := make([]byte, 20)
	for {
		s, err := i.Read(tm)
		if err != nil {
			break
		}

		b.Write(tm[:s])
	}

}

func (b *Buffer) Grow(Bytes int) {

	tem := bytes.Clone(b.date)
	b.date = make([]byte, len(tem)+Bytes)
	copy(b.date[:len(tem)], tem)
}
func (b *Buffer) WriteString(str string) {
	b.date = append(b.date, str...)
}
func (b *Buffer) WriteBytes(c byte) {
	b.date = append(b.date, c)
}

func (b *Buffer) Write(c []byte) {
	if len(c) == 0 {
		return
	}
	b.date = append(b.date, c...)
}

func (b *Buffer) ShiftRigth(s int) {
	tep := b.Bytes()
	b.date = make([]byte, s+len(tep))
	j := 0
	fmt.Println(len(b.date), s)
	for i := s; i < len(b.date); i++ {
		b.date[i] = tep[j]
		j++
	}

}

func (b *Buffer) AppendFromIdex(index int, data []byte) {

	if index >= len(b.date) {
		b.date = append(b.date, data...)
	} else if index == 0 {
		t := make([]byte, len(b.date))
		copy(t, b.date)
		b.Reset()
		b.date = append(b.date, data...)
		b.date = append(b.date, t...)
	} else {
		 L := make([]byte , len(b.Bytes()[:index]) )
		 R := make([]byte , len (b.Bytes()[index:]) )
		copy(L, b.Bytes()[:index])
		copy(R, b.Bytes()[index:])
		b.Reset()
		b.date = append(b.date, L...)
		b.date = append(b.date, data...)
		b.date = append(b.date, R...)
	}
}
