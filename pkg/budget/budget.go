package budget

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
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

// Parse parses a single month's budget
func Parse(r, tm *bufio.Reader) (Budget, error) {
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

		costStr := strings.ReplaceAll(t[0], ",", "")
		cost, err := strconv.ParseFloat(costStr, 64)
		if err != nil {
			return Budget{}, err
		}

		trans := Transaction{Cost: cost, Name: strings.TrimSpace(t[1])}

		b.Transactions = append(b.Transactions, trans)
	}

	fileTagMap, err := parseTagMapFile(tm)
	if err != nil {
		return b, err
	}

	b.TagMap = mergeTagMaps(tagMap, fileTagMap)
	b.Remaining = b.Total
	for _, trans := range b.Transactions {
		b.Remaining -= trans.Cost
	}

	return b, nil
}

func parseTagMapFile(r *bufio.Reader) (map[string]string, error) {
	tagMap := map[string]string{}

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
