---
title: plan
section: 1
footer: plan VERSION_PLACEHOLDER
header: plan User's Manual
author: Daniel Schauenberg <d@unwiredcouch.com>
date: DATE_PLACEHOLDER
---

<!-- This is the plan(1) man page, written in Markdown. -->
<!-- To generate the roff version, run `make man` -->

# NAME

plan — a command line tool for personal planning


# SYNOPSIS

`plan <subcommand> [<options>]`

# SUBCOMMANDS

`plan daily-prep`
: Create a daily note based on the configure template.

`plan validate-config --config_file=<plan.yaml>`
: Validate a provided config file.


# DESCRIPTION

`plan` is a command line utility to help with personal planning and note
taking. The premise is that notes are plain text files (assumed to be in
markdown format) in a folder. The tooling can then generate a daily note based
on a configured (or default) template. The main interaction for notes is
intended to be your favourite editor while `plan` provides tooling around it.

# CONFIGURATION FILE

`plan` reads a `.plan.yaml` configuration file in the current directory. That
way the config file can be colocated in the notes directory and different
settings can be used for different notes directories.

An example configuration file in yaml format can look like this

```
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
  - [ ] {{- .Name }}
  {{- end }}

  You have {{ len .AssignedTasks }} issues assigned:
  {{- range .AssignedTasks }}
  - [ ] {{ .Name }} [link]({{ .URL }})
  {{- end }}

  ## Log

date_format: "2006-01-02"

github:
  task_query: "assignee:mrtazz org:github state:open"
```

Validity of a configuration file can be checked with the `validate-config`
subcommand.

## RECURRING TASKS

`plan` supports simple weekly recurring tasks which can be added to the
configuration file. This is an array of strings that are available in the
template via the `{{ .RecurringTasks }}` templating macro.

`plan` also provides assigned tasks via the `{{ .AssignedTasks }}` macro in
the same way. Tasks here are taken from GitHub and have a `{{  .Name }}` and
`{{ .URL }}` attribute available for template rendering. In order for this to
work, a task query in GitHub search format needs to be configured.

## NOTE TEMPLATE RENDERING

The daily note is rendered with a small set of context variables that can be
used in templating. The templating engine is Go's `text/template` and should
be expected to behave in the same way:

- `{{ .Weekday }}`: Prints the current day of the week
- `{{ .FormattedDate }}`: Prints the current date formatted based on the
  `date_format` config key (default: "2006-01-02")
- `{[ .RecurringTasks }}`: An array of recurring tasks for the day as set in the
  configuration file. Can be iterated over and accessed with `{{ .Name }}`
  within
- `{{ .AssignedTasks }}`: An array of assigned issues from GitHub, based on
  the configured `github.task_query`. Can be iterated over and accessed with
  `{{ .Name }}` and `{{ .URL }}`


# META OPTIONS AND COMMANDS

`--help`
: Show list of command-line options.

`version`
: Show version of plan.



# AUTHOR

plan is maintained by mrtazz.

**Source code:** `https://github.com/mrtazz/plan`

# REPORTING BUGS

- https://github.com/mrtazz/plan/issues

# SEE ALSO

- https://github.com/mrtazz/vim-plan
