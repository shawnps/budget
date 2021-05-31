package main

import (
	"bufio"
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

	"github.com/shawnps/budget/pkg/budget"
)

var (
	dir       = flag.String("d", "budget", "budget data directory")
	yearMonth = flag.String("m", "", "year/month in format YYYYMM")
	short     = flag.Bool("short", false, "use tagged results when available")
)

func main() {
	flag.Parse()

	var (
		year  = time.Now().Year()
		month = time.Now().Month()
		err   error
	)

	if ym := *yearMonth; ym != "" {
		year, err = strconv.Atoi(ym[0:4])
		if err != nil {
			log.Fatal(err)
		}

		monthStr, err := strconv.Atoi(ym[4:6])
		if err != nil {
			log.Fatal(err)
		}

		month = time.Month(monthStr)
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

	b, err := budget.Parse(r, tm)
	if err != nil {
		log.Fatal(err)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 0, ' ', 0)

	fmt.Fprintf(w, "Total:\t %.2f\n", b.Total)
	fmt.Fprintf(w, "Remaining:\t %.2f\n", b.Remaining)

	dr := daysIn(month, year) - time.Now().Day() + 1
	rpd := b.Remaining / float64(dr)
	fmt.Fprintf(w, "Remaining/day:\t %.2f\n", rpd)

	top := map[string]float64{}
	for _, t := range b.Transactions {
		_, inTagMap := b.TagMap[t.Name]

		if !*short || !inTagMap {
			top[t.Name] += t.Cost
			continue
		}

		top[b.TagMap[t.Name]] += t.Cost
	}

	fmt.Fprintf(w, "Costs:\n")
	pl := sortMapByValue(top)
	for _, p := range pl {
		fmt.Fprintf(w, "    %s:\t %.2f\n", p.Key, p.Value)
	}

	w.Flush()
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
