package main

import (
	"fmt"
	"time"
)

func convertChannels(channels ...chan interface{}) <-chan interface{} {
	outputChan := make(chan interface{})

	for _, channel := range channels {
		go func(listen, done chan interface{}) {
			done <- <-listen
		}(channel, outputChan)
	}

	return outputChan
}

func main() {

	signal := func(duration time.Duration) chan interface{} {
		channel := make(chan interface{})

		go func() {
			time.Sleep(duration)
			close(channel)
		}()

		return channel
	}

	start := time.Now()
	<-convertChannels(
		signal(2*time.Hour),
		signal(5*time.Minute),
		signal(3*time.Second),
		signal(1*time.Hour),
		signal(1*time.Minute),
		signal(10*time.Second),
	)
	fmt.Println(time.Since(start))
}
