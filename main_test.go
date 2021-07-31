package main

import (
	"reflect"
	"testing"
	"time"
)

func TestTimeRange(t *testing.T) {
	cases := []struct {
		in   string
		want []budgetMonth
	}{
		{
			"202105-202107",
			[]budgetMonth{
				{2021, time.Month(5)},
				{2021, time.Month(6)},
				{2021, time.Month(7)},
			},
		},
		{
			"202005-202109",
			[]budgetMonth{
				{2020, time.Month(5)},
				{2020, time.Month(6)},
				{2020, time.Month(7)},
				{2020, time.Month(8)},
				{2020, time.Month(9)},
				{2020, time.Month(10)},
				{2020, time.Month(11)},
				{2020, time.Month(12)},
				{2021, time.Month(1)},
				{2021, time.Month(2)},
				{2021, time.Month(3)},
				{2021, time.Month(4)},
				{2021, time.Month(5)},
				{2021, time.Month(6)},
				{2021, time.Month(7)},
				{2021, time.Month(8)},
				{2021, time.Month(9)},
			},
		},

		{
			"201911-202203",
			[]budgetMonth{
				{2019, time.Month(11)},
				{2019, time.Month(12)},
				{2020, time.Month(1)},
				{2020, time.Month(2)},
				{2020, time.Month(3)},
				{2020, time.Month(4)},
				{2020, time.Month(5)},
				{2020, time.Month(6)},
				{2020, time.Month(7)},
				{2020, time.Month(8)},
				{2020, time.Month(9)},
				{2020, time.Month(10)},
				{2020, time.Month(11)},
				{2020, time.Month(12)},
				{2021, time.Month(1)},
				{2021, time.Month(2)},
				{2021, time.Month(3)},
				{2021, time.Month(4)},
				{2021, time.Month(5)},
				{2021, time.Month(6)},
				{2021, time.Month(7)},
				{2021, time.Month(8)},
				{2021, time.Month(9)},
				{2021, time.Month(10)},
				{2021, time.Month(11)},
				{2021, time.Month(12)},
				{2022, time.Month(1)},
				{2022, time.Month(2)},
				{2022, time.Month(3)},
			},
		},
	}

	for _, tt := range cases {
		got, err := timeRange(tt.in)
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("timeRange(%q) =\n %v\n, want\n %v", tt.in, got, tt.want)
		}
	}
}
