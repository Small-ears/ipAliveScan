package extension

import (
	"fmt"
	"net"
	"os"
	"sync"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

const (
	icmpProtocol = 1
)

var (
	wg              sync.WaitGroup // 等待组
	lock            sync.Mutex     // 互斥锁
	workerNum       int            //并发数
	icmpTaskCh      chan string
	icmpaliveIpList []string
)

func ScanIcmp(iplist []string, num int) []string {
	workerNum = num
	icmpTaskCh = make(chan string)

	//启动消费者
	for i := 0; i < workerNum; i++ {
		go icmpWorker()
	}

	//提交任务到channel
	wg.Add(len(iplist))
	for _, ip := range iplist {
		icmpTaskCh <- ip
	}
	wg.Wait()
	close(icmpTaskCh)
	return icmpaliveIpList
}

func icmpWorker() {
	for ip := range icmpTaskCh { //channal是可以遍历的
		isTrue := icmpScan(ip)
		if isTrue {
			lock.Lock()
			icmpaliveIpList = append(icmpaliveIpList, ip)
			lock.Unlock()
		}
		wg.Done()
	}
}

func icmpScan(ip string) (isTrue bool) {
	// 创建 ICMP 连接
	conn, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		fmt.Printf("ListenPacket error: %s\n", err)
		return
	}
	defer conn.Close()

	// 设置目标 IP 地址
	dstIPAddr, err := net.ResolveIPAddr("ip4", ip)
	if err != nil {
		fmt.Printf("ResolveIPAddr error: %s\n", err)
		return false
	}

	// 发送多个 ICMP Echo 请求
	for i := 1; i <= 4; i++ {
		// 构建 ICMP 报文
		message := icmp.Message{
			Type: ipv4.ICMPTypeEcho, Code: 0,
			Body: &icmp.Echo{
				ID:   os.Getpid() & 0xffff,
				Seq:  i,
				Data: []byte("Hello"),
			},
		}

		messageBytes, err := message.Marshal(nil)
		if err != nil {
			fmt.Printf("Message marshal error: %s\n", err)
			continue
		}

		// 发送 ICMP 报文

		_, err = conn.WriteTo(messageBytes, dstIPAddr)

		if err != nil {
			fmt.Printf("WriteTo error: %s\n", err)
			continue
		}

		// 等待接收 ICMP 回复
		replyBytes := make([]byte, 1500)
		err = conn.SetReadDeadline(time.Now().Add(time.Second * 4))
		if err != nil {
			fmt.Printf("SetReadDeadline error: %s\n", err)
			continue
		}

		n, _, err := conn.ReadFrom(replyBytes)

		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				fmt.Printf("No reply from %v\n", dstIPAddr)
				continue
			}
			fmt.Printf("ReadFrom error: %s\n", err)
			continue
		}

		// 解析 ICMP 回复
		replyMessage, err := icmp.ParseMessage(icmpProtocol, replyBytes[:n])
		if err != nil {
			fmt.Printf("ParseMessage error: %s\n", err)
			continue
		}

		// 判断回复消息类型为 ICMP Echo Reply
		if replyMessage.Type == ipv4.ICMPTypeEchoReply {
			isTrue = true // 如果收到 ICMP Echo Reply，表示 IP 存活，立即返回 true
			break
		}
	}

	return isTrue
}
