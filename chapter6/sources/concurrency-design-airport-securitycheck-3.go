package main

import "time"

// 版本3：并发方案
// 总耗时：2160ms
const (
	idCheckTmCost   = 60  // 身份证检查耗时(毫秒)
	bodyCheckTmCost = 120 // 身体检查耗时(毫秒)
	xRayCheckTmCost = 180 // X光检查耗时(毫秒)
)

// idCheck 身份证检查函数
// id: 安检通道标识符
// 返回: 检查耗时
func idCheck(id string) int {
	time.Sleep(time.Millisecond * time.Duration(idCheckTmCost))
	print("\tgoroutine-", id, "-idCheck: idCheck ok\n")
	return idCheckTmCost
}

// bodyCheck 身体检查函数
// id: 安检通道标识符
// 返回: 检查耗时
func bodyCheck(id string) int {
	time.Sleep(time.Millisecond * time.Duration(bodyCheckTmCost))
	print("\tgoroutine-", id, "-bodyCheck: bodyCheck ok\n")
	return bodyCheckTmCost
}

// xRayCheck X光检查函数
// id: 安检通道标识符
// 返回: 检查耗时
func xRayCheck(id string) int {
	time.Sleep(time.Millisecond * time.Duration(xRayCheckTmCost))
	print("\tgoroutine-", id, "-xRayCheck: xRayCheck ok\n")
	return xRayCheckTmCost
}

// start 启动一个安检环节的goroutine
// id: 安检环节标识符
// f: 要执行的安检函数
// next: 下一个安检环节的输入通道
// 返回:
//   - queue: 用于接收工作任务的通道
//   - quit: 用于通知退出的通道
//   - result: 用于获取累计耗时的通道
func start(id string, f func(string) int, next chan<- struct{}) (chan<- struct{}, chan<- struct{}, <-chan int) {
	queue := make(chan struct{}, 10) // 工作任务通道(带缓冲)
	quit := make(chan struct{})      // 退出信号通道
	result := make(chan int)         // 结果通道

	go func() {
		total := 0 // 累计耗时
		for {
			select {
			case <-quit: // 收到退出信号
				result <- total // 发送累计耗时
				return          // 退出goroutine
			case v := <-queue: // 收到工作任务
				total += f(id)   // 执行安检函数并累计耗时
				if next != nil { // 如果有下一个环节
					next <- v // 将工作传递给下一个环节
				}
			}
		}
	}()
	return queue, quit, result
}

// newAirportSecurityCheckChannel 创建完整的安检通道
// id: 安检通道标识符
// queue: 输入工作任务的通道
func newAirportSecurityCheckChannel(id string, queue <-chan struct{}) {
	go func(id string) {
		print("goroutine-", id, ": airportSecurityCheckChannel is ready...\n")
		// 启动X光检查环节(最后一个环节)
		queue3, quit3, result3 := start(id, xRayCheck, nil)

		// 启动身体检查环节(中间环节)
		queue2, quit2, result2 := start(id, bodyCheck, queue3)

		// 启动身份证检查环节(第一个环节)
		queue1, quit1, result1 := start(id, idCheck, queue2)

		for {
			select {
			case v, ok := <-queue: // 从主通道接收工作任务
				if !ok { // 如果通道已关闭
					// 关闭所有环节的quit通道
					close(quit1)
					close(quit2)
					close(quit3)
					// 获取各环节的最大耗时
					total := max(<-result1, <-result2, <-result3)
					print("goroutine-", id, ": airportSecurityCheckChannel time cost:", total, "\n")
					print("goroutine-", id, ": airportSecurityCheckChannel closed\n")
					return // 退出goroutine
				}
				queue1 <- v // 将工作传递给第一个安检环节
			}
		}
	}(id)
}

// max 返回一组整数中的最大值
func max(args ...int) int {
	n := 0
	for _, v := range args {
		if v > n {
			n = v
		}
	}
	return n
}

// 主要特点说明

// 1. 并发设计：使用goroutine实现多个安检通道并行工作

// 2. 流水线模式：每个安检通道内部采用idCheck→bodyCheck→xRayCheck的流水线设计

// 3. 优雅退出：通过关闭quit通道通知goroutine结束工作

// 4. 结果收集：每个安检环节会返回累计耗时，最终取最大值作为总耗时

// 5. 缓冲通道：使用带缓冲的通道提高并发性能

// 这个实现展示了Go语言中goroutine和channel的高级用法，包括并发控制、通道通信和优雅退出等模式。

// 模拟开启了三条通道（newAirportSecurityCheckChannel）​，每条通道创建三个goroutine，
// 分别负责处理idCheck、bodyCheck和xRayCheck，三个goroutine之间通过Go提供的原生channel相连
func main() {
	passengers := 30                 // 乘客数量
	queue := make(chan struct{}, 30) // 主工作通道(缓冲大小为乘客数量)

	// 创建3个并发的安检通道
	newAirportSecurityCheckChannel("channel1", queue)
	newAirportSecurityCheckChannel("channel2", queue)
	newAirportSecurityCheckChannel("channel3", queue)

	time.Sleep(5 * time.Second) // 等待所有goroutine准备就绪

	// 发送30个工作任务(模拟30个乘客)
	for i := 0; i < passengers; i++ {
		queue <- struct{}{}
	}

	time.Sleep(5 * time.Second)
	close(queue)                 // 关闭通道，通知安检通道结束工作
	time.Sleep(10 * time.Second) // 等待所有处理完成
}
