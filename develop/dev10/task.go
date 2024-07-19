package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/websocket"
)

var dialer = websocket.Dialer{}

func main() {
	timeout := flag.Duration("timeout", 10*time.Second, "set timeout")
	flag.Parse()

	urlString, err := getUrl(flag.Args())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	dialer.HandshakeTimeout = *timeout

	connection, _, err := dialer.Dial(urlString, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error while connecting to websocket: %v\n", err)
		os.Exit(1)
	}
	defer connection.Close()

	fmt.Fprintln(os.Stdout, "connection established")

	connClosed := make(chan struct{})
	go func() {
		defer close(connClosed)
		for {
			_, message, err := connection.ReadMessage()
			if err != nil {
				fmt.Fprintf(os.Stderr, "error while reading message: %v\n", err)
				return
			}

			fmt.Fprintln(os.Stdout, string(message))
		}
	}()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	reader := bufio.NewReader(os.Stdin)

	for {
		select {
		case <-connClosed:
			fmt.Fprintln(os.Stderr, "conn is closed")
			return
		case <-interrupt:
			fmt.Fprintln(os.Stderr, "syscall interrupt")

			if err := connection.WriteMessage(websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")); err != nil {
				fmt.Fprintf(os.Stderr, "Write error: %v", err)
			}
			select {
			case <-connClosed:
			case <-time.After(2 * time.Second):
			}
		default:
			text, err := reader.ReadBytes('\n')
			if err != nil {
				if err == io.EOF {
					fmt.Fprintln(os.Stdout, "EOF ending...")

					if err := connection.WriteMessage(websocket.CloseMessage,
						websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")); err != nil {
						fmt.Fprintf(os.Stderr, "Write error: %v", err)
					}
					select {
					case <-connClosed:
					case <-time.After(2 * time.Second):
					}
					return
				}
				fmt.Fprintf(os.Stderr, "error while scanning text: %v", err)
				return
			}

			if err := connection.WriteMessage(websocket.TextMessage, text); err != nil {
				fmt.Fprintf(os.Stderr, "error while writing message: %v", err)
				return
			}
		}
	}

}

func getUrl(arguments []string) (string, error) {
	if len(arguments) != 2 {
		return "", fmt.Errorf("need HOST + PATH")
	}

	url := url.URL{
		Scheme: "wss",
		Host:   arguments[0],
		Path:   arguments[1],
	}

	fmt.Fprintf(os.Stdout, "connecting: %s\n", url.String())

	return url.String(), nil
}
