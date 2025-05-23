package bench

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
)

var (
	m     map[int64]int = make(map[int64]int, 20)
	mu    sync.Mutex
	round int64 = 1
)

// 替代 tls.ID() 的标准库实现
func getGoroutineID() int64 {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	id, _ := strconv.Atoi(idField)
	return int64(id)
}

func BenchmarkParalell(b *testing.B) {
	goroutineID := getGoroutineID()
	fmt.Printf("\ngoroutine[%d] enter BenchmarkParalell: round[%d], b.N[%d]\n",
		goroutineID, atomic.LoadInt64(&round), b.N)
	defer func() {
		atomic.AddInt64(&round, 1)
	}()

	b.RunParallel(func(pb *testing.PB) {
		id := goroutineID
		fmt.Printf("goroutine[%d] enter loop func in BenchmarkParalell: round[%d], b.N[%d]\n", goroutineID, atomic.LoadInt64(&round), b.N)
		for pb.Next() {
			mu.Lock()
			_, ok := m[id]
			if !ok {
				m[id] = 1
			} else {
				m[id] = m[id] + 1
			}
			mu.Unlock()
		}

		mu.Lock()
		count := m[id]
		mu.Unlock()

		fmt.Printf("goroutine[%d] exit loop func in BenchmarkParalell: round[%d], loop[%d]\n", goroutineID, atomic.LoadInt64(&round), count)
	})

	fmt.Printf("goroutine[%d] exit BenchmarkParalell: round[%d], b.N[%d]\n",
		goroutineID, atomic.LoadInt64(&round), b.N)
}
