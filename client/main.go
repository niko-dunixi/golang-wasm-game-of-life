//+build js,wasm

package main

import (
	"fmt"
	"syscall/js"
)

var (
	messages   chan string
	window     = js.Global()
	canvas     js.Value
	context    js.Value
	windowSize struct{ width, height float64 }
)

func main() {
	messages = make(chan string)
	go func() {
		for message := range messages {
			println(message)
		}
	}()
	messages <- "WASM::main"
	setupCanvas()

	messages <- "WASM::This will now run forever!"
	runForever := make(chan bool)
	<-runForever
}

func setupCanvas() {
	messages <- "WASM::setupCanvas"
	document := window.Get("document")
	canvas = document.Call("createElement", "canvas")

	body := document.Get("body")
	body.Call("appendChild", canvas)

	context = canvas.Call("getContext", "2d")

	updateWindowSizeJSCallback := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		resetWindowSize()
		return nil
	})
	window.Call("addEventListener", "resize", updateWindowSizeJSCallback)
	resetWindowSize()
}

func resetWindowSize() {
	windowSize.width = window.Get("innerWidth").Float()
	windowSize.height = window.Get("innerHeight").Float()
	canvas.Set("width", windowSize.width)
	canvas.Set("height", windowSize.height)
	messages <- fmt.Sprintf("WASM::resetWindowSize (%f x %f)", windowSize.width, windowSize.height)
}
