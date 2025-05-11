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
	Header         Header     `json:"header"`
	Summary        Text       `json:"summary"`
	Skills         []Text     `json:"skills"`
	Positions      []Position `json:"positions"`
	PositionsCount int        `json:"positionMaxCount"`
	Projects       []Project  `json:"projects"`
}

func (r *Resume) Fill(variables map[string]string) {
	r.Header.Fill(variables)
	r.Summary.Fill(variables)
	for i, skill := range r.Skills {
		skill.Fill(variables)
		r.Skills[i] = skill
	}
	for i, pos := range r.Positions {
		pos.Fill(variables)
		r.Positions[i] = pos
	}
	for i, proj := range r.Projects {
		proj.Fill(variables)
		r.Projects[i] = proj
	}
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

func (t Text) html() string {
	// Define a regex to find text between ** **
	re := regexp.MustCompile(`\*\*(.*?)\*\*`)

	// Replace all matches with the desired HTML span
	return re.ReplaceAllStringFunc(string(t), func(match string) string {
		// Extract the content without the **
		content := strings.TrimPrefix(strings.TrimSuffix(match, "**"), "**")
		return `<span class="font-bold">` + content + `</span>`
	})
}

func (t Text) Render(ctx context.Context, w io.Writer) error {
	_, err := io.WriteString(w, t.html())
	return err
}

func (t Text) String() string {
	return t.html()
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
