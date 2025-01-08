package output

import (
	"cf_test/internal/types"
	"fmt"
	"os"
	"sort"

	"github.com/olekukonko/tablewriter"
	"github.com/schollz/progressbar/v3"
)

var bar *progressbar.ProgressBar

// InitProgress 初始化进度条
func InitProgress(total int) {
	bar = progressbar.Default(int64(total))
}

// UpdateProgress 更新进度
func UpdateProgress() {
	if bar != nil {
		bar.Add(1)
	}
}

// DisplayTable 在控制台显示结果表格
func DisplayTable(results interface{}) {
	switch v := results.(type) {
	case []*types.TestResult:
		// 按节点分组显示结果
		nodeMap := make(map[string][]*types.TestResult)
		for _, r := range v {
			if r.AvgLatency > 0 {
				nodeMap[r.Node.Name] = append(nodeMap[r.Node.Name], r)
			}
		}

		for nodeName, nodeResults := range nodeMap {
			fmt.Printf("\n节点: %s\n", nodeName)
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"IP", "最低延迟(ms)", "最高延迟(ms)", "平均延迟(ms)", "丢包率(%)"})

			sort.Slice(nodeResults, func(i, j int) bool {
				return nodeResults[i].AvgLatency < nodeResults[j].AvgLatency
			})

			for _, r := range nodeResults {
				table.Append([]string{
					r.IP,
					fmt.Sprintf("%d", r.MinLatency),
					fmt.Sprintf("%d", r.MaxLatency),
					fmt.Sprintf("%d", r.AvgLatency),
					fmt.Sprintf("%.1f", r.PacketLoss),
				})
			}
			table.Render()
		}
	case []*types.HTTPResult:
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"IP", "最低延迟(ms)", "最高延迟(ms)", "平均延迟(ms)", "丢包率(%)", "状态码"})

		// 过滤并排序
		var validResults []*types.HTTPResult
		for _, r := range v {
			if r.AvgLatency > 0 {
				validResults = append(validResults, r)
			}
		}
		sort.Slice(validResults, func(i, j int) bool {
			return validResults[i].AvgLatency < validResults[j].AvgLatency
		})
		for _, r := range validResults {
			table.Append([]string{
				r.IP,
				fmt.Sprintf("%d", r.MinLatency),
				fmt.Sprintf("%d", r.MaxLatency),
				fmt.Sprintf("%d", r.AvgLatency),
				fmt.Sprintf("%.1f", r.PacketLoss),
				fmt.Sprintf("%d", r.StatusCode),
			})
		}
		table.Render()
	case []*types.TCPingResult:
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"IP", "最低延迟(ms)", "最高延迟(ms)", "平均延迟(ms)", "丢包率(%)"})

		// 过滤并排序
		var validResults []*types.TCPingResult
		for _, r := range v {
			if r.AvgLatency > 0 {
				validResults = append(validResults, r)
			}
		}
		sort.Slice(validResults, func(i, j int) bool {
			return validResults[i].AvgLatency < validResults[j].AvgLatency
		})
		for _, r := range validResults {
			table.Append([]string{
				r.IP,
				fmt.Sprintf("%d", r.MinLatency),
				fmt.Sprintf("%d", r.MaxLatency),
				fmt.Sprintf("%d", r.AvgLatency),
				fmt.Sprintf("%.1f", r.PacketLoss),
			})
		}
		table.Render()
	}
}

// SaveResults 保存结果到文件
func SaveResults(filename string, results interface{}) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	switch v := results.(type) {
	case []*types.TestResult:
		// 按节点分组保存
		nodeMap := make(map[string][]*types.TestResult)
		for _, r := range v {
			if r.AvgLatency > 0 {
				nodeMap[r.Node.Name] = append(nodeMap[r.Node.Name], r)
			}
		}

		for nodeName, nodeResults := range nodeMap {
			fmt.Fprintf(file, "# 节点: %s\n", nodeName)
			sort.Slice(nodeResults, func(i, j int) bool {
				return nodeResults[i].AvgLatency < nodeResults[j].AvgLatency
			})
			for _, r := range nodeResults {
				fmt.Fprintln(file, r.IP)
			}
			fmt.Fprintln(file)
		}
	case []*types.HTTPResult:
		// 过滤并排序
		var validResults []*types.HTTPResult
		for _, r := range v {
			if r.AvgLatency > 0 { // 只保留有效结果
				validResults = append(validResults, r)
			}
		}
		sort.Slice(validResults, func(i, j int) bool {
			return validResults[i].AvgLatency < validResults[j].AvgLatency
		})
		for _, r := range validResults {
			fmt.Fprintln(file, r.IP)
		}
	case []*types.TCPingResult:
		// 过滤并排序
		var validResults []*types.TCPingResult
		for _, r := range v {
			if r.AvgLatency > 0 { // 只保留有效结果
				validResults = append(validResults, r)
			}
		}
		sort.Slice(validResults, func(i, j int) bool {
			return validResults[i].AvgLatency < validResults[j].AvgLatency
		})
		for _, r := range validResults {
			fmt.Fprintln(file, r.IP)
		}
	}

	return nil
}
