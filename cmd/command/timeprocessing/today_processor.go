package timeprocessing

import (
	"fmt"
	. "github.com/talbx/go-clockodo/pkg/model"
	. "github.com/talbx/go-clockodo/pkg/util"
	"golang.design/x/clipboard"
	"log/slog"
	"strings"
	"time"
)

type TodayProcessor struct{}

func (p TodayProcessor) Process(mode string, last int) error {
	now := time.Now()
	y := now.Year()
	m := now.Month()
	d := now.Day()
	start := time.Date(y, m, d, 0, 0, 0, 0, time.Local)
	s := SOB(start).Format("2006-01-02T15:04:05Z")
	e := EOB(now).Format("2006-01-02T15:04:05Z")

	var entriesRoot = "v2/entries"

	var repo TimeEntriesResponse
	query := fmt.Sprintf("%s?time_since=%s&time_until=%s", entriesRoot, s, e)
	_, err := CallApi(query, &repo)

	if err != nil {
		return err
	}
	desc := make([]string, 0)
	for _, entry := range repo.Entries {
		desc = append(desc, entry.Description)
	}

	temp := strings.Join(desc, ",")
	tasksOfTheDay := strings.ReplaceAll(temp, ",,", ",")

	slog.Info(fmt.Sprintf("your tasks of the day were %v.", tasksOfTheDay))
	slog.Info("the tasks were added to your clipboard!")

	err = clipboard.Init()
	if err != nil {
		return err
	}

	clipboard.Write(clipboard.FmtText, []byte(tasksOfTheDay))

	return nil
}
