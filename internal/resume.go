package internal

import (
	"context"
	"encoding/json"
	"io"
	"regexp"
	"strings"
)

func LoadResumeJSON(data []byte) (Resume, error) {
	var resume Resume
	err := json.Unmarshal(data, &resume)
	if err != nil {
		return Resume{}, err
	}
	return resume, nil
}

type Resume struct {
	Header  Header  `json:"header"`
	Content Content `json:"content"`
	Sidebar Sidebar `json:"sidebar"`
}

func (r *Resume) Fill(variables map[string]string) {
	r.Header.Fill(variables)
	r.Content.Fill(variables)
	r.Sidebar.Fill(variables)
}

type Header struct {
	Title    Text `json:"title"`
	Subtitle Text `json:"subtitle"`
	Email    Text `json:"email"`
	GitHub   Text `json:"github"`
	LinkedIn Text `json:"linkedIn"`
}

func (h *Header) Fill(variables map[string]string) {
	h.Title.Fill(variables)
	h.Subtitle.Fill(variables)
}

type Sidebar struct {
	Skills   [][]string `json:"skills"`
	Projects []Project  `json:"projects"`
}

func (s *Sidebar) Fill(variables map[string]string) {
	for i, p := range s.Projects {
		p.Fill(variables)
		s.Projects[i] = p
	}
}

type Project struct {
	Entry
	Link    string `json:"link,omitempty"`
	Bullets []Text `json:"bullets"`
}

func (p *Project) Fill(variables map[string]string) {
	p.Title.Fill(variables)
	for i, b := range p.Bullets {
		b.Fill(variables)
		p.Bullets[i] = b
	}
}

type Content struct {
	Summary   Text       `json:"summary"`
	Positions []Position `json:"positions"`
	ExpandURL Text       `json:"expandUrl"`
}

func (c *Content) Fill(variables map[string]string) {
	c.Summary.Fill(variables)
	for i, p := range c.Positions {
		p.Fill(variables)
		c.Positions[i] = p
	}
	c.ExpandURL.Fill(variables)
}

type Position struct {
	Entry
	Company    string `json:"company,omitempty"`
	Bullets    []Text `json:"bullets"`
	StartedOn  Text   `json:"startedOn"`
	FinishedOn Text   `json:"finishedOn"`
}

func (p *Position) Fill(variables map[string]string) {
	p.Title.Fill(variables)
	for i, b := range p.Bullets {
		b.Fill(variables)
		p.Bullets[i] = b
	}
	p.StartedOn.Fill(variables)
	p.FinishedOn.Fill(variables)
}

type Linkable struct {
	URL   Text `json:"url"`
	Label Text `json:"label"`
}

func (l *Linkable) Fill(variables map[string]string) {
	l.Label.Fill(variables)
	l.URL.Fill(variables)
}

type Text string

func (t Text) Render(ctx context.Context, w io.Writer) error {
	// Define a regex to find text between ** **
	re := regexp.MustCompile(`\*\*(.*?)\*\*`)

	// Replace all matches with the desired HTML span
	result := re.ReplaceAllStringFunc(string(t), func(match string) string {
		// Extract the content without the **
		content := strings.TrimPrefix(strings.TrimSuffix(match, "**"), "**")
		return `<span class="font-bold">` + content + `</span>`
	})

	// Write the result to the writer
	_, err := io.WriteString(w, result)
	return err
}

func (t Text) String() string {
	return string(t)
}

func (t *Text) Fill(variables map[string]string) {
	if t == nil || string(*t) == "" || string(*t)[0] != '$' {
		return
	}
	parts := strings.SplitN(string(*t), "/", 2)
	if value, ok := variables[parts[0][1:]]; ok {
		*t = Text(value)
		return
	}
	if len(parts) == 2 {
		*t = Text(parts[1])
		return
	}
}

type Entry struct {
	Title        Text     `json:"title"`
	Location     string   `json:"location,omitempty"`
	Technologies []string `json:"technologies,omitempty"`
}
