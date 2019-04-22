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
	setupRenderLoop()

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

func setupRenderLoop() {
	messages <- "WASM::setupRenderLoop"
	var renderJSCallback js.Func
	renderJSCallback = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		window.Call("requestAnimationFrame", renderJSCallback)
		messages <- "WASM::requestAnimationFrame"
		drawFrame()
		return nil
	})
	window.Call("requestAnimationFrame", renderJSCallback)
}

func resetWindowSize() {
	windowSize.width = window.Get("innerWidth").Float()
	windowSize.height = window.Get("innerHeight").Float()
	canvas.Set("width", windowSize.width)
	canvas.Set("height", windowSize.height)
	messages <- fmt.Sprintf("WASM::resetWindowSize (%f x %f)", windowSize.width, windowSize.height)
}

func drawFrame() {
	clearCanvas()
	strokeStyle("white")
	fillStyle("white")
	lineWidth(0.5)

	drawStrokeRect(5, 5, windowSize.width-10, windowSize.height-10)
}

func clearCanvas() {
	context.Call("clearRect", 0, 0, windowSize.width, windowSize.height)
}

func strokeStyle(style string) {
	context.Set("strokeStyle", style)
}

func fillStyle(style string) {
	context.Set("fillStyle", style)
}

func lineWidth(width float64) {
	context.Set("lineWidth", width)
}

func drawStrokeRect(x, y, width, height float64) {
	context.Call("strokeRect", x, y, width, height)
}

func drawFillRect(x, y, width, height float64) {
	context.Call("fillRect", x, y, width, height)
}
