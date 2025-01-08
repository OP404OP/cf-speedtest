package tester

import (
	"cf_test/internal/output"
	"cf_test/internal/types"
	"fmt"
	"math"
	"net"
	"sync"
	"time"
)

func tcping(nodeIP string, nodePort int, targetIP string, targetPort int) (latency int64, err error) {
	start := time.Now()
	// 直接使用IP地址和端口号，没有DNS解析
	nodeConn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", nodeIP, nodePort), time.Second*2)
	if err != nil {
		return 0, fmt.Errorf("连接节点失败: %v", err)
	}
	defer nodeConn.Close()

	// 直接使用IP地址和端口号，没有DNS解析
	targetConn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", targetIP, targetPort), time.Second*2)
	if err != nil {
		return 0, fmt.Errorf("通过节点连接目标失败: %v", err)
	}
	defer targetConn.Close()

	return time.Since(start).Milliseconds(), nil
}

// TestTCPing 测试TCPing
func TestTCPing(ips []string, port int, concurrent int) ([]*types.TCPingResult, error) {
	results := make([]*types.TCPingResult, len(ips))
	sem := make(chan bool, concurrent)
	var wg sync.WaitGroup

	for i, ip := range ips {
		wg.Add(1)
		sem <- true

		go func(index int, ip string) {
			defer wg.Done()
			defer func() { <-sem }()

			result := &types.TCPingResult{
				IP:         ip,
				MinLatency: math.MaxInt64,
			}

			var totalLatency int64
			var successCount int
			const testCount = 4

			for i := 0; i < testCount; i++ {
				if latency, err := tcping("", 0, ip, port); err == nil {
					successCount++
					totalLatency += latency
					if latency < result.MinLatency {
						result.MinLatency = latency
					}
					if latency > result.MaxLatency {
						result.MaxLatency = latency
					}
				} else {
					fmt.Printf("IP %s TCP连接失败: %v\n", ip, err)
				}
				time.Sleep(time.Millisecond * 200)
			}

			if successCount > 0 {
				result.AvgLatency = totalLatency / int64(successCount)
			} else {
				result.MinLatency = 0
				result.MaxLatency = 0
				result.AvgLatency = 0
			}
			result.PacketLoss = float64(testCount-successCount) / float64(testCount) * 100

			results[index] = result
			output.UpdateProgress()
		}(i, ip)
	}

	wg.Wait()
	return results, nil
}

// TestTCPingWithNodes 使用节点进行TCPing测试
func TestTCPingWithNodes(ips []string, nodes []*types.Node, concurrent int) ([]*types.TestResult, error) {
	results := make([]*types.TestResult, 0, len(ips)*len(nodes))
	sem := make(chan bool, concurrent)
	var wg sync.WaitGroup

	for _, node := range nodes {
		for _, ip := range ips {
			wg.Add(1)
			sem <- true

			go func(node *types.Node, ip string) {
				defer wg.Done()
				defer func() { <-sem }()

				result := &types.TestResult{
					Node:       node,
					IP:         ip,
					MinLatency: math.MaxInt64,
				}

				var totalLatency int64
				var successCount int
				const testCount = 4

				for i := 0; i < testCount; i++ {
					if latency, err := tcping(node.IP, node.Port, ip, 443); err == nil {
						successCount++
						totalLatency += latency
						if latency < result.MinLatency {
							result.MinLatency = latency
						}
						if latency > result.MaxLatency {
							result.MaxLatency = latency
						}
					} else {
						fmt.Printf("节点 %s 测试 IP %s: 第 %d 次重试失败: %v\n", node.Name, ip, i+1, err)
					}
					time.Sleep(time.Millisecond * 200)
				}

				if successCount > 0 {
					result.AvgLatency = totalLatency / int64(successCount)
				} else {
					result.MinLatency = 0
					result.MaxLatency = 0
					result.AvgLatency = 0
				}
				result.PacketLoss = float64(testCount-successCount) / float64(testCount) * 100

				results = append(results, result)
				output.UpdateProgress()
			}(node, ip)
		}
	}

	wg.Wait()
	return results, nil
}
