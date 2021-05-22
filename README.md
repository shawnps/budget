[![Go Report Card](https://goreportcard.com/badge/github.com/shawnps/budget)](https://goreportcard.com/report/github.com/shawnps/budget)

# budget
Simple command-line budget tool.

## How to use
`go get github.com/shawnps/budget`

`mkdir budget`

Filenames must use format `YYYYMM.txt`:

`budget/202105.txt`

Files must be in following format:

```
1000
5.50 sushi
10.50 sushi
15.00 sushi
8.00 book
3.00 coffee
10.00 book
8.00 sushi
4.00 coffee
3.00 coffee
```

The top number is the total you want to spend per month. Each line after is something you bought.

Either call the command from one directory outside of `budget/`, or provide the directory with `-d`.

```
$ budget -m 202105
Total:         1000.00
Remaining:     933.00
Remaining/day: 42.41
Top costs:
    sushi:  39.00
    book:   18.00
    coffee: 10.00
```

## Tags

You can also create tags:

```
1000

# food: sushi, coffee
# books: book

5.50 sushi
10.50 sushi
15.00 sushi
8.00 book
3.00 coffee
10.00 book
8.00 sushi
4.00 coffee
3.00 coffee
```

then pass the `-short` flag:

```
$ budget -m 202105 -short
Total:         1000.00
Remaining:     933.00
Remaining/day: 42.41
Costs:
    food:  49.00
    books: 18.00
```

You can also specify tags in a `tags.txt` file:

```
food: sushi, coffee
books: book
```

## You might also like
[Ledger](http://www.ledger-cli.org/index.html)

[YNAB](https://www.youneedabudget.com/)

[Personal Finance Reddit](http://personalfinance.reddit.com/)
