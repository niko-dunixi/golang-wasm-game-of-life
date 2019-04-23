//+build js,wasm

package main

import (
	"fmt"
	"math"
	"math/rand"
	"net/url"
	"strconv"
	"syscall/js"
	"time"
)

var (
	messages      chan string
	window        = js.Global()
	canvas        js.Value
	context       js.Value
	windowSize    struct{ width, height float64 }
	rows, columns int
	random        *rand.Rand
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

	pageUrl := document.Get("location").Get("href").String()
	params := parseUrlQueryParams(pageUrl)
	messages <- fmt.Sprintf("WASM::setupCanvas Params: %+v", params)
	random = initializeRandom(params["seed"])

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

func parseUrlQueryParams(pageUrl string) (params map[string]int64) {
	currentTimeAsSeed := time.Now().UnixNano()
	params = map[string]int64{
		"rows":    10,
		"columns": 10,
		"seed":    currentTimeAsSeed,
	}
	parse, e := url.Parse(pageUrl)
	if e != nil {
		return
	}
	for paramKey, paramValues := range parse.Query() {
		if len(paramValues) > 0 {
			if value, e := strconv.ParseInt(paramValues[0], 10, 64); e == nil {
				params[paramKey] = value
			} else {
				params[paramKey] = -1
			}
		}
	}
	return
}

func initializeRandom(seed int64) *rand.Rand {
	messages <- fmt.Sprintf("WASM::initializeRandom using seed: %d", seed)
	source := rand.NewSource(seed)
	return rand.New(source)
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
	// https://stackoverflow.com/a/8486324/1478636
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
	padding := float64(4)

	squareSize := math.Min(windowSize.width/10, windowSize.height/10)
	for row := 0; row < 10; row ++ {
		for column := 0; column < 10; column++ {
			x := float64(column)*squareSize + padding
			y := float64(row)*squareSize + padding
			side := squareSize - padding*2
			drawStrokeRect(x, y, side, side)
		}
	}
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
