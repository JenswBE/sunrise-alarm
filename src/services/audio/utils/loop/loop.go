// Based on https://github.com/metakeule/loop

package loop

type Loop struct {
	data    []byte
	pointer int
	length  int
}

func New(data []byte) *Loop {
	return &Loop{
		data:    data,
		pointer: 0,
		length:  len(data),
	}
}

func (l *Loop) Reset() { l.pointer = 0 }

// return bytes as long as you are asked for it
func (l *Loop) Read(p []byte) (n int, err error) {
	for i := 0; i < len(p); i++ {
		dataPos := (l.pointer + i) % l.length
		p[i] = l.data[dataPos]
	}
	l.pointer = (l.pointer + len(p)) % l.length
	return len(p), nil
}
