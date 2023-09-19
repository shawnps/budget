package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/shawnps/budget/budget"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var (
	dir            = flag.String("d", "budget", "budget data directory")
	yearMonth      = flag.String("m", "", "year/month in format YYYYMM")
	yearMonthRange = flag.String("r", "", "year/month range in format YYYYMM-YYYYMM")
	short          = flag.Bool("short", false, "use tagged results when available")
	tag            = flag.String("tag", "", "get results for a given tag")
)

func main() {
	flag.Parse()

	var (
		year  = time.Now().Year()
		month = time.Now().Month()
		err   error
	)

	if *yearMonth != "" && *yearMonthRange != "" {
		log.Fatal("can only use -m or -r")
	}

	if ym := *yearMonth; ym != "" {
		parsed, err := parseYearMonth(ym)
		if err != nil {
			log.Fatal(err)
		}

		year = parsed.year
		month = parsed.month
	}

	name := fmt.Sprintf("%s/%d%02d.txt", *dir, year, month)
	file, err := os.Open(name)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	r := bufio.NewReader(file)

	var tm *bufio.Reader
	tagMapFilename := filepath.Join(*dir, "tags.txt")
	fileTagMap, err := os.Open(tagMapFilename)
	if os.IsNotExist(err) {
		tm = bufio.NewReader(strings.NewReader(""))
		// tag file doesn't exist, ignore
	} else if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if fileTagMap != nil {
		tm = bufio.NewReader(fileTagMap)
	}

	if *yearMonthRange != "" {
		yearMonths, err := timeRange(*yearMonthRange)
		if err != nil {
			log.Fatal(err)
		}

		var all []budget.Budget
		for _, ym := range yearMonths {
			name := fmt.Sprintf("%s/%d%02d.txt", *dir, ym.year, ym.month)
			file, err := os.Open(name)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			r := bufio.NewReader(file)
			b, err := budget.Parse(r, tm)
			if err != nil {
				log.Fatal(err)
			}

			all = append(all, b)
		}

		combined := budget.CombineBudgets(all)
		printBudget(combined, 0, 0)
		return
	}

	b, err := budget.Parse(r, tm)
	if err != nil {
		log.Fatal(err)
	}

	printBudget(b, year, month)
}

// borrowed from golang src time pkg
// daysBefore[m] counts the number of days in a non-leap year
// before month m begins.  There is an entry for m=12, counting
// the number of days before January of next year (365).
var daysBefore = [...]int32{
	0,
	31,
	31 + 28,
	31 + 28 + 31,
	31 + 28 + 31 + 30,
	31 + 28 + 31 + 30 + 31,
	31 + 28 + 31 + 30 + 31 + 30,
	31 + 28 + 31 + 30 + 31 + 30 + 31,
	31 + 28 + 31 + 30 + 31 + 30 + 31 + 31,
	31 + 28 + 31 + 30 + 31 + 30 + 31 + 31 + 30,
	31 + 28 + 31 + 30 + 31 + 30 + 31 + 31 + 30 + 31,
	31 + 28 + 31 + 30 + 31 + 30 + 31 + 31 + 30 + 31 + 30,
	31 + 28 + 31 + 30 + 31 + 30 + 31 + 31 + 30 + 31 + 30 + 31,
}

func isLeap(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}

func daysIn(m time.Month, year int) int {
	if m == time.February && isLeap(year) {
		return 29
	}
	return int(daysBefore[m] - daysBefore[m-1])
}

// pair is a data structure to hold a key/value pair.
// modified from https://groups.google.com/forum/#!topic/golang-nuts/FT7cjmcL7gw
type pair struct {
	Key   string
	Value float64
}

// pairList is a slice of pairs that implements sort.Interface to sort by Value.
type pairList []pair

func (p pairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p pairList) Len() int           { return len(p) }
func (p pairList) Less(i, j int) bool { return p[i].Value < p[j].Value }

