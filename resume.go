package main

import (
	"encoding/json"
	"strings"

	"github.com/volatiletech/null/v8"
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
	Title    Text       `json:"title"`
	Subtitle Text       `json:"subtitle"`
	Labels   []Linkable `json:"labels"`
	Links    []Linkable `json:"links"`
}

func (h *Header) Fill(variables map[string]string) {
	h.Title.Fill(variables)
	h.Subtitle.Fill(variables)
	for i, l := range h.Labels {
		l.Fill(variables)
		h.Labels[i] = l
	}
	for i, l := range h.Links {
		l.Fill(variables)
		h.Links[i] = l
	}
}

type Sidebar struct {
	Skills   [][]Text  `json:"skills"`
	Projects []Project `json:"projects"`
}

func (s *Sidebar) Fill(variables map[string]string) {
	for i, skills := range s.Skills {
		for j, skill := range skills {
			skill.Fill(variables)
			skills[j] = skill
		}
		s.Skills[i] = skills
	}
	for i, p := range s.Projects {
		p.Fill(variables)
		s.Projects[i] = p
	}
}

type Project struct {
	Entry
	Description Text `json:"description"`
}

func (p *Project) Fill(variables map[string]string) {
	p.Title.Fill(variables)
	for i, l := range p.Labels {
		l.Fill(variables)
		p.Labels[i] = l
	}
	for i, t := range p.Technologies {
		t.Fill(variables)
		p.Technologies[i] = t
	}
	p.Description.Fill(variables)
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
	Bullets    []Text `json:"bullets"`
	StartedOn  Text   `json:"startedOn"`
	FinishedOn Text   `json:"finishedOn"`
}

func (p *Position) Fill(variables map[string]string) {
	p.Title.Fill(variables)
	for i, l := range p.Labels {
		l.Fill(variables)
		p.Labels[i] = l
	}
	for i, t := range p.Technologies {
		t.Fill(variables)
		p.Technologies[i] = t
	}
	for i, b := range p.Bullets {
		b.Fill(variables)
		p.Bullets[i] = b
	}
	p.StartedOn.Fill(variables)
	p.FinishedOn.Fill(variables)
}

type Linkable struct {
	URL   null.String `json:"url"`
	Label Text        `json:"label"`
}

func (l *Linkable) Fill(variables map[string]string) {
	l.Label.Fill(variables)
}

type Text string

func (t Text) String() string {
	return string(t)
}

func (t *Text) Fill(variables map[string]string) {
	if t == nil || string(*t) == "" || string(*t)[0] != '$' {
		return
	}
	parts := strings.SplitN(string(*t), "/", 2)
	if value, ok := variables[parts[0]]; ok {
		*t = Text(value)
		return
	}
	if len(parts) == 2 {
		*t = Text(parts[1])
		return
	}
}

type Entry struct {
	Title        Text       `json:"title"`
	Labels       []Linkable `json:"labels,omitempty"`
	Location     Text       `json:"location,omitempty"`
	Technologies []Text     `json:"technologies,omitempty"`
}
