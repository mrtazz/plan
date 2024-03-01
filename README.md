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
  add-todo <text> ...
    Add a todo to the current daily note.

  add-note <text> ...
    Add a note to the current daily note.

  daily-prep
    create the daily note file

  get-assigned-issues
    Retrieve assigned issues

  import-screenshots
    Import screenshots to the plan folder

  validate-config
    Validate the passed config and return.

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

screenshot_import:
  source: "~/Desktop"
  destination: "dailies/20060102_attachments"
  file_format: "Screenshot 2006-01-02 at 15.04.05.png"
```

## Installation

There are pre-built binaries [on the releases page](https://github.com/mrtazz/plan/releases/)
or the go standard way should also work, e.g.:

```
go install github.com/mrtazz/plan@latest
```
