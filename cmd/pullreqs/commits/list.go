package commits

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/urfave/cli/v2"
)

func List(c *cli.Context) error {
	path := fmt.Sprintf("%v/%v/%v/%v", c.String("dir"), c.String("owner"), c.String("repo"), Filename)
	list, err := Deserialize(path)
	if err != nil {
		return fmt.Errorf("deserialize: %v", err)
	}

	sort.Slice(list, func(i, j int) bool { return list[i].Commit.Author.Date.Unix() < list[i].Commit.Author.Date.Unix() }) // desc

	format := strings.ToLower(c.String("format"))
	if err := print(format, list); err != nil {
		return fmt.Errorf("print: %v", err)
	}

	return nil
}

func Deserialize(path string) ([]CommitWithPRID, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("file not found: %v", path)
	}

	read, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read %s: %v", path, err)
	}

	out := make([]CommitWithPRID, 0)
	for _, r := range strings.Split(string(read), "\n") {
		if len(r) < 1 {
			// skip empty line
			continue
		}

		var commit CommitWithPRID
		if err := json.Unmarshal([]byte(r), &commit); err != nil {
			return nil, fmt.Errorf("unmarshal: %v", err)
		}

		out = append(out, commit)
	}

	return out, nil
}

func print(format string, list []CommitWithPRID) error {
	if format == "json" {
		for _, r := range list {
			fmt.Println(JSON(r))
		}

		return nil
	}

	if format == "csv" {
		fmt.Println("pr_id, pr_number, sha, login, date, message, ")

		for _, r := range list {
			fmt.Println(r.CSV())
		}

		return nil
	}

	return fmt.Errorf("invalid format=%v", format)
}
