package bufferpool

import (
	"sync"
	"testing"
)

func TestGetPut(t *testing.T) {
	buf := Get()
	if buf == nil {
		t.Fatal("Get returned nil")
	}
	buf.WriteString("hello")
	if buf.Len() != 5 {
		t.Fatalf("expected len 5, got %d", buf.Len())
	}
	Put(buf)

	buf2 := Get()
	if buf2.Len() != 0 {
		t.Fatalf("expected empty buffer after Put, got len %d", buf2.Len())
	}
	Put(buf2)
}

func TestConcurrentAccess(t *testing.T) {
	var wg sync.WaitGroup
	for range 100 {
		wg.Go(func() {
			buf := Get()
			buf.WriteString("data")
			Put(buf)
		})
	}
	wg.Wait()
}
