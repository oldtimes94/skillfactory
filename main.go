package main

import (
	"container/ring"
	"fmt"
	"time"
)

func ReadInts() chan int {
	out := make(chan int)
	var value int

	go func() {
		for {
			_, err := fmt.Scanln(&value)
			if err != nil {
				fmt.Println(err)
				continue
			}
			out <- value
		}
	}()
	return out
}

func NegativeNumberFilter(ints chan int) chan int {
	out := make(chan int)

	go func() {
		for value := range ints {
			if value < 0 {
				fmt.Println("Value < 0, filtered:", value)
				continue
			}
			out <- value
		}
	}()
	return out
}

func MultipleOfThreeFilter(ints chan int) chan int {
	out := make(chan int)

	go func() {
		for value := range ints {
			if value%3 != 0 || value == 0 {
				fmt.Println("Value % 3 != 0, filtered:", value)
				continue
			}
			out <- value
		}
	}()
	return out
}

func main() {

	const (
		RingBufferSize     = 5
		TruncateBufferTime = 15 * time.Second
	)

	ReadingChannel := ReadInts()
	FilterNegativeChannel := NegativeNumberFilter(ReadingChannel)
	MultpleOfThreeChannel := MultipleOfThreeFilter(FilterNegativeChannel)

	RingBuffer := ring.New(RingBufferSize)
	TruncateBufferTicker := time.NewTicker(TruncateBufferTime)

	go func() {
		for {
			select {
			case <-TruncateBufferTicker.C:
				RingBuffer = ring.New(RingBufferSize)
				fmt.Println("It's time to truncate buffer")
			}
		}
	}()

	for value := range MultpleOfThreeChannel {
		RingBuffer.Value = value
		RingBuffer = RingBuffer.Next()
		fmt.Println("Received Value", value)

		n := RingBuffer.Len()
		fmt.Printf("Current buffer: ")
		for j := 0; j < n; j++ {
			fmt.Printf("%v ", RingBuffer.Value)
			RingBuffer = RingBuffer.Next()
		}
		fmt.Println("")

	}

}
