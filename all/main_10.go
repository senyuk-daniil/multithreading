//Fan-Out (разделение работы между несколькими воркерами)
//Задача: Создайте функцию, которая принимает канал с задачами и распределяет их между N горутинами.
//// Разделить канал на n каналов,
//// которые получают сообщения в циклическом порядке.
//func Split(ch <-chan int, n int) []<-chan int

package main

import (
	"fmt"
	"sync"
)

func Split(ch <-chan int, n int) []<-chan int {
	out := make([]chan int, n)
	for i := range out {
		out[i] = make(chan int)
	}

	go func() {
		defer func() {
			for _, c := range out {
				close(c)
			}
		}()

		i := 0
		for val := range ch {
			out[i] <- val
			i = (i + 1) % n
		}
	}()

	result := make([]<-chan int, n)
	for i := range out {
		result[i] = out[i]
	}
	return result
}

func main() {
	tasks := make(chan int)
	nWorkers := 3

	splits := Split(tasks, nWorkers)

	var wg sync.WaitGroup
	wg.Add(nWorkers)

	for i, ch := range splits {
		go func(id int, ch <-chan int) {
			defer wg.Done()
			for val := range ch {
				fmt.Printf("Worker %d received task %d\n", id, val)
			}
		}(i, ch)
	}

	go func() {
		for i := 1; i <= 10; i++ {
			tasks <- i
		}
		close(tasks)
	}()

	wg.Wait()
}
