//Реализуй функцию `StartBatchProcessor(ctx context.Context, input <-chan int)`, которая:
//
//- Собирает числа из канала `input` в батчи по максимум 5 элементов.

//- Если в течение 2 секунд батч не собран — обрабатывает то, что есть.

//- Обработка батча — это просто `fmt.Println("Processed batch:", batch)`.

//- Выход из функции должен происходить при отмене контекста (`ctx.Done()`).
//
// **Дополнительно**:
//- Отмена должна происходить либо чере`з context.WithTimeout`, либо вручную через `cancel()` — попробовать оба варианта
//Начальный код с вызовом(доработать)

package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func StartBatchProcessor(ctx context.Context, input <-chan int) {
	const batchSize = 5

	batch := make([]int, 0, batchSize)
	timer := time.NewTimer(time.Second * 2)
	defer timer.Stop()

	for {
		select {
		case num, ok := <-input:
			if !ok {
				if len(batch) > 0 {
					fmt.Println("Processed batch:", batch)
				}
				return
			}
			batch = append(batch, num)

			if len(batch) == batchSize {
				fmt.Println("Processed batch:", batch)
				batch = make([]int, 0, batchSize)

				if !timer.Stop() {
					<-timer.C
				}
				timer.Reset(time.Second * 2)
			}

		case <-ctx.Done():
			if len(batch) > 0 {
				fmt.Println("Processed batch:", batch)
			}
			return

		case <-timer.C:
			if len(batch) > 0 {
				fmt.Println("Processed batch:", batch)
				batch = make([]int, 0, batchSize)
			}
			timer.Reset(2 * time.Second)
		}
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	input := make(chan int)

	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		for i := 0; i < 12; i++ {
			input <- i
		}
		time.Sleep(10 * time.Second)
		cancel()
		close(input)
	}()

	go func() {
		defer wg.Done()
		StartBatchProcessor(ctx, input)
	}()

	wg.Wait()
	fmt.Println("Main: processing stopped")
}
