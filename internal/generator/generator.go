package generator

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"strings"
)

// 将IP字符串转换为uint32
func ip2uint(ip net.IP) uint32 {
	ip = ip.To4()
	return uint32(ip[0])<<24 | uint32(ip[1])<<16 | uint32(ip[2])<<8 | uint32(ip[3])
}

// 将uint32转换为IP字符串
func uint2ip(nn uint32) string {
	return fmt.Sprintf("%d.%d.%d.%d", byte(nn>>24), byte(nn>>16), byte(nn>>8), byte(nn))
}

// 处理CIDR格式
func generateFromCIDR(cidr string) ([]string, error) {
	_, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}

	var ips []string
	// 从网络地址开始
	start := ip2uint(ipnet.IP)
	// 计算最后一个IP
	mask := binary.BigEndian.Uint32(ipnet.Mask)
	end := start | (^mask)

	// 包含所有IP（从.0到.255）
	for ip := start; ip <= end; ip++ {
		ips = append(ips, uint2ip(ip))
	}
	return ips, nil
}

// 处理IP范围格式 (1.1.1.1-1.1.1.255)
func generateFromRange(ipRange string) ([]string, error) {
	parts := strings.Split(ipRange, "-")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid IP range format: %s", ipRange)
	}

	startIP := net.ParseIP(strings.TrimSpace(parts[0])).To4()
	endIP := net.ParseIP(strings.TrimSpace(parts[1])).To4()

	if startIP == nil || endIP == nil {
		return nil, fmt.Errorf("invalid IP address in range: %s", ipRange)
	}

	start := ip2uint(startIP)
	end := ip2uint(endIP)

	if end < start {
		return nil, fmt.Errorf("end IP is smaller than start IP: %s", ipRange)
	}

	var ips []string
	for ip := start; ip <= end; ip++ {
		ips = append(ips, uint2ip(ip))
	}
	return ips, nil
}

// GenerateIPs 从文件生成IP列表
func GenerateIPs(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var allIPs []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		fmt.Printf("处理IP行: %s\n", line)
		var ips []string
		var err error

		if strings.Contains(line, "/") {
			ips, err = generateFromCIDR(line)
			fmt.Printf("CIDR %s 生成了 %d 个IP\n", line, len(ips))
		} else if strings.Contains(line, "-") {
			ips, err = generateFromRange(line)
			fmt.Printf("Range %s 生成了 %d 个IP\n", line, len(ips))
		} else {
			ips = []string{line}
		}

		if err != nil {
			return nil, fmt.Errorf("line '%s': %v", line, err)
		}
		allIPs = append(allIPs, ips...)
	}

	return allIPs, scanner.Err()
}
