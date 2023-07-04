package extension

import (
	"fmt"
	"net"
	"sync"
	"time"
)

var (
	workerCount int
	taskCh      chan string
	aliveIpList []string
	aliveIpLock sync.Mutex //互斥锁
	wgb         sync.WaitGroup
)

func ScanPort(iplist []string, num int) []string {
	aliveIpList = make([]string, 0)
	workerCount = num //最大并发数
	taskCh = make(chan string)

	// 启动消费者
	for i := 0; i < workerCount; i++ {
		go worker()
	}

	// 提交任务到任务通道
	wgb.Add(len(iplist))
	for _, ip := range iplist {
		taskCh <- ip
	}

	wgb.Wait() // 等待所有任务完成

	close(taskCh) // 关闭任务通道等待消费者处理完剩余任务

	return aliveIpList
}

func worker() {
	for ip := range taskCh {
		isTrue := tcpPort(ip)
		if isTrue {
			aliveIpLock.Lock()
			aliveIpList = append(aliveIpList, ip)
			aliveIpLock.Unlock()
		}
		wgb.Done()
	}
}

func tcpPort(ip string) bool {
	tcpPorts := []int{
		22,  // SSH
		21,  // FTP
		23,  // Telnet
		25,  // SMTP
		53,  //DNS
		80,  // HTTP
		110, // POP3
		135,
		139,
		143,  // IMAP
		443,  // HTTPS
		3306, // MySQL
		3389, // RDP
	}

	for _, port := range tcpPorts {
		conn, err := net.DialTimeout("tcp", fmt.Sprintf("%v:%v", ip, port), time.Millisecond*600)
		if err != nil {
			continue
		}
		defer conn.Close()
		return true
	}

	return false
}
