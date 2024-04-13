package backend

import (
	"fmt"
	"log"
	"net/http"
)

func helloWorld(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(res, "Hello World")
}

func Run(addr string) {
	http.HandleFunc("/", helloWorld)
	fmt.Println("Server started and listening on the port", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
