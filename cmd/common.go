/*
Copyright © 2025 ganlinden@gmail.com
*/
package cmd

import (
	// "path/filepath"
	"fmt"
	"math"
	"os"
	"sort"
	"strings"
	"text/tabwriter"
	"time"

	"gopkg.in/yaml.v3"
)

var DATA_PATH = "example.yaml"
// var DATA_PATH = filepath.Join(os.Getenv("HOME"), "go/bin/data/tasks.yaml")
const MAX_DAYS = 60

type YamlTask struct {
	Name string `yaml:"name"` // unique
	Block string `yaml:"block"`
	Dates []string `yaml:"dates,flow"` // [12/25, 12/30, 1/2, 1/6]
	Alias []string `yaml:"alias,flow"`
	Frequency int `yaml:"frequency"` // in days, optional
}

type CliTask struct {
	Name string
	Block string
	Dates []time.Time
	Alias []string
	Frequency int
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
		} else if _, ok := query[yamlTask.Block]; ok {
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
	if len(res) == 0 {
		fmt.Println("Warning: no matching tasks found.")
	}
	return res
}

func Tick(yamlTasks []*YamlTask) {
	today := date2String(time.Now())
	for _, yamlTask := range yamlTasks {
		if len(yamlTask.Dates) == 0 ||
				yamlTask.Dates[len(yamlTask.Dates) - 1] != today {
			yamlTask.Dates = append(yamlTask.Dates, today)
			fmt.Println("Ticked", yamlTask.Name)
		} else {
			fmt.Println("Already ticked", yamlTask.Name, "today.")
		}
	}
}

func Untick(yamlTasks []*YamlTask) {
	today := date2String(time.Now())
	for _, yamlTask := range yamlTasks {
		if len(yamlTask.Dates) > 0 && yamlTask.Dates[len(yamlTask.Dates) - 1] == today {
			yamlTask.Dates = yamlTask.Dates[:len(yamlTask.Dates)-1]
			fmt.Println("Unticked", yamlTask.Name)
		} else {
			fmt.Println(yamlTask.Name, "had not been ticked today.")
		}
	}
}

func yamlTask2CliTask(yamlTask *YamlTask) *CliTask {
	dates := strings2Dates(yamlTask.Dates)
	frequency := yamlTask.Frequency
	if frequency == 0 && len(dates) >= 2 { // Can deduce frequency
		frequency = countDays(dates[len(dates) - 2], dates[len(dates) - 1])
		if len(dates) > 2 {
			frequency = int(math.Round(
				0.7 * float64(frequency) +
				0.3 * float64(countDays(dates[len(dates) - 3], dates[len(dates) - 2]))))
		}
	}
	countDown := 0
	if frequency != 0 && len(dates) >= 2 {
		countDown = frequency - countDays(dates[len(dates) - 1], time.Now())
	}
	cliTask := CliTask {
		Name: yamlTask.Name,
		Block: yamlTask.Block,
		Dates: dates,
		Alias: yamlTask.Alias,
		Frequency: frequency,
		CountDown: countDown,
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
		frequency = 0.7 * frequency +
			0.3 * float64(countDays(dates[len(dates) - 3], dates[len(dates) - 2]))
	}
	countDown := frequency - float64(countDays(dates[len(dates) - 1], time.Now()))
	return int(math.Round(countDown))
}

// Sprint header like  - - | - - - - - - | - - - - - - |
func sprintHeader() string {
	return strings.Repeat(" -", (MAX_DAYS - 1) % 7) +
		strings.Repeat(" | - - - - - -", (MAX_DAYS - 1) / 7) + " |"
}

// Sprint date tracker from past dates to today like  . . . # . # . . # . . # . # . .
func sprintTracker(cliTask *CliTask) string {
	trackerLen := MAX_DAYS * 2
	sb := strings.Builder{}
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
	return string(tracker)
}

func Print(yamlTasks []*YamlTask, urgency bool, frequency bool) {
	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	cliTasks := yamlTasks2CliTasks(yamlTasks)
	header := sprintHeader()

	blocks2Tasks := make(map[string][]*CliTask)
	orderedBlocks := []string{}
	for _, cliTask := range cliTasks {
		if _, ok := blocks2Tasks[cliTask.Block]; !ok {
			orderedBlocks = append(orderedBlocks, cliTask.Block)
		}
		blocks2Tasks[cliTask.Block] = append(blocks2Tasks[cliTask.Block], cliTask)
	}
	
	for _, block := range orderedBlocks {
		tasks := blocks2Tasks[block]
		if urgency {
			sort.Slice(tasks, func(i, j int) bool {
				return tasks[i].CountDown < tasks[j].CountDown
			})
		} else if frequency {
			sort.Slice(tasks, func(i, j int) bool {
				return tasks[i].Frequency < tasks[j].Frequency
			})
		}
		
		fmt.Fprintln(tw, "Block:", block)
		fmt.Fprintln(tw, header, "\t", "Next", "\t", "Freq", "\t", "Task")
		for _, task := range tasks {
			fmt.Fprintln(tw, sprintTracker(task), "\t", task.CountDown, "\t", task.Frequency, "\t", task.Name)
		}
	}
	tw.Flush()
}
