package dailies

import (
	"bytes"
	"fmt"
	"text/template"
	"time"

	"github.com/benbjohnson/clock"

	"github.com/mrtazz/plan/pkg/task"
)

const (
	defaultNotesDir      = "./dailies"
	defaultNotesSuffix   = ".md"
	defaultNotesFormat   = "20060102"
	defaultNotesTemplate = `## Overview

It's {{ .Weekday }} today.

## Tasks
You have {{ len .RecurringTasks }} recurring tasks today:
{{ range .RecurringTasks }}
{{- .Name }}
{{- end }}

You have {{ len .AssignedTasks }} assigned tasks:
{{- range .AssignedTasks }}
- [ ] {{ .Name }} [link]({{ .URL }})
{{- end }}

## Log

`
)

type Content struct {
	Weekday        string
	RecurringTasks []task.Task
	AssignedTasks  []task.Task
}

type Note struct {
	clock    clock.Clock
	template string
	content  Content
}

func NewNote(assignedTasks, recurringTasks []task.Task) *Note {
	clk := clock.New()
	n := &Note{
		clock:    clk,
		template: defaultNotesTemplate,
		content: Content{
			Weekday:        clk.Now().Weekday().String(),
			RecurringTasks: recurringTasks,
			AssignedTasks:  assignedTasks,
		},
	}
	return n
}

func (n *Note) WithTemplate(tpl string) *Note {
	newNote := n
	n.template = tpl
	return newNote
}

func (n *Note) WithDate(t time.Time) *Note {
	newNote := n
	clk := clock.NewMock()
	clk.Set(t)
	n.clock = clk
	n.content.Weekday = clk.Now().Weekday().String()
	return newNote
}

func (n *Note) Filepath() string {
	return fmt.Sprintf("%s/%s%s",
		defaultNotesDir,
		n.clock.Now().Format(defaultNotesFormat),
		defaultNotesSuffix)
}

func (n *Note) Render() (string, error) {
	tmpl, err := template.New("dailyNote").Parse(n.template)
	if err != nil {
		return "", fmt.Errorf("template error: %w", err)
	}
	var b bytes.Buffer
	err = tmpl.Execute(&b, n.content)
	if err != nil {
		return "", fmt.Errorf("template error: %w", err)
	}
	return b.String(), nil
}
