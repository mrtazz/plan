package dailies

import (
	"bytes"
	"fmt"
	"text/template"
	"time"

	"github.com/mrtazz/plan/pkg/task"
)

const (
	defaultNotesDir    = "./dailies"
	defaultNotesSuffix = ".md"
)

type Content struct {
	Weekday        string
	RecurringTasks []task.Task
	AssignedIssues []task.Task
}

type Note struct {
	Filepath string
	Content  string
}

func CreateDailyNote(tpl string, ghTasks, recurringTasks []task.Task) (Note, error) {

	now := time.Now()
	weekday := now.Weekday().String()
	c := Content{
		Weekday:        weekday,
		RecurringTasks: recurringTasks,
		AssignedIssues: ghTasks,
	}

	dailyNotePath := fmt.Sprintf("%s/%s%s",
		defaultNotesDir,
		now.Format("20060102"),
		defaultNotesSuffix)

	out, err := renderTemplate(tpl, c)

	return Note{
		Filepath: dailyNotePath,
		Content:  out,
	}, err
}

func renderTemplate(tpl string, c Content) (string, error) {
	tmpl, err := template.New("dailyNote").Parse(tpl)
	if err != nil {
		return "", err
	}
	var b bytes.Buffer
	err = tmpl.Execute(&b, c)
	if err != nil {
		return "", err
	}
	return b.String(), nil
}
