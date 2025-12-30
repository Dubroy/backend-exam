package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Employee struct {
	ID           int
	ProcessCount int
	mu           sync.Mutex
}

type Item1 struct{}

type Item2 struct{}

type Item3 struct{}

type Item interface {
	// Process 這是一個耗時操作
	Process()
}

// Process 實現 Item1 的處理，處理時間為 100ms
func (i Item1) Process() {
	time.Sleep(100 * time.Millisecond)
}

// Process 實現 Item2 的處理，處理時間為 200ms
func (i Item2) Process() {
	time.Sleep(200 * time.Millisecond)
}

// Process 實現 Item3 的處理，處理時間為 300ms
func (i Item3) Process() {
	time.Sleep(300 * time.Millisecond)
}

func main() {
	// 初始化隨機數生成器
	rand.Seed(time.Now().UnixNano())

	// 創建 30 個物品（每種 10 個）
	items := make([]Item, 0, 30)
	for i := 0; i < 10; i++ {
		items = append(items, Item1{})
		items = append(items, Item2{})
		items = append(items, Item3{})
	}

	// 隨機打亂物品順序
	rand.Shuffle(len(items), func(i, j int) {
		items[i], items[j] = items[j], items[i]
	})

	// 創建 5 個員工
	employees := make([]*Employee, 5)
	for i := 0; i < 5; i++ {
		employees[i] = &Employee{ID: i + 1}
	}

	// 創建工作 channel
	jobChan := make(chan Item, len(items))
	doneChan := make(chan *Employee, len(items))

	// 記錄開始時間
	startTime := time.Now()
	fmt.Printf("開始處理時間: %s\n", startTime.Format("2006-01-02 15:04:05.000"))

	// 啟動 5 個 worker
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go worker(employees[i], jobChan, doneChan, &wg)
	}

	// 中央調度器：將物品發送到 worker pool
	go func() {
		for _, item := range items {
			jobChan <- item
		}
		close(jobChan)
	}()

	// 等待所有 worker 完成
	wg.Wait()
	close(doneChan)

	// 記錄結束時間
	endTime := time.Now()
	totalTime := endTime.Sub(startTime)
	fmt.Printf("結束處理時間: %s\n", endTime.Format("2006-01-02 15:04:05.000"))
	fmt.Printf("總處理時間: %v\n", totalTime)

	// 統計每個員工處理的物品數量
	fmt.Println("\n各員工處理統計:")
	for _, emp := range employees {
		fmt.Printf("員工 %d: 處理了 %d 個物品\n", emp.ID, emp.ProcessCount)
	}
}

// worker 處理物品的函數
func worker(emp *Employee, jobChan <-chan Item, doneChan chan<- *Employee, wg *sync.WaitGroup) {
	defer wg.Done()

	for item := range jobChan {
		// 開始處理
		fmt.Printf("員工 %d 開始處理物品 (類型: %T)\n", emp.ID, item)
		start := time.Now()

		// 處理物品
		item.Process()

		// 結束處理
		duration := time.Since(start)
		fmt.Printf("員工 %d 結束處理物品 (類型: %T, 耗時: %v)\n", emp.ID, item, duration)

		// 更新處理計數
		emp.mu.Lock()
		emp.ProcessCount++
		emp.mu.Unlock()

		doneChan <- emp
	}
}
