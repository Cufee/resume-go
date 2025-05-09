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
	icon, err := os.Open("static/favicon.svg")
	if err != nil {
		panic(err)
	}
	buildIcon, err := os.Create(filepath.Join("build", "favicon.svg"))
	if err != nil {
		panic(err)
	}
	_, err = io.Copy(buildIcon, icon)
	if err != nil {
		panic(err)
	}

	data, err := os.ReadFile("static/resume.json")
	if err != nil {
		panic(err)
	}

	resume, err := internal.LoadResumeJSON(data)
	if err != nil {
		panic(err)
	}
	if len(resume.Content.Positions) > 4 {
		resume.Content.Positions = resume.Content.Positions[:4]
	}

	browser := rod.New().MustConnect()
	var wg sync.WaitGroup

	// default variation
	err = newVariation(browser, "", resume, nil)
	if err != nil {
		panic(err)
	}

	// other variations
	func() {
		data, err := os.ReadFile("static/variations.json")
		if err != nil {
			return
		}
		var variations map[string]map[string]string
		err = json.Unmarshal(data, &variations)
		if err != nil {
			return
		}

		for name, vars := range variations {
			err = newVariation(browser, name, resume, vars)
			if err != nil {
				return
			}
		}
	}()

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

	path, err := filepath.Rel(filepath.Join("build", name), filepath.Join("build"))
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
		Scale:        p(0.63),
		MarginTop:    p(0.25),
		MarginLeft:   p(0.35),
		MarginRight:  p(0.35),
		MarginBottom: p(0.0),
		PaperHeight:  p(11.0),
		PaperWidth:   p(8.27),
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
