package jobs

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/google/go-github/v40/github"
	"github.com/urfave/cli/v2"
)

func List(c *cli.Context) error {
	path := fmt.Sprintf("%v/%v/%v/%v", c.String("dir"), c.String("owner"), c.String("repo"), Filename)
	list, err := Deserialize(path)
	if err != nil {
		return fmt.Errorf("deserialize: %v", err)
	}

	format := strings.ToLower(c.String("format"))
	if err := print(format, list); err != nil {
		return fmt.Errorf("print: %v", err)
	}

	return nil
}

func Deserialize(path string) ([]github.WorkflowJob, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("file not found: %v", path)
	}

	read, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read %s: %v", path, err)
	}

	jobs := make([]github.WorkflowJob, 0)
	for _, r := range strings.Split(string(read), "\n") {
		if len(r) < 1 {
			// skip empty line
			continue
		}

		var job github.WorkflowJob
		if err := json.Unmarshal([]byte(r), &job); err != nil {
			return nil, fmt.Errorf("unmarshal: %v", err)
		}

		jobs = append(jobs, job)
	}

	return jobs, nil
}

func print(format string, list []github.WorkflowJob) error {
	if format == "json" {
		for _, r := range list {
			fmt.Println(JSON(r))
		}

		return nil
	}

	if format == "csv" {
		fmt.Println("job_id, job_name, status, conclusion, started_at, completed_at, duration(minutes)")

		for _, r := range list {
			fmt.Println(CSV(r))
		}

		return nil
	}

	return fmt.Errorf("invalid format=%v", format)
}

func CSV(j github.WorkflowJob) string {
	return fmt.Sprintf(
		"%v, %v, %v, %v, %v, %v, %v",
		*j.ID,
		*j.Name,
		*j.Status,
		*j.Conclusion,
		j.StartedAt.Format("2006-01-02 15:04:05"),
		j.CompletedAt.Format("2006-01-02 15:04:05"),
		j.CompletedAt.Sub(j.StartedAt.Time).Minutes(),
	)
}