package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"
)

var (
	dir       = flag.String("d", "budget", "budget data directory")
	yearMonth = flag.String("m", "", "year/month in format YYYYMM")
	short     = flag.Bool("short", false, "use tagged results when available")
)

// Transaction is a single transaction with a cost and name
type Transaction struct {
	Cost float64
	Name string
}

// Budget is a monthly budget
type Budget struct {
	Total        float64
	Remaining    float64
	TagMap       map[string]string
	Transactions []Transaction
}

func parseFile(dir string, year int, month time.Month) (Budget, error) {
	name := fmt.Sprintf("%s/%d%02d.txt", dir, year, month)
	file, err := os.Open(name)
	if err != nil {
		return Budget{}, err
	}

	r := bufio.NewReader(file)
	b := Budget{}
	// first line is the total for the month
	line, _, err := r.ReadLine()
	if err != nil {
		return Budget{}, err
	}

	total, err := strconv.ParseFloat(string(line), 64)
	if err != nil {
		return Budget{}, err
	}

	b.Total = total
	tagMap := map[string]string{}

	for {
		line, _, err := r.ReadLine()
		if err == io.EOF {
			break
		}

		if err != nil {
			return Budget{}, err
		}

		if strings.HasPrefix(string(line), "#") {
			sp := strings.Split(string(line), ":")
			if len(sp) != 2 {
				return Budget{}, fmt.Errorf("invalid tag line %q", string(line))
			}

			tn := strings.TrimPrefix(sp[0], "# ")
			items := strings.Split(sp[1], ",")

			for _, item := range items {
				tagMap[strings.TrimSpace(item)] = tn
			}

			continue
		}

		t := strings.SplitN(string(line), " ", 2)
		if strings.TrimSpace(string(line)) == "" {
			continue
		}

		if len(t) != 2 {
			return Budget{}, fmt.Errorf("invalid line %q", line)
		}

		cost, err := strconv.ParseFloat(t[0], 64)
		if err != nil {
			return Budget{}, err
		}

		trans := Transaction{Cost: cost, Name: strings.TrimSpace(t[1])}

		b.Transactions = append(b.Transactions, trans)
	}

	tagMapFilename := filepath.Join(dir, "tags.txt")
	fileTagMap, err := parseTagMapFile(tagMapFilename)
	if os.IsNotExist(err) {
		// tag file doesn't exist, ignore
	} else if err != nil {
		return Budget{}, err
	}

	b.TagMap = mergeTagMaps(tagMap, fileTagMap)
	b.Remaining = b.Total
	for _, trans := range b.Transactions {
		b.Remaining -= trans.Cost
	}

	return b, nil
}

// Pair is a data structure to hold a key/value pair.
// modified from https://groups.google.com/forum/#!topic/golang-nuts/FT7cjmcL7gw
type Pair struct {
	Key   string
	Value float64
}

// PairList is a slice of Pairs that implements sort.Interface to sort by Value.
type PairList []Pair

func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }

// A function to turn a map into a PairList, then sort and return it.
func sortMapByValue(m map[string]float64) PairList {
	p := make(PairList, len(m))
	i := 0

	for k, v := range m {
		p[i] = Pair{k, v}
		i++
	}

	sort.Sort(sort.Reverse(p))
	return p
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

func parseTagMapFile(filename string) (map[string]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return map[string]string{}, err
	}

	tagMap := map[string]string{}

	r := bufio.NewReader(file)

	for {
		line, _, err := r.ReadLine()
		if err == io.EOF {
			break
		}

		if err != nil {
			return map[string]string{}, err
		}

		if strings.TrimSpace(string(line)) == "" {
			continue
		}

		sp := strings.Split(string(line), ":")
		if len(sp) != 2 {
			return map[string]string{}, fmt.Errorf("invalid tag line %q", string(line))
		}

		tn := sp[0]
		items := strings.Split(sp[1], ",")

		for _, item := range items {
			tagMap[strings.TrimSpace(item)] = tn
		}
	}

	return tagMap, nil
}

func mergeTagMaps(tm, ftm map[string]string) map[string]string {
	for k, v := range ftm {
		if tm[k] == "" {
			tm[k] = v
		}
	}

	return tm
}

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

	b, err := parseFile(*dir, year, month)
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
