package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/a-h/templ"
)

func main() {
	data, err := os.ReadFile("resume.json")
	if err != nil {
		panic(err)
	}

	resume, err := LoadResumeJSON(data)
	if err != nil {
		panic(err)
	}
	resume.Fill(nil)

	http.Handle("/", templ.Handler(Index(resume)))

	fmt.Println("Listening on :8081")
	panic(http.ListenAndServe(":8081", nil))

}
