package config

import (
	"io/ioutil"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/mrtazz/plan/pkg/task"
)

type Config struct {
	RecurringTasks map[string][]string `yaml:"recurring_tasks"`
	DailyTemplate  string              `yaml:"daily_template"`
	DateFormat     string              `yaml:"date_format"`
	GitHub         struct {
		Token     string `yaml:"token,omitempty"`
		TaskQuery string `yaml:"task_query"`
	} `yaml:"github,omitempty"`
	ScreenshotImport struct {
		Source      string `yaml:"source"`
		Destination string `yaml:"destination"`
		FileFormat  string `yaml:"file_format"`
	} `yaml:"screenshot_import,omitempty"`
}

type ValidationError struct {
	message string
}

func (v ValidationError) Error() string {
	return v.message
}

func (c Config) GetRecurringTasks(weekday string) []task.Task {
	ret := make([]task.Task, 0, len(c.RecurringTasks))
	for _, t := range c.RecurringTasks[weekday] {
		ret = append(ret, task.New(t, ""))
	}
	return ret
}

func LoadConfigFromFile(filename string) (Config, error) {
	var c Config
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return c, err
	}
	err = yaml.Unmarshal(yamlFile, &c)
	return c, err
}

func ValidateConfig(filename string) error {
	c, err := LoadConfigFromFile(filename)
	if err != nil {
		return ValidationError{message: err.Error()}
	}

	if _, err = time.Parse(c.DateFormat, c.DateFormat); err != nil {
		return ValidationError{message: err.Error()}
	}

	return nil
}
