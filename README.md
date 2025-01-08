# CF Speed Test

Cloudflare CDN 节点测速工具，支持多节点测试。

> 基于 Cursor 的 AI 助手实现。可批量测试cloudflare IP地址的最低延迟、最高延迟、丢包率、平均延迟、支持自定义节点测试、默认按平均延迟排序

> 测试结果准确性有待检验，咱也不知道准不准

## 功能特点

- 支持 HTTP/TCPING 测试模式
- 支持多节点并发测试
- 支持 CIDR 和 IP 范围格式
- 跨平台支持 (Windows/Linux/MacOS)
- 默认按平均延迟排序(从低到高)

## 下载安装

从 [Releases](https://github.com/OP404OP/cf-speedtest/releases) 页面下载对应平台的二进制文件：

- Windows: `cfspeedtest_v1.0.0_windows_amd64.exe`
- Linux: `cfspeedtest_v1.0.0_linux_amd64`
- MacOS: `cfspeedtest_v1.0.0_darwin_amd64`

## 使用方法

### 命令行参数
```bash
选项：
-mode string 测试模式: http/tcping (默认 "tcping")

-port int 端口号 (默认 443)

-c int 并发数 (默认 10)

-i string IP列表文件 (默认 "ip.txt")

-o string 结果输出文件 (默认 "result.txt")

-n 启用节点测试
```

### 使用示例
```bash
HTTP测试
./cfspeedtest -mode http -port 443 -i ip.txt -o result.txt

TCPing测试
./cfspeedtest -mode tcping -port 443 -i ip.txt -o result.txt

节点测试
./cfspeedtest -mode tcping -n -i ip.txt -o result.txt
```

## 配置文件

### 1. IP列表 (ip.txt)
```bash
CIDR格式
1.1.1.1/24

IP范围格式
1.0.0.1-1.0.0.255
```
### 2. 节点配置 (configs/nodes.yaml)
```yaml
nodes:
- name: "北京电信"
  ip: "1.2.3.4"
  location: "北京"
  isp: "电信"
  port: 80
```

## 开源许可

本项目采用 [MIT License](LICENSE) 开源。
