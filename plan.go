package main

import (
	"fmt"
	"os"
	"time"

	"github.com/alecthomas/kong"
	"github.com/mrtazz/plan/pkg/github"
	log "github.com/sirupsen/logrus"

	"github.com/mrtazz/plan/pkg/config"
	"github.com/mrtazz/plan/pkg/dailies"
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

	version = "0.1.0"
)

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

	cfg, err := config.LoadConfigFromFile(flags.ConfigFile)
	if err != nil {
		log.WithFields(log.Fields{
			"error":       err.Error(),
			"config-file": flags.ConfigFile,
		}).Error("failed to parse config")
		os.Exit(1)
	}

	token := os.Getenv("ISSUES_TOKEN_GITHUB")

	assignedIssues, err := github.GetAssignedTasks(token, cfg.GitHub.TaskQuery)

	if err != nil {
		fmt.Println("error getting assigned issues")
		fmt.Println(err)
	}

	todayPlan := dailies.NewNote(assignedIssues,
		cfg.GetRecurringTasks(time.Now().Weekday().String()))

	if cfg.DailyTemplate != "" {
		todayPlan = todayPlan.WithTemplate(cfg.DailyTemplate)
	}

	dailyNoteString, err := todayPlan.Render()
	if err != nil {
		fmt.Println("error rendering daily note")
		fmt.Println(err)
		os.Exit(1)
	}

	if flags.DailyPrep.NoDryRun {
		printAndAppendToFile(todayPlan.Filepath(), dailyNoteString)
	} else {
		fmt.Println(dailyNoteString)
	}

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
