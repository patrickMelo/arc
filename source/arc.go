package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"arc/database"
	"arc/server"
	"arc/vm"
)

func printUsage() {
	println("Usage: arc [mode]")
	println("")
	println("Available modes:")
	println("- client: run in client mode and connect to server at <localhost:8080>")
	println("- server: run in server mode at <localhost:8080>")
	println("- standalone: run in standalone mode")
}

func runServer() {
	var db = database.Create()
	log.Print("ARC: database created.")

	var runtime = vm.CreateRuntime(vm.StandardLibrary, db)
	log.Print("ARC: runtime created with standard library.")

	var server = server.Create(":8080", runtime)
	log.Print("ARC: server created to run at :8080.")

	log.Print("ARC: running...")
	defer log.Print("ARC: done.")

	log.Fatal(server.Run())
}

func runClient(standalone bool) {
	var db *database.Database
	var runtime *vm.Runtime

	if standalone {
		db = database.Create()
		log.Print("ARC: database created.")

		runtime = vm.CreateRuntime(vm.StandardLibrary, db)
		log.Print("ARC: runtime created with standard library.")

		log.Print("ARC: running in standalone mode, type HELP for help and EXIT to exit.")
	} else {
		log.Print("ARC: running in client mode, type HELP for help and EXIT to exit.")
	}

	defer log.Print("ARC: done.")

	var commandLine = ""
	var commandLineScanner = bufio.NewScanner(os.Stdin)

	print("> ")

	for commandLineScanner.Scan() {
		commandLine = commandLineScanner.Text()

		if strings.ToUpper(commandLine) == "EXIT" {
			break
		}

		if strings.ToUpper(commandLine) == "HELP" {
			for index := range vm.StandardLibrary {
				println(vm.StandardLibrary[index].GetHelp())
			}
		} else {
			if standalone {
				if result := runtime.Execute(commandLine); result != nil {
					println(strings.Join(result, " "))
				}
			} else {
				if httpResponse, err := http.Get("http://localhost:8080/?cmd=" + url.QueryEscape(commandLine)); err == nil {
					var bodyBuffer = new(bytes.Buffer)

					if bodySize, bodyError := bodyBuffer.ReadFrom(httpResponse.Body); (bodyError == nil) && (bodySize > 0) {
						println(bodyBuffer.String())
					}

					httpResponse.Body.Close()
				} else {
					fmt.Printf("HTTP ERROR: %v.\n", err)
				}

			}
		}

		print("> ")
	}
}

func main() {
	if len(os.Args) != 2 {
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "client":
		runClient(false)

	case "server":
		runServer()

	case "standalone":
		runClient(true)

	default:
		log.Fatalf("Unknown mode: %s\n", os.Args[1])
		os.Exit(1)
	}

	return
}
