//go:build ignore

package main

import (
	"context"
	"net/http"
	"os"
	"path/filepath"

	"github.com/cufee/resume-go/internal"
)

//go:generate templ generate
//go:generate bunx tailwindcss -i ./tailwind.css -o ./build/main.css --minify

func main() {
	data, err := os.ReadFile("static/resume.json")
	if err != nil {
		panic(err)
	}

	resume, err := internal.LoadResumeJSON(data)
	if err != nil {
		panic(err)
	}

	resume.Fill(nil)

	err = os.MkdirAll("build", os.ModePerm)
	if err != nil {
		panic(err)
	}

	out, err := os.Create(filepath.Join("build", "index.html"))
	if err != nil {
		panic(err)
	}
	defer out.Close()

	err = internal.Index("/main.css", resume).Render(context.Background(), out)
	if err != nil {
		panic(err)
	}

	http.Handle("/", http.FileServer(http.Dir("build/")))
	println("Listening on :8081")
	panic(http.ListenAndServe(":8081", nil))
}
