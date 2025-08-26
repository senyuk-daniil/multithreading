package main

import (
	"fmt"
	"sync"
)

func main() {
	var wa sync.WaitGroup
	wa.Add(1)
	go func() {
		defer wa.Done()
		fmt.Println("Hello from goroutine!")
	}()
	wa.Wait()
}

//Какие способы есть ещё кроме waitGroup, чтобы дождаться выполнения горутины? Приведи хотя бы 2 примера.

//*) Можно использовать time.Sleep
//*) Можно использовать каналы
