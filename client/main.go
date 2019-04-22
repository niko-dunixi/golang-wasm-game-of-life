//+build js,wasm

package main

func main() {
	messages := make(chan string)
	go func() {
		for message := range messages {
			println(message)
		}
	}()
}
