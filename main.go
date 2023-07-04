package main

import (
	"flag"
	"fmt"

	"gocode.com/IPScan/extension"
)

func main() {
	var target string
	var save bool
	var tcp bool
	var icmp bool
	var goNum int

	flag.StringVar(&target, "u", "", "example: 127.0.0.1,127.0.0.1/24,127.0.0.*,127.0.0.1-10")
	flag.BoolVar(&tcp, "t", false, "TCP")
	flag.BoolVar(&icmp, "p", false, "ICMPping")
	flag.IntVar(&goNum, "c", 50, "并发数量,默认为50")
	flag.BoolVar(&save, "s", false, "save result")
	flag.Parse()

	if !flag.Parsed() || (target == "" && !tcp && !icmp) { //无参数输入的时候输出帮助信息
		flag.Usage()
		return
	}

	iplist, err := extension.GetIpList(target) //获得任务列表
	if err != nil {
		fmt.Println(err)
		return
	}

	if tcp && !icmp { //检查变量 tcp 的值为 true，且变量 icmp 的值为 false 的情况
		ips := extension.ScanPort(iplist, goNum)
		if save { //是否保存结果
			extension.SaveResult(ips)
		} else {
			for _, ip := range ips {
				fmt.Println(ip)
			}
		}
	} else if !tcp && icmp {
		aliveIpList := extension.ScanIcmp(iplist, goNum)
		if save {
			extension.SaveResult(aliveIpList)
		} else {
			for _, ip := range aliveIpList {
				fmt.Println(ip)
			}
		}
	} else if tcp && icmp {
		aliveIpList := extension.ScanIcmp(iplist, goNum)
		ips := extension.ScanPort(iplist, goNum)
		newResult := compareAndRemove(aliveIpList, ips)

		if save {
			extension.SaveResult(newResult)
		} else {
			for _, ip := range newResult {
				fmt.Println(ip)
			}
		}

	} else {
		fmt.Println("Usage:")
		flag.Usage()
	}
}

// compareAndRemove 切片组合去重
func compareAndRemove(ips, ipList []string) []string {
	uniqueMap := make(map[string]bool)

	for _, ip := range ips {
		uniqueMap[ip] = true
	}

	for _, ip := range ipList {
		uniqueMap[ip] = true
	}

	uniqueSlice := make([]string, 0, len(uniqueMap))
	for ip := range uniqueMap {
		uniqueSlice = append(uniqueSlice, ip)
	}

	return uniqueSlice
}
