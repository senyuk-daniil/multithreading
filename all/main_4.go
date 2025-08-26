// Задача: Напишите программу, где 10 горутин инкрементируют один счётчик, защищая его sync.Mutex.
package main

import (
	"fmt"
	"sync"
)

func main() {
	var mu sync.Mutex
	var wa sync.WaitGroup
	count := 0
	for i := 0; i < 10; i++ {
		wa.Add(1)
		go func(i int) {
			defer wa.Done()
			mu.Lock()
			count++
			mu.Unlock()
		}(i)
	}
	wa.Wait()
	fmt.Println("count:", count)
}

//1. Что если не обложить мютексом? Воспроизвести race condition.
//* Некоторые горутины могут не внести изменения, особенно заметно это будет, если count большой, например 100000.
//Mutex гарантирует, что доступ к счётчику будет у каждой горутины
