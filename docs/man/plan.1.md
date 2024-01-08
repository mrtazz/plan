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

`plan daily-prep [<options>]`


# EXAMPLES

`plan daily-prep`
: Create a daily note based on the configure template.


# DESCRIPTION

`plan` is a command line utility to help with personal planning and note
taking. The premise is that notes are plain text files (assumed to be in
markdown format) in a folder. The tooling can then generate a daily notes
template based on

# Configuration File

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

github:
  task_query: "assignee:mrtazz org:github state:open"
```

## Recurring tasks

`plan` supports simple weekly recurring tasks which can be added to the
configuration file. This is an array of strings that are available in the
template via the `{{ .RecurringTasks }}` templating macro.

`plan` also provides assigned tasks via the `{{ .AssignedTasks }}` macro in
the same way. Tasks here are taken from GitHub and have a `{{  .Name }}` and
`{{ .URL }}` attribute available for template rendering. In order for this to
work, a task query in GitHub search format needs to be configured.


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
