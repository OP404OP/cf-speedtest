package tester

import (
	"cf_test/internal/output"
	"cf_test/internal/types"
	"fmt"
	"math"
	"net"
	"sync"
	"time"

	"crypto/tls"

	"github.com/valyala/fasthttp"
)

// TestHTTP 测试HTTP
func TestHTTP(ips []string, port int, concurrent int) ([]*types.HTTPResult, error) {
	results := make([]*types.HTTPResult, len(ips))
	sem := make(chan bool, concurrent)
	var wg sync.WaitGroup

	for i, ip := range ips {
		wg.Add(1)
		sem <- true

		go func(index int, ip string) {
			defer wg.Done()
			defer func() { <-sem }()

			result := &types.HTTPResult{
				IP:         ip,
				MinLatency: math.MaxInt64,
			}

			client := &fasthttp.Client{
				MaxConnsPerHost:     1,
				ReadTimeout:         time.Second * 15,
				WriteTimeout:        time.Second * 15,
				DialDualStack:       true,
				MaxIdleConnDuration: time.Second * 10,
				TLSConfig: &tls.Config{
					InsecureSkipVerify: true,
					ServerName:         "cloudflare.com",
					MinVersion:         tls.VersionTLS10,
					MaxVersion:         tls.VersionTLS13,
					CipherSuites: []uint16{
						tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
						tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
						tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
						tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
						tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
						tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
					},
				},
				NoDefaultUserAgentHeader: true,
				DisablePathNormalizing:   true,
			}

			url := fmt.Sprintf("https://%s:%d", ip, port)
			req := fasthttp.AcquireRequest()
			resp := fasthttp.AcquireResponse()
			defer fasthttp.ReleaseRequest(req)
			defer fasthttp.ReleaseResponse(resp)

			req.SetRequestURI(url)
			req.Header.SetHost("cloudflare.com")
			req.URI().SetHost("cloudflare.com")
			req.URI().SetPath("/cdn-cgi/trace")
			req.Header.SetUserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/120.0.0.0 Safari/537.36")
			req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
			req.Header.Set("Accept-Language", "en-US,en;q=0.5")
			req.Header.Set("Accept-Encoding", "gzip, deflate, br")
			req.Header.Set("Connection", "keep-alive")
			req.Header.Set("Upgrade-Insecure-Requests", "1")

			var totalLatency int64
			var successCount int
			const testCount = 4

			for i := 0; i < testCount; i++ {
				// 1. 测量TCP连接时间
				tcpStart := time.Now()
				conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", ip, port), time.Second*5)
				tcpLatency := time.Since(tcpStart).Milliseconds()

				if err != nil {
					fmt.Printf("IP %s TCP连接失败: %v\n", ip, err)
					continue
				}
				conn.Close()

				// 2. 测量HTTPS总时间
				start := time.Now()
				err = client.Do(req, resp)
				httpsLatency := time.Since(start).Milliseconds()

				if err == nil {
					successCount++
					// 使用TCP延迟而不是HTTPS总延迟
					if tcpLatency < result.MinLatency {
						result.MinLatency = tcpLatency
					}
					if tcpLatency > result.MaxLatency {
						result.MaxLatency = tcpLatency
					}
					totalLatency += tcpLatency
					result.StatusCode = resp.StatusCode()
					fmt.Printf("IP %s TCP延迟: %dms, HTTPS延迟: %dms\n", ip, tcpLatency, httpsLatency)
				} else {
					fmt.Printf("IP %s HTTPS请求失败: %v\n", ip, err)
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
