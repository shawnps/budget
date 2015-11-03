package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

var (
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

	for _, trans := range b.Transactions {
		b.Remaining = b.Total + trans.Cost
	}
	return b, nil
}

func main() {
	flag.Parse()
	b, err := parseFile(fmt.Sprintf("budget/%s.txt", *month))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Total:", b.Total)
	fmt.Println("Remaining:", b.Remaining)
}
