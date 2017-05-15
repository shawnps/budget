[![Go Report Card](https://goreportcard.com/badge/github.com/shawnps/budget)](https://goreportcard.com/report/github.com/shawnps/budget)

# budget
super simple budget thing

## how to use
`go get github.com/shawnps/budget`

`mkdir budget`

filenames must use format `YYYYMM.txt`, for example:

`budget/201510.txt`

files must be in following format:

```
100000
-550,sushi
-320,train
-8000,shoes
-1050,sushi
-3000,sushi
-800,book
-300,coffee
-1000,book
-500,sushi
```

the top number is the total you want to spend per month. the lines after are things you have bought.

you must either call the command from one directory outside of `budget/`, or provide the directory with `-d`

if everything is setup properly,

`budget -m 201510`

should show:

```
➜  ~  budget -m 201510
Total: 100000
Remaining: 84480
Remaining/day: 16896.00
Top costs:
	shoes: -8000.00
	sushi: -5100.00
	book: -1800.00
	train: -320.00
	coffee: -300.00
```

## you might also like
[Ledger](http://www.ledger-cli.org/index.html)

[YNAB](https://www.youneedabudget.com/)

[Personal Finance Reddit](http://personalfinance.reddit.com/)
