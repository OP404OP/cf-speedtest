package main

import (
	"cf_test/internal/generator"
	"cf_test/internal/output"
	"cf_test/internal/tester"
	"cf_test/internal/types"
	"flag"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

const (
	AppName    = "CF Speed Test"
	AppVersion = "v1.0.0"
)

func main() {
	fmt.Printf("%s %s\n\n", AppName, AppVersion)

	mode := flag.String("mode", "tcping", "测试模式: http/tcping")
	port := flag.Int("port", 443, "端口")
	concurrent := flag.Int("c", 10, "并发数")
	inputFile := flag.String("i", "ip.txt", "IP列表文件")
	outputFile := flag.String("o", "result.txt", "结果输出文件")
	useNodes := flag.Bool("n", false, "是否使用节点测试")
	flag.Parse()

	ips, err := generator.GenerateIPs(*inputFile)
	if err != nil {
		fmt.Printf("读取IP文件失败: %v\n", err)
		os.Exit(1)
	}

	output.InitProgress(len(ips))

	var results interface{}
	if *useNodes {
		nodes, err := loadNodes("configs/nodes.yaml")
		if err != nil {
			fmt.Printf("加载节点配置失败: %v\n", err)
			os.Exit(1)
		}
		switch *mode {
		case "tcping":
			results, err = tester.TestTCPingWithNodes(ips, nodes, *concurrent)
			if err != nil {
				fmt.Printf("节点测试失败: %v\n", err)
				os.Exit(1)
			}
		default:
			fmt.Println("节点测试仅支持tcping模式")
			os.Exit(1)
		}
	} else {
		switch *mode {
		case "http":
			results, err = tester.TestHTTP(ips, *port, *concurrent)
			if err != nil {
				fmt.Printf("HTTP测试失败: %v\n", err)
				os.Exit(1)
			}
		case "tcping":
			results, err = tester.TestTCPing(ips, *port, *concurrent)
			if err != nil {
				fmt.Printf("TCPing测试失败: %v\n", err)
				os.Exit(1)
			}
		default:
			fmt.Println("不支持的测试模式")
			os.Exit(1)
		}
	}

	output.DisplayTable(results)

	if err := output.SaveResults(*outputFile, results); err != nil {
		fmt.Printf("保存结果失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("测试完成，结果已保存到", *outputFile)
}

func loadNodes(filename string) ([]*types.Node, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config struct {
		Nodes []*types.Node `yaml:"nodes"`
	}
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return config.Nodes, nil
}
