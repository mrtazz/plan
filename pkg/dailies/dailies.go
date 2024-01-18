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
	defaultNotesDir          = "./dailies"
	defaultNotesSuffix       = ".md"
	defaultNotesFormat       = "20060102"
	defaultContentDateFormat = "2006-01-02"
	defaultNotesTemplate     = `## Overview

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
	date           time.Time
	dateFormat     string
	RecurringTasks []task.Task
	AssignedTasks  []task.Task
}

func (c Content) Weekday() string {
	return c.date.Weekday().String()
}

func (c Content) FormattedDate() string {
	return c.date.Format(c.dateFormat)
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
			date:           clk.Now(),
			dateFormat:     defaultContentDateFormat,
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

func (n *Note) WithDateFormat(format string) *Note {
	newNote := n
	newNote.content.dateFormat = format
	return newNote
}

func (n *Note) WithDate(t time.Time) *Note {
	newNote := n
	clk := clock.NewMock()
	clk.Set(t)
	newNote.clock = clk
	newNote.content.date = clk.Now()
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
