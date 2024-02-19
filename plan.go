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
	"github.com/mrtazz/plan/pkg/screenshots"
)

var (
	flags struct {
		Debug      bool   `help:"Enable debug mode."`
		ConfigFile string `default:".plan.yaml" help:"path to the config file."`
		DailyPrep  struct {
			Day      string `help:"day to write daily note for in '2022-12-31' format"`
			NoDryRun bool   `help:"whether to actually write the daily note"`
		} `cmd:"" help:"create the daily note file"`
		GetAssignedIssues struct {
			Format string `default:"markdown" help:"which format to use for outputting issues"`
		} `cmd:"" help:"Retrieve assigned issues"`
		ImportScreenshots struct {
			NoDryRun bool `help:"whether to actually create folders and move files"`
		} `cmd:"" help:"Import screenshots to the plan folder"`
		ValidateConfig struct {
		} `cmd:"" help:"Validate the passed config and return."`
		Version struct {
		} `cmd:"" help:"print version and exit."`
	}

	version   = "dev"
	goversion = "na"
)

const (
	dayFlagFormat = "2006-01-02"
)

func main() {

	ctx := kong.Parse(&flags, kong.UsageOnError())
	switch ctx.Command() {

	case "daily-prep":
		dailyPrep()

	case "get-assigned-issues":
		getAssignedIssues()

	case "validate-config":
		if err := config.ValidateConfig(flags.ConfigFile); err != nil {
			fmt.Printf("Failed to validate config '%s', got error: %s\n", flags.ConfigFile, err.Error())
			os.Exit(1)
		}
		fmt.Printf("Config file '%s' is valid.\n", flags.ConfigFile)
		return

	case "import-screenshots":
		cfg, err := config.LoadConfigFromFile(flags.ConfigFile)
		if err != nil {
			log.WithFields(log.Fields{
				"error":       err.Error(),
				"config-file": flags.ConfigFile,
			}).Error("failed to parse config")
			os.Exit(1)
		}
		i := screenshots.NewImporter(
			cfg.ScreenshotImport.Source,
			cfg.ScreenshotImport.Destination,
			cfg.ScreenshotImport.FileFormat,
			flags.ImportScreenshots.NoDryRun,
		)

		if err = i.ImportToPlanFolder(); err != nil {
			fmt.Printf("Failed to import screenshots: %s\n", err.Error())
			os.Exit(1)
		}

	case "version":
		fmt.Printf("plan: version %s %s", version, goversion)
		return

	default:
		ctx.FatalIfErrorf(fmt.Errorf("Unknown command: " + ctx.Command()))
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

	if cfg.DateFormat != "" {
		todayPlan = todayPlan.WithDateFormat(cfg.DateFormat)
	}

	if cfg.DailyTemplate != "" {
		todayPlan = todayPlan.WithTemplate(cfg.DailyTemplate)
	}

	if flags.DailyPrep.Day != "" {
		t, err := time.Parse(dayFlagFormat, flags.DailyPrep.Day)
		if err != nil {
			fmt.Println("error parsing provided date")
			fmt.Println(err)
			os.Exit(1)
		}
		todayPlan = todayPlan.WithDate(t)
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

func getAssignedIssues() {
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

	switch flags.GetAssignedIssues.Format {
	case "markdown":
		for _, task := range assignedIssues {
			fmt.Printf("- [ ] [%s](%s)\n", task.Name(), task.URL())

		}

	default:
		fmt.Printf("Unknown format: '%s'\n", flags.GetAssignedIssues.Format)
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
