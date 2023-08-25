package render

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/talbx/go-clockodo/pkg/model"
	. "github.com/talbx/go-clockodo/pkg/model"
	"github.com/talbx/go-clockodo/pkg/util"
)

func Render(mappy map[string][]DayByCustomer, clock model.ClockResponse, clockProcessor func(clock *model.ClockResponse) (int, int)) {

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"#", "ID", "Customer", "Tasks", "Times"})
	rowConfigAutoMerge := table.RowConfig{AutoMerge: true}

	taskCount := 0
	tt := 0
	keys := make([]string, 0, len(mappy))
	for k := range mappy {
		keys = append(keys, k)
	}

	for _, key := range keys {
		for _, entry := range mappy[key] {
			alterTime(&entry)
			taskCount += len(strings.Split(entry.AggregatedTasks, ","))
			tt += entry.TotalTime
			t.AppendRow(table.Row{key, entry.CustomerId, entry.Customer, entry.AggregatedTasks, fmt.Sprintf("(%v) - %v", entry.RoundedTime, entry.AggregatedTime)}, rowConfigAutoMerge)
			t.AppendSeparator()
		}
	}

	t.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, AutoMerge: true},
		{Number: 2, AutoMerge: true},
		{Number: 3, AutoMerge: true},
		{Number: 4, AutoMerge: true},
		{Number: 5, AutoMerge: true},
		{Number: 6, AutoMerge: true},
		{Number: 7, AutoMerge: true},
	})

	th, tm := util.DurationToHM(tt)
	t.AppendSeparator()
	t.SetStyle(table.StyleLight)
	t.AppendFooter(table.Row{"TOTAL", "", "", fmt.Sprintf("Total tasks: %v", taskCount), fmt.Sprintf("%v:%v", util.AddLeadingZero(th), util.AddLeadingZero(tm))})
	t.Render()

	h, m := clockProcessor(&clock)
	slog.Info(fmt.Sprintf("Also, you have a task running for %vh:%vm right now.\n", h, m))
}

func alterTime(entry *DayByCustomer) {
	var hs, m = 0, 0
	if strings.Contains(entry.AggregatedTime, "h") {
		hRest := strings.Split(entry.AggregatedTime, "h")
		mRest := strings.Split(hRest[1], "m")
		hs, _ = strconv.Atoi(hRest[0])
		m, _ = strconv.Atoi(mRest[0])

	} else if strings.Contains(entry.AggregatedTime, "m") {
		mRest := strings.Split(entry.AggregatedTime, "m")
		m, _ = strconv.Atoi(mRest[0])
	}
	r1, r2 := util.Round(hs, m)
	entry.RoundedTime = fmt.Sprintf("%v:%v", r1, r2)
}
