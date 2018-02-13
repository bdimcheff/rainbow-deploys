package main

import (
	"fmt"
	"net/http"
	"os"
)

func getColor() string {
	return os.Getenv("COLOR")
}

func colorHandler(w http.ResponseWriter, r *http.Request) {
	color := getColor()

	fmt.Printf("Serving color: #%s", color)
	fmt.Println()

	fmt.Fprintf(w, "<body bgcolor=\"#%s\"><h1>#%s</h1></body>", color, color)
}

func main() {
	color := getColor()

	fmt.Printf("Booted with color: #%s", color)
	fmt.Println()

	http.HandleFunc("/", colorHandler)

	fmt.Println("listening on :8080")
	http.ListenAndServe(":8080", nil)
}
