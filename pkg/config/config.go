package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v3"

	"github.com/mrtazz/plan/pkg/task"
)

type Config struct {
	RecurringTasks map[string][]string `yaml:"recurring_tasks"`
	DailyTemplate  string              `yaml:"daily_template"`
	GitHub         struct {
		Token     string `yaml:"token,omitempty"`
		TaskQuery string `yaml:"task_query"`
	} `yaml:"github,omitempty"`
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
