package github

import (
	"context"
	"fmt"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

type Task struct {
	title string
	url   string
}

// task.Task interface implementation
func (t *Task) Name() string { return t.title }
func (t *Task) URL() string  { return t.url }

var (
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
)

func GetAssignedTasks(token, searchQuery string) ([]Task, error) {

	tasks := make([]Task, 0, 10)

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
		if edge.Node.Issue.Title != "" &&
			edge.Node.Issue.Url != "" {
			tasks = append(tasks, Task{
				title: fmt.Sprintf("%s", edge.Node.Issue.Title),
				url:   fmt.Sprintf("%s", edge.Node.Issue.Url),
			})
		}
		if edge.Node.PullRequest.Title != "" &&
			edge.Node.PullRequest.Url != "" {
			tasks = append(tasks, Task{
				title: fmt.Sprintf("%s", edge.Node.PullRequest.Title),
				url:   fmt.Sprintf("%s", edge.Node.PullRequest.Url),
			})
		}
	}

	return tasks, nil

}
