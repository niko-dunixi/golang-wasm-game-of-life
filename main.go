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
	messages   chan string
	window     = js.Global()
	canvas     js.Value
	context    js.Value
	windowSize struct{ width, height float64 }
	random     *rand.Rand
	universe   Universe
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
	universe = NewRandomBufferedUniverse(int(params["rows"]), int(params["columns"]), random)

	canvas = document.Call("getElementById", "canvas")
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
		"rows":    50,
		"columns": 50,
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
	return NewRand(seed)
}

func setupRenderLoop() {
	messages <- "WASM::setupRenderLoop"
	var renderJSCallback js.Func
	renderJSCallback = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		messages <- "WASM::requestAnimationFrame"
		if universe.Generation() < 4 {
			messages <- fmt.Sprintf("%s", universe)
		}
		draw()
		update()
		window.Call("setTimeout", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			window.Call("requestAnimationFrame", renderJSCallback)
			return nil
		}), 75)
		return nil
	})
	window.Call("requestAnimationFrame", renderJSCallback)
}

func resetWindowSize() {
	// https://stackoverflow.com/a/8486324/1478636
	squareWindowSide := math.Min(window.Get("innerWidth").Float(), window.Get("innerHeight").Float())
	windowSize.width = squareWindowSide
	windowSize.height = squareWindowSide
	canvas.Set("width", windowSize.width)
	canvas.Set("height", windowSize.height)
	messages <- fmt.Sprintf("WASM::resetWindowSize (%f x %f)", windowSize.width, windowSize.height)
}

func draw() {
	clearCanvas()
	color := "#ffffff"
	strokeStyle(color)
	fillStyle(color)
	lineWidth(0.75)
	padding := float64(4)
	innerPadding := 0.4 * padding

	squareSize := math.Min(windowSize.width/float64(universe.ColumnCount()), windowSize.height/float64(universe.RowCount()))
	side := squareSize - padding*2

	for row := 0; row < universe.RowCount(); row++ {
		for column := 0; column < universe.ColumnCount(); column++ {
			x := float64(column)*squareSize + padding
			y := float64(row)*squareSize + padding
			drawStrokeRect(x, y, side, side)
			if universe.IsAlive(row, column) {
				drawFillRect(x+innerPadding, y+innerPadding, side-(2*innerPadding), side-(2*innerPadding))
			}
		}
	}
}

func update() {
	universe.Iterate()
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
