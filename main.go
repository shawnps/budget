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
)

var (
	dir   = flag.String("d", "budget", "budget data directory")
	month = flag.String("m", "", "year/month in format YYYYMM")
)

type Transaction struct {
	Cost float64
	Name string
}

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

// modified from https://groups.google.com/forum/#!topic/golang-nuts/FT7cjmcL7gw
// Pair is a data structure to hold a key/value pair.
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

func main() {
	flag.Parse()
	b, err := parseFile(fmt.Sprintf("%s/%s.txt", *dir, *month))
	if err != nil {
		log.Fatal(err)
	}
	top := map[string]float64{}
	fmt.Println("Total:", b.Total)
	fmt.Println("Remaining:", b.Remaining)
	for _, t := range b.Transactions {
		top[t.Name] += t.Cost
	}
	fmt.Println("Top costs:")
	pl := sortMapByValue(top)
	for _, p := range pl {
		fmt.Printf("%s:%f\n", p.Key, p.Value)
	}
}
