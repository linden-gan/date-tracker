/*
Copyright © 2025 ganlinden@gmail.com
*/
package cmd

import (
	"path/filepath"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"gopkg.in/yaml.v3"
)

// var DATA_PATH = "example.yaml"
var DATA_PATH = filepath.Join(os.Getenv("HOME"), "go/bin/data/tasks.yaml")
const MAX_DAYS = 50
const MAX_FUTURE_DAYS = 10

func Unmarshal() (*map[string][]*YamlTask, []string) {
	data, err := os.ReadFile(DATA_PATH)
	if err != nil {
		panic(err)
	}

	yamlTask := make(map[string][]*YamlTask)
	if err := yaml.Unmarshal(data, &yamlTask); err != nil {
		panic(err)
	}

	names := []string{}
	for group, tasks := range yamlTask {
		names = append(names, group)
		for _, task := range tasks {
			names = append(names, task.Alias...)
		}
	}

	return &yamlTask, names
}

func Marshal(yamlTask *map[string][]*YamlTask) {
	// enc := yaml.NewEncoder(os.Stdout)
    // enc.SetIndent(2)
    // defer enc.Close()

    // enc.Encode(output)
	data, err := yaml.Marshal(yamlTask)
	err = os.WriteFile(DATA_PATH, data, 0644)
	if err != nil {
		panic(err)
	}
}

type YamlTask struct {
	Name string `yaml:"name"` // unique
	Dates []string `yaml:"dates,flow"` // [12/25, 12/30, 1/2, 1/6]
	Alias []string `yaml:"alias,flow"`
}

type Task struct {
	Name string
	Dates []time.Time
	Alias []string
}

func Query(query map[string]struct{}, yamlTask *map[string][]*YamlTask) *map[string][]*YamlTask {
	res := make(map[string][]*YamlTask)
	for groupName, tasks := range *yamlTask {
		if _, ok := query[groupName]; ok {
			res[groupName] = append(res[groupName], tasks...)
		} else {
			for _, task := range tasks {
				if _, ok := query[task.Name]; ok {
					res[groupName] = append(res[groupName], task)
				} else {
					for _, alias := range task.Alias {
						if _, ok := query[alias]; ok {
							res[groupName] = append(res[groupName], task)
						}
					}
				}
			}
		}
	}
	return &res
}

func Check(yamlTask *map[string][]*YamlTask) {
	today := date2String(time.Now())
	for _, tasks := range *yamlTask {
		for _, task := range tasks {
			if len(task.Dates) == 0 || task.Dates[len(task.Dates) - 1] != today {
				task.Dates = append(task.Dates, today)
				fmt.Println("Checking", task.Name)
			} else {
				fmt.Println("Already checked", task.Name, "today.")
			}
		}
	}
}

func YamlTasks2Tasks(yamlTask *map[string][]*YamlTask) *map[string][]*Task {
	res := make(map[string][]*Task)
	for groupName, yamlTasks := range *yamlTask {
		for _, yamlTask := range yamlTasks {
			dates := strings2Dates(yamlTask.Dates)
			res[groupName] = append(res[groupName], &Task{
				Name: yamlTask.Name,
				Dates: dates,
			})
		}
	}
	return &res
}

func Tasks2YamlTasks(task *map[string][]*Task) *map[string][]*YamlTask {
	res := make(map[string][]*YamlTask)
	for groupName, tasks := range *task {
		for _, task := range tasks {
			dates := dates2Strings(task.Dates)
			res[groupName] = append(res[groupName], &YamlTask{
				Name: task.Name,
				Dates: dates,
				Alias: task.Alias,
			})
		}
	}
	return &res
}

func Print(YamlTask *map[string][]*YamlTask) {
	ts := YamlTasks2Tasks(YamlTask)
	var tw = tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	// Print the header like - - * - - - - - - * - - - - - - |
	var sb strings.Builder
	sb.WriteString(strings.Repeat(" -", (MAX_DAYS - 1) % 7))
	sb.WriteString(strings.Repeat(" * - - - - - -", (MAX_DAYS - 1) / 7))
	sb.WriteString(" |")
	for groupName, tasks := range *ts {
		fmt.Fprintln(tw, groupName)
		fmt.Fprintln(tw, sb.String(), "\t", "Item")
		for _, task := range tasks {
			fmt.Fprintln(tw, dates2Tracker(task.Dates), "\t", task.Name)
		}
		fmt.Fprintln(tw, "")
	}
	tw.Flush()
}

func dates2Tracker(dates []time.Time) string {
	maxLength := MAX_DAYS * 2
	var sb strings.Builder
	// Print date tracker like . . . # . # # . . . .   # means checked.
	laterDate := time.Now().AddDate(0, 0, 1)
	for i := len(dates) - 1; i >= 0 && sb.Len() <= maxLength; i-- {
		gap := countDays(dates[i], laterDate)
		if gap == 0 {
			panic("dates2Tracker: duplicate dates")
		}
		sb.WriteString(strings.Repeat(". ", gap - 1))
		sb.WriteString("# ")
		laterDate = dates[i]
	}
	// If there are still days to fill, fill with dots.
	if sb.Len() < maxLength {
		sb.WriteString(strings.Repeat(". ", (maxLength - sb.Len()) / 2))
	}
	tracker := []byte(sb.String())
	tracker = tracker[:maxLength]
	// Reverse the tracker to make it in the right order.
	for i, j := 0, len(tracker) - 1; i < j; i, j = i + 1, j - 1 {
		tracker[i], tracker[j] = tracker[j], tracker[i]
	}
	return string(tracker)
}