// A function to turn a map into a pairList, then sort and return it.
func sortMapByValue(m map[string]float64) pairList {
	p := make(pairList, len(m))
	i := 0

	for k, v := range m {
		p[i] = pair{k, v}
		i++
	}

	sort.Sort(sort.Reverse(p))
	return p
}

func printBudget(b budget.Budget, year int, month time.Month) {
	mp := message.NewPrinter(language.English)
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 0, ' ', 0)

	mp.Fprintf(w, "Total:\t %.2f\n", b.Total)
	mp.Fprintf(w, "Remaining:\t %.2f\n", b.Remaining)

	if year >= time.Now().Year() && month >= time.Now().Month() {
		dr := daysIn(month, year) - time.Now().Day() + 1
		rpd := b.Remaining / float64(dr)
		mp.Fprintf(w, "Remaining/day:\t %.2f\n", rpd)
	} else {
		mp.Fprintf(w, "Spent:\t %.2f\n", b.Total-b.Remaining)
	}

	top := map[string]float64{}
	for _, t := range b.Transactions {
		tmEntry, inTagMap := b.TagMap[t.Name]

		if *tag != "" {
			if tmEntry == *tag {
				top[t.Name] += t.Cost
			}

			continue
		}

		if !*short || !inTagMap {
			top[t.Name] += t.Cost
			continue
		}

		top[b.TagMap[t.Name]] += t.Cost
	}

	fmt.Fprintf(w, "Costs:\n")
	pl := sortMapByValue(top)
	for _, p := range pl {
		mp.Fprintf(w, "    %s:\t %.2f\n", p.Key, p.Value)
	}

	w.Flush()
}

type budgetMonth struct {
	year  int
	month time.Month
}

func parseYearMonth(ym string) (budgetMonth, error) {
	year, err := strconv.Atoi(ym[0:4])
	if err != nil {
		return budgetMonth{}, err
	}

	monthStr, err := strconv.Atoi(ym[4:6])
	if err != nil {
		return budgetMonth{}, err
	}

	month := time.Month(monthStr)

	return budgetMonth{year, month}, nil
}

func timeRange(timeRange string) ([]budgetMonth, error) {
	var bm []budgetMonth

	sp := strings.Split(timeRange, "-")
	if len(sp) != 2 {
		return []budgetMonth{}, fmt.Errorf("invalid time range %q", timeRange)
	}

	first, second := sp[0], sp[1]

	if len(first) != 6 || len(second) != 6 {
		return []budgetMonth{}, fmt.Errorf("invalid time range %q", timeRange)
	}

	parsedFirst, err := parseYearMonth(first)
	if err != nil {
		return []budgetMonth{}, err
	}

	parsedSecond, err := parseYearMonth(second)
	if err != nil {
		return []budgetMonth{}, err
	}

	if parsedSecond.year < parsedFirst.year {
		return []budgetMonth{}, errors.New("second year must be > first year")
	}

	if parsedFirst.year == parsedSecond.year {
		for m := parsedFirst.month; m <= parsedSecond.month; m++ {
			bm = append(bm, budgetMonth{parsedFirst.year, m})
		}

		return bm, nil
	}

	for m := parsedFirst.month; m <= 12; m++ {
		bm = append(bm, budgetMonth{parsedFirst.year, m})
	}

	if parsedSecond.year-parsedFirst.year == 1 {
		for m := time.Month(1); m <= parsedSecond.month; m++ {
			bm = append(bm, budgetMonth{parsedSecond.year, m})
		}

		return bm, nil
	}

	for y := parsedFirst.year + 1; y <= parsedSecond.year; y++ {
		var endMonth time.Month = time.Month(12)
		if parsedSecond.year == y {
			endMonth = parsedSecond.month
		}

		for m := time.Month(1); m <= endMonth; m++ {
			bm = append(bm, budgetMonth{y, m})
		}
	}

	return bm, nil
}
