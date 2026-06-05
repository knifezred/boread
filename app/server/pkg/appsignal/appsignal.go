// Package appsignal 提供进程级信号协调机制，用于服务内部触发热重启
package appsignal

// restartCh 带缓冲 1 的通道，避免发送方阻塞
var restartCh = make(chan struct{}, 1)

// RequestRestart 发送重启信号（非阻塞），配置保存成功后调用
func RequestRestart() {
	select {
	case restartCh <- struct{}{}:
	default:
	}
}

// WaitRestart 阻塞等待重启信号，main 函数中调用
func WaitRestart() {
	<-restartCh
}
