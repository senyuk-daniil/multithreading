//## 8. Практическая задача:
//Конкурентный HTTP-запрос к списку URL-адресов с использованием горутин и синхронизации
//
//> **Задача:**
//**Напишите функцию FetchURLs(urls []string) map[string]string, которая:**
//
//*Принимает слайс URL-адресов.
//Конкурентно делает HTTP-запросы к каждому URL.
//Собирает результаты (код ответа и часть тела) в map[string]string, где:
//ключ — URL
//значение — содержимое ответа (ограниченное, например, 100 символами)
//Использует sync.WaitGroup и sync.Mutex для защиты записи в map.
//В случае ошибки записывает "error" как значение.*
//>

package main

import (
	"fmt"
	"net/http"
	"sync"
)

func FetchURLs(urls []string) map[string]string {
	var wg sync.WaitGroup
	var mu sync.Mutex

	result := make(map[string]string)

	for _, url := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			resp, err := http.Get(url)
			if err != nil {
				mu.Lock()
				result[url] = err.Error()
				mu.Unlock()
				return
			}
			defer resp.Body.Close()

			buf := make([]byte, 100)
			n, _ := resp.Body.Read(buf)
			mu.Lock()
			result[url] = string(buf[:n])
			mu.Unlock()
		}(url)
	}
	wg.Wait()
	return result
}

func main() {

	urls := []string{
		"https://www.example1.com",
		"https://www.example2.com",
		"https://www.example3.com",
	}

	results := FetchURLs(urls)

	for url, result := range results {
		fmt.Println(url, result)
	}
}
