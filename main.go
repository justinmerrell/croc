package main

import (
	"fmt"
	"log"
	"os"
	"syscall/js"

	"github.com/schollz/croc/v9/src/cli"
)

func main() {

	web_socket := js.Global().Get("WebSocket").New("ws://localhost:8080")
	web_socket.Call("addEventListener", "open", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		fmt.Println("open")

		web_socket.Call("send", "Hello from Go WASM")
		return nil
	}))

	wait_channel := make(chan struct{}, 0)

	println("Go croc WebAssembly Initialized")

	js.Global().Set("crocJS", js.FuncOf(crocJS))

	println("Go croc WebAssembly Ready")

	<-wait_channel
}

func crocJS(this js.Value, args []js.Value) interface{} {
	println("Go croc WebAssembly Called")
	println("Go croc WebAssembly Args: ", args[0].String())

	os.Setenv("CROC_PASSPHRASE", args[0].String())
	os.Setenv("CROC_CONFIG_DIR", "/") // TEST: Manually sets config dir for WASM

	// if len(args) < 1 {
	// 	return nil
	// }

	// Run the CLI
	if err := cli.Run(); err != nil {
		log.Fatalln(err)
	}

	return nil
}
