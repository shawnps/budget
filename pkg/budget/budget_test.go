package budget

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	file, err := os.Open("../../testdata/budget/202105.txt")
	if err != nil {
		t.Fatal(err)
	}

	var tm *bufio.Reader
	tagMapFilename := filepath.Join("../../testdata/budget/tags.txt")
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

	r := bufio.NewReader(file)

	b, err := Parse(r, tm)
	if err != nil {
		t.Fatal(err)
	}

	// {1000 933 map[book:books coffee:food sushi:food] [{5.5 sushi} {10.5 sushi} {15 sushi} {8 book} {3 coffee} {10 book} {8 sushi} {4 coffee} {3 coffee}]}

	want := Budget{
		Total:     1000,
		Remaining: 933,
		TagMap:    map[string]string{"book": "books", "coffee": "food", "sushi": "food"},
		Transactions: []Transaction{
			{5.5, "sushi"},
			{10.5, "sushi"},
			{15, "sushi"},
			{8, "book"},
			{3, "coffee"},
			{10, "book"},
			{8, "sushi"},
			{4, "coffee"},
			{3, "coffee"},
		},
	}

	if !reflect.DeepEqual(b, want) {
		t.Errorf("got Budget =\n %v\n, want\n %v", b, want)
	}
}
