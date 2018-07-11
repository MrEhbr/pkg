package pool

import (
	"bufio"
	"bytes"
)

// BufferPool implements a pool of bytes.Buffers in the form of a bounded channel.
type WriterPool struct {
	c chan *bufio.Writer
}

// NewBufferPool creates a new BufferPool bounded to the given size.
func NewWriterPool(size int) (wp *WriterPool) {
	return &WriterPool{
		c: make(chan *bufio.Writer, size),
	}
}

// Get gets a Buffer from the BufferPool, or creates a new one if none are available in the pool.
func (wp *WriterPool) Get() (buf *bytes.Buffer, w *bufio.Writer) {
	select {
	case w = <-wp.c:
	default:
		buf = bytes.NewBuffer([]byte{})
		w = bufio.NewWriter(buf)
	}
	return
}

// Put returns the given Buffer to the BufferPool.
func (wp *WriterPool) Put(b *bufio.Writer) {
	b.Reset(b)
	select {
	case wp.c <- b:
	default:
	}
}
