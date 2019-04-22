//+build js,wasm

package main

func main() {
	messages := make(chan string)
	go func() {
		for message := range messages {
			println(message)
		}
	}()
	messages <- "WASM::main"

	

	messages <- "WASM::This will now run forever!"
	runForever := make(chan bool)
	<-runForever
}
