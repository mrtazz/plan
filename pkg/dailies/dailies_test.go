package dailies

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/mrtazz/plan/pkg/task"
)

func TestDailyNoteRender(t *testing.T) {
	assert := assert.New(t)
	tests := map[string]struct {
		assignedTasks   []task.Task
		recurringTasks  []task.Task
		expectedContent string
		day             string
	}{
		"default": {
			day: "2023-01-20",
			assignedTasks: []task.Task{
				task.New("assigned foo", "http://example.com/tasks/1"),
			},
			recurringTasks: []task.Task{
				task.New("recurring foo", "http://example.com/tasks/2"),
			},
			expectedContent: `## Overview

It's Friday today.

## Tasks
You have 1 recurring tasks today:
recurring foo

You have 1 assigned tasks:
- [ ] assigned foo [link](http://example.com/tasks/1)

## Log

`,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			day, err := time.Parse("2006-01-02", tc.day)
			n := NewNote(tc.assignedTasks, tc.recurringTasks).WithDate(day)
			content, err := n.Render()
			assert.Equal(nil, err)
			assert.Equal(tc.expectedContent, content)
		})
	}
}
