package bufferpool

import (
	"bytes"
	"sync"
)

var pool = sync.Pool{
	New: func() any {
		return &bytes.Buffer{}
	},
}

func Get() *bytes.Buffer {
	return pool.Get().(*bytes.Buffer) //nolint:errcheck //best effort
}

func Put(buf *bytes.Buffer) {
	buf.Reset()
	pool.Put(buf)
}
