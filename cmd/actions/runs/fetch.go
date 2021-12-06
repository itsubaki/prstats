package runs

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sort"

	"github.com/google/go-github/v40/github"
	"github.com/itsubaki/ghstats/pkg/actions/runs"
	"github.com/urfave/cli/v2"
)

const Filename = "runs.json"

func Fetch(c *cli.Context) error {
	dir := fmt.Sprintf("%v/%v/%v", c.String("dir"), c.String("owner"), c.String("repo"))
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, os.ModePerm)
	}

	path := fmt.Sprintf("%s/%s", dir, Filename)
	lastID, err := scanLastID(path)
	if err != nil {
		return fmt.Errorf("last id: %v", err)
	}

	in := runs.FetchInput{
		Owner:   c.String("owner"),
		Repo:    c.String("repo"),
		PAT:     c.String("pat"),
		Page:    c.Int("page"),
		PerPage: c.Int("perpage"),
		LastID:  lastID,
	}

	fmt.Printf("target: %v/%v\n", in.Owner, in.Repo)
	fmt.Printf("last_id: %v\n", lastID)

	ctx := context.Background()
	list, err := runs.Fetch(ctx, &in)
	if err != nil {
		return fmt.Errorf("fetch: %v", err)
	}

	if err := serialize(path, list); err != nil {
		return fmt.Errorf("serialize: %v", err)
	}

	if len(list) > 0 {
		fmt.Printf("%v %v\n", *list[len(list)-1].ID, *list[0].ID)
	}

	return nil
}

func JSON(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}

	return string(b)
}

func serialize(path string, list []*github.WorkflowRun) error {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("open file: %v", err)
	}
	defer file.Close()

	sort.Slice(list, func(i, j int) bool { return *list[i].ID < *list[j].ID }) // asc

	for _, r := range list {
		fmt.Fprintln(file, JSON(r))
	}

	return nil
}

func scanLastID(path string) (int64, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return -1, nil
	}

	file, err := os.Open(path)
	if err != nil {
		return -1, fmt.Errorf("open %v: %v", path, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lastID int64
	for scanner.Scan() {
		var run github.WorkflowRun
		if err := json.Unmarshal([]byte(scanner.Text()), &run); err != nil {
			return -1, fmt.Errorf("unmarshal: %v", err)
		}

		lastID = *run.ID
	}

	return lastID, nil
}
