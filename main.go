package main

import (
	"flag"
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	url := flag.String("url", "", "URL do serviço a ser testado")
	requests := flag.Int("requests", 0, "número total de requisições")
	concurrency := flag.Int("concurrency", 1, "número de chamadas simultâneas")
	flag.Parse()

	if *url == "" || *requests <= 0 || *concurrency <= 0 {
		fmt.Println("uso: --url=<url> --requests=<n> --concurrency=<n>")
		return
	}

	jobs := make(chan struct{}, *requests)
	for i := 0; i < *requests; i++ {
		jobs <- struct{}{}
	}
	close(jobs)

	var total int64
	var status200 int64
	statusCounts := make(map[int]int64)
	var mu sync.Mutex
	var wg sync.WaitGroup

	client := &http.Client{Timeout: 30 * time.Second}

	start := time.Now()

	for w := 0; w < *concurrency; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range jobs {
				atomic.AddInt64(&total, 1)

				resp, err := client.Get(*url)
				if err != nil {
					mu.Lock()
					statusCounts[0]++
					mu.Unlock()
					continue
				}
				resp.Body.Close()

				if resp.StatusCode == http.StatusOK {
					atomic.AddInt64(&status200, 1)
				}

				mu.Lock()
				statusCounts[resp.StatusCode]++
				mu.Unlock()
			}
		}()
	}

	wg.Wait()
	elapsed := time.Since(start)

	fmt.Println("Relatório de Stress Test")
	fmt.Printf("Tempo total gasto: %s\n", elapsed)
	fmt.Printf("Total de requests realizados: %d\n", total)
	fmt.Printf("Requests com status 200: %d\n", status200)
	fmt.Println("Distribuição de outros status codes:")
	for code, count := range statusCounts {
		if code == http.StatusOK {
			continue
		}
		if code == 0 {
			fmt.Printf("  erro de conexão: %d\n", count)
			continue
		}
		fmt.Printf("  %d: %d\n", code, count)
	}
}
