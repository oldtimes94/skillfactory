package main

import (
	"container/ring"
	"fmt"
	"go.uber.org/zap"
	"time"
)

var logger *zap.Logger

func init() {
	var err error
	logger, err = zap.NewProduction()
	if err != nil {
		panic(err)
	}

	defer logger.Sync()

}

func ReadInts() chan int {

	defer logger.Info("ReadInts channel inited")

	out := make(chan int)
	var value int

	go func() {
		for {
			_, err := fmt.Scanln(&value)
			if err != nil {
				logger.Error(err.Error(), zap.String("func", "ReadInts"))
				continue
			}
			logger.Info("[ReadInts] Sending value to out chan", zap.Int("value", value))
			out <- value
		}
	}()
	return out
}

func NegativeNumberFilter(ints chan int) chan int {

	defer logger.Info("NegativeNumberFilter channel inited")

	out := make(chan int)

	go func() {
		for value := range ints {
			if value < 0 {
				logger.Info("Value < 0, filtered: ", zap.Int("value", value))
				continue
			}
			out <- value
		}
	}()
	return out
}

func MultipleOfThreeFilter(ints chan int) chan int {
	defer logger.Info("MultipleOfThreeFilter channel inited")

	out := make(chan int)

	go func() {
		for value := range ints {
			if value%3 != 0 || value == 0 {
				logger.Info("Value % 3 != 0 filtered", zap.Int("value", value))
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
				logger.Info("Truncating bugger (timer)")
			}
		}
	}()

	for value := range MultpleOfThreeChannel {
		RingBuffer.Value = value
		RingBuffer = RingBuffer.Next()
		logger.Info("Received Value", zap.Int("value", value))

		n := RingBuffer.Len()
		logger.Info("Current buffer")
		for j := 0; j < n; j++ {
			logger.Info("Buffer", zap.Int("index", j), zap.Any("value", RingBuffer.Value))
			RingBuffer = RingBuffer.Next()
		}

	}

}
