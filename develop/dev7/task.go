package main

import (
	"context"
	"fmt"
	"os"
	"time"
)

func convertChannels(channels ...chan interface{}) <-chan interface{} {
	outputChan := make(chan interface{})

	ctx, cancel := context.WithCancel(context.Background())
	for _, channel := range channels {
		go func(listen, done chan interface{}, cancelFunc context.CancelFunc) {
			select {
			case <-ctx.Done():
				fmt.Fprintf(os.Stdout, "closing\n")
				return
			case s := <-listen:
				done <- s
				cancelFunc()
				close(done)
			}
		}(channel, outputChan, cancel)
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
