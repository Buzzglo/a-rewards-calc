# A-rewards-calc

I followed this guide to learn how to deal with CSV's with goLang
https://earthly.dev/blog/golang-csv-files/

This is my first go program outside of crash course I did years ago but I worked with Node, Angular and python as part
of my day to day. Please dont judge me.

##

The purpose of this project is to calculate the rewards until today and project how much more can be spent or not at the rewards partners to
optimise the return. 

Should I buy Meds, Toilet paper or Milk from rewards parter A, now or tomorrow. If the price is the same between partner A and B but I have 
capped my spend at partner A get it at partner B. etc.

The effect of checking this has allowed me to recover 15% of my credit card spend in total absa Rewards or about 10% of my take home every month.

## Scope

I have not investigated API support from the bank I use, however they do allow CSV export. Rewards are calculated from the 16th
of the previous month to the 15th of the current month and awarded on the 1st of the next Month. I export the CSV and rename the file cheque.csv or credit.csv . The file is picked up from the root and uses the adjusted date in the code. ( I need to add a function that looks up the current date and calculates based on where I an now, possibly parameterise for last months caluclations and a input param on the run) 

It basic and assumes the top tier. Functionality to be extended and UI mattured incrementally.

## How to run

```
go mod init a-rewards-calc
```

running the my-rewards-calc
```
go run my-rewards-calc.go
```

testing commit
