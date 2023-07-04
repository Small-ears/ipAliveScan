package extension

import (
	"fmt"
	"os"
)

func SaveResult(iplist []string) {
	// 创建文件
	file, err := os.OpenFile("scan_results.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer file.Close()

	// 将每个IP地址写入文件，每个IP地址占据一行，使用换行符
	for _, ip := range iplist {
		_, err = file.WriteString(ip + "\n")
		if err != nil {
			fmt.Printf("Error writing to file: %v\n", err)
			return
		}
	}
}
