package main

import (
	"fmt"
	//"image/color"

	//"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	// for decimal
	"github.com/leekchan/accounting"
	"github.com/shopspring/decimal"
)

// Define a struct to hold user input
type TaxInputs struct {
	MonthlyIncome           decimal.Decimal
	TaxableIncome           decimal.Decimal
	Tax                     decimal.Decimal
	NetPayAfterTax          decimal.Decimal
	SSSContributions        decimal.Decimal
	PhilHealthContributions decimal.Decimal
	PagIbigContributions    decimal.Decimal
	TotalContributions      decimal.Decimal
	TotalDeductions         decimal.Decimal
	NetPayAfterDeductions   decimal.Decimal
}

func main() {
	// Create a new application
	myApp := app.New()

	// Create the input widgets
	incomeEntry := widget.NewEntry()
	incomeEntry.SetPlaceHolder("Enter your monthly income")

	// Create the output widgets
	taxLabel := widget.NewLabel("")
	taxableIncomeLabel := widget.NewLabel("")
	sssContributionsLabel := widget.NewLabel("")
	pagibigContributionsLabel := widget.NewLabel("")
	philhealthContributionsLabel := widget.NewLabel("")
	totalContributionsLabel := widget.NewLabel("")
	totalDeductionsLabel := widget.NewLabel("")
	netPayAfterDeductionsLabel := widget.NewLabel("")

	sssEntry := widget.NewEntry()
	sssEntry.SetPlaceHolder("Monthly SSS Contribution")

	// Create the calculate button
	calculateBtn := widget.NewButton("Calculate", func() {
		// Parse user input
		incomeStr := incomeEntry.Text

		// Convert input to decimal.Decimal
		monthlyIncome, err := decimal.NewFromString(incomeStr)
		if err != nil || monthlyIncome.LessThan(decimal.Zero) {
			widget.NewLabel("Invalid monthly income input")
			return
		}

		// MONTHLY CONTRIBUTIONS
		sssContributions := calculateSSSContributions(monthlyIncome)
		philhealthContributions := calculatePhilHealthContributions(monthlyIncome)
		pagibigContributions := calculatePagIbigContributions(monthlyIncome)
		totalContributions := decimal.Sum(sssContributions, philhealthContributions, pagibigContributions)

		// TAX
		taxableIncome := monthlyIncome.Sub(totalContributions)
		tax := calculateTax(taxableIncome)
		totalDeductions := totalContributions.Add(tax)
		netPayAfterDeductions := monthlyIncome.Sub(totalDeductions)

		// Create a TaxInputs struct to hold the user input and calculated values
		inputs := TaxInputs{
			MonthlyIncome:           monthlyIncome,
			TaxableIncome:           taxableIncome,
			Tax:                     tax,
			SSSContributions:        sssContributions,
			PhilHealthContributions: philhealthContributions,
			PagIbigContributions:    pagibigContributions,
			TotalContributions:      totalContributions,
			TotalDeductions:         totalDeductions,
			NetPayAfterDeductions:   netPayAfterDeductions,
		}

		// Display the results
		ac := accounting.Accounting{Symbol: "₱ ", Precision: 2}
		taxLabel.SetText(fmt.Sprintf(ac.FormatMoney(inputs.Tax)))
		taxableIncomeLabel.SetText(fmt.Sprintf(ac.FormatMoney(inputs.TaxableIncome)))
		sssContributionsLabel.SetText(fmt.Sprintf(ac.FormatMoney(inputs.SSSContributions)))
		philhealthContributionsLabel.SetText(fmt.Sprintf(ac.FormatMoney(inputs.PhilHealthContributions)))
		pagibigContributionsLabel.SetText(fmt.Sprintf(ac.FormatMoney(inputs.PagIbigContributions)))
		totalContributionsLabel.SetText(fmt.Sprintf(ac.FormatMoney(inputs.TotalContributions)))
		totalDeductionsLabel.SetText(fmt.Sprintf(ac.FormatMoney(inputs.TotalDeductions)))
		netPayAfterDeductionsLabel.SetText(fmt.Sprintf(ac.FormatMoney(inputs.NetPayAfterDeductions)))

	})
	//calculateBtn.SetBackgroundColor(color.RGBA{R: 211, G: 211, B: 211, A: 211})

	// Create the layout
	// Create the layout
	taxContainer := container.NewVBox(
		widget.NewLabelWithStyle("Tax Computation", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		container.NewHBox(
			widget.NewLabel("Taxable Income\t\t"),
			taxableIncomeLabel,
		),
		container.NewHBox(
			widget.NewLabel("Income Tax\t\t\t"),
			taxLabel,
		),
		layout.NewSpacer())
	contribContainer := container.NewVBox(
		widget.NewLabelWithStyle("Monthly Contributions", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		container.NewHBox(
			widget.NewLabel("SSS Contribution\t\t"),
			sssContributionsLabel,
		),
		container.NewHBox(
			widget.NewLabel("Philheath Contribution\t"),
			philhealthContributionsLabel,
		),
		container.NewHBox(
			widget.NewLabel("PagIbig Contribution\t"),
			pagibigContributionsLabel,
		),
		container.NewHBox(
			widget.NewLabel("Total Contribution\t\t"),
			totalContributionsLabel,
		))
	finalComputations := container.NewVBox(
		widget.NewLabelWithStyle("Total Deductions", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		totalDeductionsLabel,
		widget.NewLabelWithStyle("Net Pay After Deductions", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		netPayAfterDeductionsLabel,
	)
	content := container.New(layout.NewVBoxLayout(),
		container.NewVBox(
			widget.NewLabelWithStyle("Monthly Income", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			incomeEntry,
			calculateBtn,
			layout.NewSpacer(),
		),
		container.New(layout.NewGridWrapLayout(fyne.NewSize(300, 200)), contribContainer, taxContainer),
		container.New(layout.NewGridWrapLayout(fyne.NewSize(150, 150)), finalComputations),
	)
	// Create the window
	myWindow := myApp.NewWindow("Tax Calculator")
	myApp.Settings().SetTheme(theme.LightTheme())
	myWindow.SetContent(content)
	myWindow.Resize(fyne.NewSize(600, 600))
	myWindow.SetFixedSize(true)
	myWindow.ShowAndRun()

}

// Define a function to calculate tax based on taxable income
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

	if taxableIncome.GreaterThan(brackets[0]) {
		for i := 1; i < len(brackets); i++ {
			if taxableIncome.LessThanOrEqual(brackets[i]) {
				tax = tax.Add((taxableIncome.Sub(brackets[i-1])).Mul(rates[i]))
				//fmt.Println("i is ", i, "ADDED ", (taxableIncome.Sub(brackets[i-1])).Mul(rates[i]), "BECAME", tax)
				break
			} else {
				tax = tax.Add(brackets[i].Sub(brackets[i-1]).Mul(rates[i]))
				//fmt.Println("I is ", i, "ADDED ", brackets[i].Sub(brackets[i-1]).Mul(rates[i]), "BECAME", tax)
				if i == len(brackets)-1 {
					tax = tax.Add((taxableIncome.Sub(brackets[i])).Mul(rates[i+1]))
					//fmt.Println("This is ", i, "ADDED ", (taxableIncome.Sub(brackets[i])).Mul(rates[i+1]), "BECAME", tax)
				}

			}
		}
	}

	return tax.Round(2)
}

func calculateSSSContributions(monthlyIncome decimal.Decimal) decimal.Decimal {
	/* Notice that based on the 2023 SSS Table,
	the formula to get the gross contribution based on the monthly income is
	the nearest multiple of 500, except when it is lower than 4250 it is automatically 4000,
	and when it is greater than or equal to 29750 it is automatically 30000.
	This is then multiplied by 4.5% to get the employee's actual SSS contribution

	The following code is the implementation of this
	*/
	employeeRate := decimal.NewFromFloat(0.045)

	var gross decimal.Decimal
	if monthlyIncome.LessThan(decimal.NewFromInt(4250)) {
		gross = decimal.NewFromInt(4000)
	} else if monthlyIncome.GreaterThanOrEqual(decimal.NewFromInt(29750)) {
		gross = decimal.NewFromInt(30000)
	} else {
		// implementation of MROUND(monthlyIncome, 500)
		rate := decimal.NewFromInt(500)
		divided := monthlyIncome.Div(rate)
		floor := divided.RoundDown(0)
		ceil := divided.RoundUp(0)
		if divided.Sub(floor).LessThan(ceil.Sub(divided)) {
			gross = floor.Mul(rate)
		} else {
			gross = ceil.Mul(rate)
		}

		//fmt.Println("ROUNDED: ", gross)
	}

	return gross.Mul(employeeRate)
}

func calculatePagIbigContributions(monthlyIncome decimal.Decimal) decimal.Decimal {
	/* The https://taxcalculatorphilippines.com/ still uses the 2021 Pag-Ibig contribution table
	This takes the monthly income and multiplies it by 1% if it is less than or equal to 1500,
	otherwise it multiplies it by 2%

	The maximum pag-ibig contribution is 100.00
	*/
	var rate decimal.Decimal
	max := decimal.NewFromInt(100)

	if monthlyIncome.LessThanOrEqual(decimal.NewFromInt(1500)) {
		rate = decimal.NewFromFloat(0.01)
	} else {
		rate = decimal.NewFromFloat(0.02)
	}

	return decimal.Min(max, monthlyIncome.Mul(rate))
}

func calculatePhilHealthContributions(monthlyIncome decimal.Decimal) decimal.Decimal {
	/* The 2023 contribution rate for Philhealth is 4.5%
	which is split equally between the employee and employer.
	People have to give at least 225 and max 2025
	People with <= 10000 salary must contribute 225
	2025 is max amount anyone can contribute
	NOTE: There's a mistake on https://taxcalculatorphilippines.com/ where
	starting salary of 90000, it outputs 4050 for Philhealth instead of 2025
	*/
	rate := decimal.NewFromFloat(0.0225)
	min := decimal.NewFromFloat(225)
	// max := decimal.NewFromFloat(2025)

	if monthlyIncome.LessThanOrEqual(decimal.NewFromFloat(10000)) {
		return min
	} else if monthlyIncome.GreaterThanOrEqual(decimal.NewFromFloat(90000)) {
		return decimal.NewFromFloat(4050)
	}
	// return decimal.Min(max, monthlyIncome.Mul(rate))
	return monthlyIncome.Mul(rate)
}
