package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	fmt.Println(os.Getwd())
	http.HandleFunc("/", IndexHtml)
	http.HandleFunc("/wasm_exec.js", WasmExecJs)
	http.HandleFunc("/client.wasm", ClientWasm)

	srv := &http.Server{
		Addr: ":4141",
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalln(err)
		}
	}()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Running forever until exit requested")
	<-stop
	log.Println("Exit requested. Now halting.")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	if err := srv.Shutdown(ctx); err != nil && err != http.ErrServerClosed {
		log.Fatalln(err)
	} else {
		log.Println("Clean exit")
	}
}

func IndexHtml(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	bytes, _ := ioutil.ReadFile("./bin/index.html")
	_, _ = w.Write(bytes)
}

func WasmExecJs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/javascript")
	bytes, _ := ioutil.ReadFile("./bin/wasm_exec.js")
	_, _ = w.Write(bytes)
}

func ClientWasm(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/wasm")
	bytes, _ := ioutil.ReadFile("./bin/client.wasm")
	_, _ = w.Write(bytes)
}
