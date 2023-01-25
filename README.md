# plan

cli tool for some planning tasks

This is mostly support tooling for https://github.com/mrtazz/vim-plan


## Usage

```
% plan --help
Usage: plan <command>

Flags:
  -h, --help                        Show context-sensitive help.
      --debug                       Enable debug mode.
      --config-file=".plan.yaml"    path to the config file.

Commands:
  daily-prep
    create the daily note file

  version
    print version and exit.

Run "plan <command> --help" for more information on a command.
```

## Configuration

`plan` looks for a `.plan.yaml` config file in the current directory or can be
passed a config file to use as a flag. There are a couple of things that can
be configured via the config file, like recurring tasks, the daily template to
use, and the query to use to get GitHub issues.


### Example config file


```yaml
recurring_tasks:
  Monday:
    - plan the week

  Friday:
    - update achievements doc

daily_template: |-
  ## Misc

  It's {{ .Weekday }} today.
  You have {{ len .RecurringTasks }} recurring tasks today:
  {{ range .RecurringTasks }}
  {{- .Name }}
  {{- end }}

  You have {{ len .AssignedTasks }} issues assigned:
  {{- range .AssignedTasks }}
  - [ ] {{ .Name }} [link]({{ .URL }})
  {{- end }}

  ## Log

github:
  task_query: "assignee:mrtazz org:github state:open"
```
