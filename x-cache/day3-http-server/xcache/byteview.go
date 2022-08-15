package xcache

// 只读的缓存值，byte 存储任意类型
type ByteView struct {
	b []byte
}

func (v ByteView) Len() int64 {
	return int64(len(v.b))
}

func (v ByteView) String() string {
	return string(v.b)
}

func (v ByteView) ByteSlice() []byte {
	return cloneByte(v.b)
}

// 返回拷贝，避免被外部程序修改，
func cloneByte(oldest []byte) []byte {
	newByte := make([]byte, len(oldest))
	copy(newByte, oldest)
	return newByte
}
