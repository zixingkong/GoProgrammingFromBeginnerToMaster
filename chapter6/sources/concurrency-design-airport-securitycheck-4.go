package main

import (
	"fmt"
	"sync"
	"time"
)

// Passenger 旅客结构体
type Passenger struct {
	ID int
}

func main() {

	// 创建安检步骤的管道
	idCheckChan := make(chan Passenger)
	bodyCheckChan := make(chan Passenger)
	luggageCheckChan := make(chan Passenger)
	done := make(chan struct{})

	// 启动安检流水线
	go idCheck(idCheckChan, bodyCheckChan)
	go bodyCheck(bodyCheckChan, luggageCheckChan)
	go luggageCheck(luggageCheckChan, done)

	// 模拟旅客到达（并发）
	var wg sync.WaitGroup
	for i := 1; i <= 30; i++ { // 生成5个旅客
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			fmt.Printf("旅客 %d 到达安检口\n", id)
			idCheckChan <- Passenger{ID: id}
		}(i)
	}

	// 等待所有旅客进入安检流程
	wg.Wait()
	close(idCheckChan) // 关闭管道触发流水线结束
	// 这段代码 `<-done` 的含义是 **等待安检流水线完成信号**，具体解释如下：

	// 1. **`done`通道的作用**：
	//    - 这是一个`chan struct{}`类型的通道，专门用于传递完成信号
	//    - 当行李检查环节(`luggageCheck`)处理完所有旅客后，会执行`close(done)`关闭通道

	// 2. **`<-done`的行为**：
	//    - 从`done`通道接收值会阻塞，直到通道被关闭
	//    - 通道关闭后，`<-done`会立即返回零值（这里是`struct{}`）
	//    - 这种模式是Go中常用的"等待完成"同步机制

	// 3. **在整个流程中的意义**：
	//    - 确保主goroutine等待所有旅客完成全部安检流程（身份验证→人身检查→行李检查）
	//    - 只有收到完成信号后，程序才会继续执行（这里是结束程序）

	// 4. **与`sync.WaitGroup`的区别**：
	//    - 前文的`wg.Wait()`是等待所有旅客进入安检流程
	//    - 这里的`<-done`是等待所有旅客完成全部安检流程

	// 这种设计模式在Go并发编程中非常常见，用于精确控制goroutine的执行顺序和生命周期。
	// 等待所有安检流程完成
	<-done
}

// 身份验证
func idCheck(in <-chan Passenger, out chan<- Passenger) {
	for p := range in {
		processTime := time.Duration(60) * time.Millisecond
		time.Sleep(processTime)
		fmt.Printf("旅客 %d 完成身份验证 (耗时 %v)\n", p.ID, processTime)
		out <- p
	}
	close(out)
}

// 人身检查
func bodyCheck(in <-chan Passenger, out chan<- Passenger) {
	for p := range in {
		processTime := time.Duration(20) * time.Millisecond
		time.Sleep(processTime)
		fmt.Printf("旅客 %d 完成人身检查 (耗时 %v)\n", p.ID, processTime)
		out <- p
	}
	close(out)
}

// 行李检查
func luggageCheck(in <-chan Passenger, done chan<- struct{}) {
	for p := range in {
		processTime := time.Duration(180) * time.Millisecond
		time.Sleep(processTime)
		fmt.Printf("旅客 %d 完成行李检查 (耗时 %v)\n", p.ID, processTime)
	}
	close(done)
}
