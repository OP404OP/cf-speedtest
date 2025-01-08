package types

// HTTPResult 存储HTTP测试结果
type HTTPResult struct {
	IP         string
	MinLatency int64
	MaxLatency int64
	AvgLatency int64
	PacketLoss float64
	StatusCode int
}

// TCPingResult 存储TCPing测试结果
type TCPingResult struct {
	IP         string
	MinLatency int64
	MaxLatency int64
	AvgLatency int64
	PacketLoss float64
}

// Node 测试节点信息
type Node struct {
	Name     string // 如：北京电信、上海联通
	IP       string // 节点IP地址
	Location string // 地理位置
	ISP      string // 运营商
	Port     int    // 端口号，默认80
}

// TestResult 测试结果
type TestResult struct {
	Node       *Node
	IP         string
	MinLatency int64
	MaxLatency int64
	AvgLatency int64
	PacketLoss float64
}
