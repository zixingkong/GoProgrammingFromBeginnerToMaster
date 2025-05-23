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
	m     = make(map[int64]struct{}, 10)
	mu    sync.Mutex
	round int64 = 1
)

// 替代 tls.ID() 的标准库实现
func getGoroutineID() int {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	id, _ := strconv.Atoi(idField)
	return id
}

func BenchmarkSequential(b *testing.B) {
	goroutineID := getGoroutineID()
	fmt.Printf("\ngoroutine[%d] enter BenchmarkSequential: round[%d], b.N[%d]\n",
		goroutineID, atomic.LoadInt64(&round), b.N)
	defer func() {
		atomic.AddInt64(&round, 1)
	}()

	for i := 0; i < b.N; i++ {
		mu.Lock()
		_, ok := m[round]
		if !ok {
			m[round] = struct{}{}
			fmt.Printf("goroutine[%d] enter loop in BenchmarkSequential: round[%d], b.N[%d]\n",
				goroutineID, atomic.LoadInt64(&round), b.N)
		}
		mu.Unlock()
	}
	fmt.Printf("goroutine[%d] exit BenchmarkSequential: round[%d], b.N[%d]\n",
		goroutineID, atomic.LoadInt64(&round), b.N)
}
