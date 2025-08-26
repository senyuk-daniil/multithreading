//Fan-In (объединение данных из нескольких каналов)
//Задача: Напишите функцию, которая объединяет два входных канала в один выходной.
//// Объединяем разные каналы в один канал
//func Merge(cs ...<-chan int) <-chan int

package main

import (
	"fmt"
	"sync"
)

func Merge(cs ...<-chan int) <-chan int {
	var wg sync.WaitGroup

	out := make(chan int)

	output := func(c <-chan int) {
		defer wg.Done()
		for n := range c {
			out <- n
		}
	}

	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	go func() {
		defer close(out)
		wg.Wait()
	}()

	return out
}

func main() {
	ch1 := make(chan int)
	go func() {
		defer close(ch1)
		ch1 <- 1
	}()
	ch2 := make(chan int)
	go func() {
		defer close(ch2)
		ch2 <- 2
	}()

	for v := range Merge(ch1, ch2) {
		fmt.Println(v)
	}
}
