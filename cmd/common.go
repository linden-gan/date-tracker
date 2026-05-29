/*
Copyright © 2025 ganlinden@gmail.com
*/
package cmd

import (
	"path/filepath"
	"fmt"
	"math"
	"os"
	"sort"
	"strings"
	"text/tabwriter"
	"time"

	"gopkg.in/yaml.v3"
)

// var DATA_PATH = "test.yaml"
var DATA_PATH = filepath.Join(os.Getenv("HOME"), "go/bin/data/tasks.yaml")
const MAX_DAYS = 60

type YamlTask struct {
	Name string `yaml:"name"` // unique
	Dates []string `yaml:"dates,flow"` // [12/25, 12/30, 1/2, 1/6]
	Alias []string `yaml:"alias,flow"`
}

type CliTask struct {
	Name string
	Dates []time.Time
	Alias []string
	CountDown int
}

func Unmarshal() ([]*YamlTask, []string) {
	data, err := os.ReadFile(DATA_PATH)
	if err != nil {
		panic(err)
	}

	yamlTasks := []*YamlTask{}
	if err := yaml.Unmarshal(data, &yamlTasks); err != nil {
		panic(err)
	}

	names := []string{}
	for _, yamlTask := range yamlTasks {
		names = append(names, yamlTask.Name)
		names = append(names, yamlTask.Alias...)
	}

	return yamlTasks, names
}

func Marshal(yamlTasks []*YamlTask) {
	data, err := yaml.Marshal(yamlTasks)
	err = os.WriteFile(DATA_PATH, data, 0644)
	if err != nil {
		panic(err)
	}
}

func Query(query map[string]struct{}, yamlTasks []*YamlTask) []*YamlTask {
	if len(query) == 0 {
		return yamlTasks
	}

	res := []*YamlTask{}
	for _, yamlTask := range yamlTasks {
		if _, ok := query[yamlTask.Name]; ok {
			res = append(res, yamlTask)
		} else {
			for _, alias := range yamlTask.Alias {
				if _, ok := query[alias]; ok {
					res = append(res, yamlTask)
					break
				}
			}
		}
	}
	return res
}

func Check(yamlTasks []*YamlTask) {
	today := date2String(time.Now())
	for _, yamlTask := range yamlTasks {
		if len(yamlTask.Dates) == 0 ||
				yamlTask.Dates[len(yamlTask.Dates) - 1] != today {
			yamlTask.Dates = append(yamlTask.Dates, today)
			fmt.Println("Checked", yamlTask.Name)
		} else {
			fmt.Println("Already checked", yamlTask.Name, "today.")
		}
	}
}

func yamlTask2CliTask(yamlTask *YamlTask) *CliTask {
	dates := strings2Dates(yamlTask.Dates)
	cliTask := CliTask {
		Name: yamlTask.Name,
		Dates: dates,
		Alias: yamlTask.Alias,
		CountDown: countDown(dates),
	}
	return &cliTask
}

func yamlTasks2CliTasks(yamlTasks []*YamlTask) []*CliTask {
	cliTasks := make([]*CliTask, len(yamlTasks))
	for i, yamlTask := range yamlTasks {
		cliTasks[i] = yamlTask2CliTask(yamlTask)
	}
	return cliTasks
}

func countDown(dates []time.Time) int {
	if len(dates) < 2 {
		// Return 0 if we don't have enough data to deduce countDown.
		return 0
	}
	frequency := float64(countDays(dates[len(dates) - 2], dates[len(dates) - 1]))
	if len(dates) > 2 {
		frequency = 0.7 * frequency + 0.3 * float64(countDays(dates[len(dates) - 3], dates[len(dates) - 2]))
	}
	countDown := frequency - float64(countDays(dates[len(dates) - 1], time.Now()))
	return int(math.Round(countDown))
}

func Print(yamlTasks []*YamlTask, urgency bool) {
	cliTasks := yamlTasks2CliTasks(yamlTasks)

	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	// Print the header like - - | - - - - - - | - - - - - - |
	trackerLen := MAX_DAYS * 2
	sb := strings.Builder{}
	// sb.Grow(trackerLen)
	sb.WriteString(strings.Repeat(" -", (MAX_DAYS - 1) % 7))
	sb.WriteString(strings.Repeat(" | - - - - - -", (MAX_DAYS - 1) / 7))
	sb.WriteString(" |")
	fmt.Fprintln(tw, sb.String(), "\t", "Next", "\t", "Task")
	sb.Reset()

	if urgency {
		sort.Slice(cliTasks, func(i, j int) bool {
			return cliTasks[i].CountDown < cliTasks[j].CountDown
		})
	}
    
	for _, cliTask := range cliTasks {
		laterDate := time.Now().AddDate(0, 0, 1) // sentinel
		for i := len(cliTask.Dates) - 1; i >= 0 && sb.Len() <= trackerLen; i-- {
			curDate := cliTask.Dates[i]
			gap := countDays(curDate, laterDate)
			if gap <= 0 || gap >= MAX_DAYS {
				fmt.Printf(
					"WARNING: suspicious date %s in %s.\n",
					date2String(curDate), cliTask.Name)
			}
			sb.WriteString(strings.Repeat(". ", gap - 1))
     		sb.WriteString("# ")
			laterDate = curDate
		}
		if sb.Len() < trackerLen {
			sb.WriteString(strings.Repeat(". ", (trackerLen - sb.Len()) / 2))
		}
		// Reverse the string.
		tracker := []byte(sb.String())[:trackerLen]
		for i, j := 0, trackerLen - 1; i < j; i, j = i + 1, j - 1 {
			tracker[i], tracker[j] = tracker[j], tracker[i]
		}
		fmt.Fprintln(tw, string(tracker), "\t", cliTask.CountDown, "\t", cliTask.Name)
		sb.Reset()
	}

	tw.Flush()
}
