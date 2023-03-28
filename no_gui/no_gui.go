package main

import (
	"fmt"
	"github.com/shopspring/decimal"
)

func calculateTax(taxableIncome decimal.Decimal) decimal.Decimal {
	var brackets = []decimal.Decimal{
		decimal.NewFromFloat(20833),
		decimal.NewFromFloat(33333),
		decimal.NewFromFloat(66667),
		decimal.NewFromFloat(166667),
		decimal.NewFromFloat(666667),
	}

	var rates = []decimal.Decimal{
		decimal.NewFromFloat(0.0),
		decimal.NewFromFloat(0.15),
		decimal.NewFromFloat(0.20),
		decimal.NewFromFloat(0.25),
		decimal.NewFromFloat(0.30),
		decimal.NewFromFloat(0.35),
	}

	tax := decimal.Zero

	for i := 1; i < len(brackets); i++ {
		if taxableIncome.LessThanOrEqual(brackets[i]) {
			fmt.Println(i, ":")
			tax = tax.Add((taxableIncome.Sub(brackets[i-1])).Mul(rates[i]))
			break
		} else {
			tax = tax.Add(brackets[i].Sub(brackets[i-1]).Mul(rates[i]))
			fmt.Println("\nADDED: ", brackets[i].Sub(brackets[i-1]).Mul(rates[i]))
		}
	}

	return tax.Round(2)
}

func main() {
	fmt.Println(calculateTax(decimal.NewFromFloat(20833)))
}