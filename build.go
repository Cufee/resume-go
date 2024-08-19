package main

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/cufee/resume-go/internal"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

//go:generate templ generate
//go:generate bunx tailwindcss -i ./tailwind.css -o ./build/main.css --minify

func main() {
	data, err := os.ReadFile("static/resume.json")
	if err != nil {
		panic(err)
	}
	variationsData, err := os.ReadFile("static/resume_variations.json")
	if err != nil {
		panic(err)
	}

	resume, err := internal.LoadResumeJSON(data)
	if err != nil {
		panic(err)
	}

	browser := rod.New().MustConnect()
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		println("generating base variation")
		err = newVariation(browser, "", resume, nil)
		if err != nil {
			panic(err)
		}
		println("done generating base variation")
	}()

	var variations map[string]map[string]string
	err = json.Unmarshal(variationsData, &variations)
	if err != nil {
		panic(err)
	}

	for v, data := range variations {
		wg.Add(1)
		go func() {
			defer wg.Done()
			println("generating variation", v)
			err := newVariation(browser, v, resume, data)
			if err != nil {
				panic(err)
			}
			println("done generating variation", v)
		}()
	}

	wg.Wait()
	browser.Close()
	println("Done generating static assets")
}

func newVariation(browser *rod.Browser, name string, resume internal.Resume, variables map[string]string) error {
	var copy internal.Resume
	d, _ := json.Marshal(resume)
	json.Unmarshal(d, &copy)
	copy.Fill(variables)

	err := os.MkdirAll(filepath.Join("build", name), os.ModePerm)
	if err != nil {
		return err
	}

	out, err := os.Create(filepath.Join("build", name, "index.html"))
	if err != nil {
		return err
	}
	defer out.Close()

	path, err := filepath.Rel(filepath.Join("build", name), filepath.Join("build", "main.css"))
	if err != nil {
		return err
	}

	err = internal.Index(path, copy).Render(context.Background(), out)
	if err != nil {
		return err
	}

	source, _ := filepath.Abs(filepath.Join("build", name, "index.html"))
	target, err := os.Create(filepath.Join("build", name, "resume.pdf"))
	if err != nil {
		return err
	}
	defer out.Close()

	savePageAsPDF(browser, source, target)
	return nil
}

func savePageAsPDF(browser *rod.Browser, source string, target io.Writer) error {
	page := browser.MustPage(filepath.Join("file://", source)).MustWaitStable()

	reader, err := page.PDF(&proto.PagePrintToPDF{
		Scale:        p(0.58),
		MarginTop:    p(0.0),
		MarginLeft:   p(0.0),
		MarginRight:  p(0.0),
		MarginBottom: p(0.0),
		PaperHeight:  p(11.69),
		PaperWidth:   p(8.5),
		PageRanges:   "1",
	})
	if err != nil {
		return err
	}

	_, err = io.Copy(target, reader)
	if err != nil {
		return err
	}
	return nil
}

func p[T any](v T) *T {
	return &v
}
