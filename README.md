# budget
super simple budget thing

## how to use
`go get github.com/shawnps/budget`

`mkdir budget`

filenames must use format YYYYMM.txt, for example:

`budget/201510.txt`

files must be in following format:

```
100000
-550,sushi
-320,train
-8000,shoes
```

the top number is the total you want to spend per month. the lines after are things you have bought.

if everything is setup properly,

`budget -m 201510`

should show:

```
➜  ~  budget -m 201510
Total: 100000
Remaining: 92000
```

## you might also like
http://www.ledger-cli.org/index.html
https://www.youneedabudget.com/
http://personalfinance.reddit.com/
