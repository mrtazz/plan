package main

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"text/template"
	"time"

	"github.com/alecthomas/kong"
	"github.com/shurcooL/githubv4"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"gopkg.in/yaml.v3"
)

var (
	flags struct {
		Debug      bool   `help:"Enable debug mode."`
		ConfigFile string `default:".plan.yaml" help:"path to the config file."`
		DailyPrep  struct {
			NoDryRun bool `help:"whether to actually write the daily note"`
		} `cmd:"" help:"create the daily note file"`
		Version struct {
		} `cmd:"" help:"print version and exit."`
	}

	query struct {
		Search struct {
			IssueCount githubv4.Int
			Edges      []struct {
				Node struct {
					Issue struct {
						Title string
						Url   string
					} `graphql:"... on Issue"`
					PullRequest struct {
						Title string
						Url   string
					} `graphql:"... on PullRequest"`
				}
			}
		} `graphql:"search(first: 100, type: ISSUE, query:$searchQuery)"`
	}

	version = "0.1.0"
)

type config struct {
	RecurringTasks map[string][]string `yaml:"recurring_tasks"`
	DailyTemplate  string              `yaml:"daily_template"`
	GitHub         struct {
		Token     string `yaml:"token,omitempty"`
		TaskQuery string `yaml:"task_query"`
	} `yaml:"github,omitempty"`
}

func LoadConfigFromFile(filename string) (config, error) {
	var c config
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return c, err
	}
	err = yaml.Unmarshal(yamlFile, &c)
	return c, err
}

type githubTask struct {
	Title string
	URL   string
}

type content struct {
	Weekday        string
	RecurringTasks []string
	AssignedIssues []githubTask
}

func main() {

	ctx := kong.Parse(&flags)
	switch ctx.Command() {
	case "daily-prep":
		dailyPrep()
	case "version":
		fmt.Printf(version)
		return
	default:
		log.Fatal("Unknown command: " + ctx.Command())
	}
}

func dailyPrep() {

	cfg, err := LoadConfigFromFile(flags.ConfigFile)
	if err != nil {
		log.WithFields(log.Fields{
			"error":       err.Error(),
			"config-file": flags.ConfigFile,
		}).Error("failed to parse config")
		os.Exit(1)
	}

	now := time.Now()
	token := os.Getenv("ISSUES_TOKEN_GITHUB")

	weekday := now.Weekday().String()
	dailyNote := fmt.Sprintf("./dailies/%s.md", now.Format("20060102"))

	fmt.Printf("Today's journal file is %s.\n", dailyNote)
	todayTasks := cfg.RecurringTasks[weekday]

	c := content{
		Weekday: weekday,
	}

	for _, task := range todayTasks {
		markdownTask := fmt.Sprintf("- [ ] %s\n", task)
		c.RecurringTasks = append(c.RecurringTasks, markdownTask)
	}

	if c.AssignedIssues, err = GetAssignedTasks(token, cfg.GitHub.TaskQuery); err != nil {
		fmt.Println("error getting assigned issues")
		fmt.Println(err)
	}

	out, err := renderTemplate(cfg.DailyTemplate, c)
	if err != nil {
		fmt.Println("error rendering template")
		fmt.Println(err)
		os.Exit(1)
	}

	if flags.DailyPrep.NoDryRun {
		printAndAppendToFile(dailyNote, out)
	} else {
		fmt.Println(out)
	}

}

func renderTemplate(tpl string, c content) (string, error) {
	tmpl, err := template.New("test").Parse(tpl)
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

func printAndAppendToFile(filename, content string) error {
	fmt.Println(content)
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(content)
	return err
}

func GetAssignedTasks(token, searchQuery string) ([]githubTask, error) {

	tasks := make([]githubTask, 0, 10)

	variables := map[string]interface{}{
		"searchQuery": githubv4.String(searchQuery),
	}

	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	httpClient := oauth2.NewClient(context.Background(), src)

	client := githubv4.NewClient(httpClient)

	if err := client.Query(context.Background(), &query, variables); err != nil {
		return tasks, err
	}
	for _, edge := range query.Search.Edges {
		t := githubTask{}
		if edge.Node.Issue.Title != "" &&
			edge.Node.Issue.Url != "" {
			t.Title = fmt.Sprintf("%s", edge.Node.Issue.Title)
			t.URL = fmt.Sprintf("%s", edge.Node.Issue.Url)
		}
		if edge.Node.PullRequest.Title != "" &&
			edge.Node.PullRequest.Url != "" {
			t.Title = fmt.Sprintf("%s", edge.Node.PullRequest.Title)
			t.URL = fmt.Sprintf("%s", edge.Node.PullRequest.Url)
		}
		tasks = append(tasks, t)
	}

	return tasks, nil

}
