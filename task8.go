//Написать код, который будет выводить
//коды ответов на HTTP-запросы по двум URL
//главная страница Google и главная страница WB.
//Запросы должны осуществляться параллельно.

//Ответ:

package main

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"
)

func main() {

	urls := []string{
		"https://www.google.com",
		"https://www.wildberries.ru",
	}

	// Создаем основной контекст с таймаутом 3 секунды
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	results := make(chan string, len(urls))

	for _, url := range urls {
		go func() {
			statusCode := fetchStatusCode(ctx, url)
			results <- fmt.Sprintf("%s - %d", url, statusCode)
		}()
	}

	for i := 0; i < len(urls); i++ {
		select {
		case result := <-results:
			fmt.Println(result)
		case <-ctx.Done():
			fmt.Println("Timeout exceeded")
			return
		}
	}
}

func fetchStatusCode(ctx context.Context, url string) int {
	client := &http.Client{
		Timeout: 2 * time.Second,
	}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return -1
	}

	resp, err := client.Do(req)
	if err != nil {
		return -1
	}
	defer resp.Body.Close()

	return resp.StatusCode
}

//если много запросов, то как внедрить воркер пул?
С воркер пулом:

type Task struct {
	URL string
	ID  int
}

type Result struct {
	URL        string
	StatusCode int
	Error      error
	ID         int
}

func main() {

	urls := []string{
		"https://www.google.com",
		"https://www.wildberries.ru",
	}

	workerCount := 3
	tasks := make(chan Task, len(urls))
	results := make(chan Result, len(urls))

	var wg sync.WaitGroup
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go worker(ctx, i+1, tasks, results, &wg)
	}

	for i, url := range urls {
		tasks <- Task{URL: url, ID: i}
	}
	close(tasks)

	go func() {
		wg.Wait()
		close(results)
	}()

	for result := range results {
		if result.Error != nil {
			fmt.Printf("[%d] %s - ERROR: %v\n", result.ID, result.URL, result.Error)
		} else {
			fmt.Printf("[%d] %s - %d\n", result.ID, result.URL, result.StatusCode)
		}
	}
}

func worker(ctx context.Context, id int, tasks <-chan Task, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	for task := range tasks {
		select {
		case <-ctx.Done():
			results <- Result{
				URL:   task.URL,
				Error: ctx.Err(),
				ID:    task.ID,
			}
			return
		default:
			statusCode, err := fetchStatusCode(ctx, client, task.URL)
			results <- Result{
				URL:        task.URL,
				StatusCode: statusCode,
				Error:      err,
				ID:         task.ID,
			}
		}
	}
}

func fetchStatusCode(ctx context.Context, client *http.Client, url string) (int, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return -1, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return -1, err
	}
	defer resp.Body.Close()

	return resp.StatusCode, nil
}
