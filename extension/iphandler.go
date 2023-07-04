package extension

import (
	"log"
	"net"

	"github.com/malfunkt/iprange"
)

// GetIpList 接受传参，将IP传参分解为独立的IP
func GetIpList(IP string) ([]string, error) {
	list, err := iprange.ParseList(IP)
	if err != nil {
		log.Printf("error: %s", err)
		return nil, err
	}
	ipList := list.Expand() //[]net.IP类型的切片可以直接传递给另一个切片;

	strIpList := make([]string, len(ipList))

	for i, ip := range ipList {
		strIpList[i] = ip.String()
	}

	localIPs, _ := GetSelfIPs()
	filteredIPs := compareAndRemove(localIPs, strIpList)

	return filteredIPs, nil //返回 IP 列表和 nil 错误
}

// GetSelfIPs 获取本机所有网卡的IP
func GetSelfIPs() ([]string, error) {
	interfaces, err := net.Interfaces() // 获取本机所有网络接口
	if err != nil {
		return nil, err
	}

	selfIPs := make([]string, 0) // 用于存储本地IP地址的切片
	for _, iface := range interfaces {
		addrs, err := iface.Addrs() // 获取与接口关联的地址
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			ipNet, ok := addr.(*net.IPNet)
			if ok && !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
				selfIPs = append(selfIPs, ipNet.IP.String()) // 将非回环且IPv4的地址添加到切片中
			}
		}
	}

	return selfIPs, nil
}

// CompareAndRemove从生成的IP列表中将本机IP剔除
func compareAndRemove(localIPs, ipList []string) []string {
	result := make([]string, 0)
	localIPSet := make(map[string]bool) //map

	// 将本地IP列表存储到set中，方便快速查找
	for _, ip := range localIPs {
		localIPSet[ip] = true
	}

	// 检查ipList中的元素是否存在于本地IP列表中，若不存在则将其添加到结果切片中
	for _, ip := range ipList {
		if !localIPSet[ip] { //localIPSet[ip] 表示通过 ip 在 localIPSet 中进行查找,元素存在true，反之，感叹号取反
			result = append(result, ip)
		}
	}

	return result
}
