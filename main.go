package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

var (
	dir   = flag.String("d", "budget", "budget data directory")
	month = flag.String("m", "", "year/month in format YYYYMM")
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
	Transactions []Transaction
}

func parseFile(name string) (Budget, error) {
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
	for {
		line, _, err := r.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			return Budget{}, err
		}
		t := strings.Split(string(line), ",")
		if len(t) != 2 {
			return Budget{}, fmt.Errorf("invalid line %q", line)
		}
		cost, err := strconv.ParseFloat(t[0], 64)
		if err != nil {
			return Budget{}, err
		}
		trans := Transaction{Cost: cost, Name: t[1]}

		b.Transactions = append(b.Transactions, trans)
	}

	b.Remaining = b.Total
	for _, trans := range b.Transactions {
		b.Remaining = b.Remaining + trans.Cost
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
	sort.Sort(p)
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

func main() {
	flag.Parse()
	b, err := parseFile(fmt.Sprintf("%s/%s.txt", *dir, *month))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Total:", b.Total)
	fmt.Println("Remaining:", b.Remaining)
	mon := *month
	year, err := strconv.Atoi(mon[0:4])
	if err != nil {
		log.Fatal(err)
	}
	m, err := strconv.Atoi(mon[4:6])
	if err != nil {
		log.Fatal(err)
	}
	rpd := b.Remaining / float64(daysIn(time.Month(m), year)-time.Now().Day())
	fmt.Printf("Remaining/day: %.2f\n", rpd)
	top := map[string]float64{}
	for _, t := range b.Transactions {
		top[t.Name] += t.Cost
	}
	fmt.Println("Top costs:")
	pl := sortMapByValue(top)
	for _, p := range pl {
		fmt.Printf("\t%s: %.2f\n", p.Key, p.Value)
	}
}
