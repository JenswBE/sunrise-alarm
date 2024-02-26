// Based on https://github.com/metakeule/loop

package loop

type Loop struct {
	data    []byte
	pointer int
	length  int
}

// New creates a new loop reader.
// Will panic if provided data is nil or empty.
func New(data []byte) *Loop {
	if len(data) == 0 {
		panic("Loop reader: Provided data is nil or empty")
	}

	return &Loop{
		data:    data,
		pointer: 0,
		length:  len(data),
	}
}

func (l *Loop) Reset() { l.pointer = 0 }

// Read returns bytes and loops data if needed.
// Will always return nil for err.
func (l *Loop) Read(p []byte) (n int, err error) {
	for i := range len(p) {
		dataPos := (l.pointer + i) % l.length
		p[i] = l.data[dataPos]
	}
	l.pointer = (l.pointer + len(p)) % l.length
	return len(p), nil
}
