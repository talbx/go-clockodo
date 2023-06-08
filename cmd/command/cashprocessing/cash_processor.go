package cashprocessing

import (
	"strconv"
	"strings"

	"github.com/Rhymond/go-money"
	i "github.com/talbx/go-clockodo/pkg/intercept"
	"github.com/talbx/go-clockodo/pkg/model"
	"github.com/talbx/go-clockodo/pkg/util"
)

func RevenueToANRevenue(revenue *money.Money) {

	margin := i.ClockodoConfig.Revenue.Margin
	salary := i.ClockodoConfig.Revenue.Salary

	minRevenue := money.NewFromFloat((float64((margin + salary) / 4)), money.EUR)
	myRevenue, err := revenue.Subtract(minRevenue)

	if err != nil {
		util.SugaredLogger.Error("There was an error calculating the AN revenue", err)
	}

	util.SugaredLogger.Infof("Your actual revenue is %v after margin of %v was subtracted from raw revenue of %v", myRevenue.Display(), minRevenue.Display(), revenue.Display())

}

func CashProcess(dbc *model.DayByCustomer) {
	hm := strings.Split(dbc.RoundedTime, ":")

	revenue := i.ClockodoConfig.Revenue
	rate := revenue.HourlyRate
	h, err := strconv.Atoi(hm[0])
	if err != nil {
		util.SugaredLogger.Error("There was an error calculating the revenue", err)
		dbc.AggregatedRevenue = money.New(0, money.EUR)
		return
	}

	m, err := strconv.ParseFloat("0."+hm[1], 64)
	if err != nil {
		util.SugaredLogger.Error("There was an error calculating the revenue", err)
		dbc.AggregatedRevenue = money.New(0, money.EUR)
		return
	}

	hours := float64(h) * float64(rate)
	mins := m * float64(rate)
	sum := money.NewFromFloat(hours, money.EUR)
	rest := money.NewFromFloat(mins, money.EUR)

	total, err := sum.Add(rest)

	if err != nil {
		util.SugaredLogger.Error("There was an error translating the revenue into currency", err)
		dbc.AggregatedRevenue = money.New(0, money.EUR)
		return
	}
	dbc.AggregatedRevenue = total
}
