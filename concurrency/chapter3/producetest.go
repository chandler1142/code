package main

import (
	"fmt"
	"sync"
	"time"
)

type good struct {
}

func main() {

	totalCount := 100

	w := sync.WaitGroup{}
	w.Add(totalCount)

	box := make([]interface{}, 0, 1000)
	boxAvailable := sync.NewCond(&sync.Mutex{})
	goodsAvailable := sync.NewCond(&sync.Mutex{})

	consumeCount := 0
	produceCount := 0
	consumer := func(name string) {
		time.Sleep(1 * time.Second)
		goodsAvailable.L.Lock()
		for len(box) == 0 {
			fmt.Printf("box is empty, consumer %s is waiting\n", name)
			goodsAvailable.Wait()
		}
		box = box[1:]
		consumeCount++
		fmt.Printf("consumer box len: %d\n", len(box))
		goodsAvailable.L.Unlock()
		boxAvailable.Signal()
		w.Done()
	}

	producer := func(name string) {
		time.Sleep(1 * time.Second)
		boxAvailable.L.Lock()
		for len(box) == 10 {
			fmt.Printf("box is full, producer %s is waiting\n", name)
			boxAvailable.Wait()
		}
		box = append(box, good{})
		produceCount++
		fmt.Printf("producer box len: %d\n", len(box))
		boxAvailable.L.Unlock()
		goodsAvailable.Signal()
	}

	for i := 0; i < totalCount; i++ {
		go consumer(fmt.Sprintf("c%d", i))
		go producer(fmt.Sprintf("p%d", i))
	}

	w.Wait()

	fmt.Printf("consume count: %d\n", consumeCount)
	fmt.Printf("produce count: %d\n", produceCount)

}
