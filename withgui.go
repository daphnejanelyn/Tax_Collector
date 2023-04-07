package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/leekchan/accounting"
	"github.com/shopspring/decimal"
)

/*
Define a structure for all variables needed

	to hold computational results and inputs
*/
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
	/* Create a new application along
	with its output and input widgets */
	myApp := app.New()

	// Input Widgets
	incomeEntry := widget.NewEntry()

	// Output Widgets
	taxLabel := widget.NewLabel("")
	taxableIncomeLabel := widget.NewLabel("")
	sssContributionsLabel := widget.NewLabel("")
	pagibigContributionsLabel := widget.NewLabel("")
	philhealthContributionsLabel := widget.NewLabel("")
	totalContributionsLabel := widget.NewLabel("")
	totalDeductionsLabel := widget.NewLabel("")
	netPayAfterDeductionsLabel := widget.NewLabel("")

	// Create input row for user input
	incomeEntry.SetPlaceHolder("Enter your monthly income")

	// Create the calculate button
	calculateBtn := widget.NewButton("Calculate", func() {

		// Parse user input
		incomeStr := incomeEntry.Text

		// Convert input to decimal format
		monthlyIncome, err := decimal.NewFromString(incomeStr)

		// Function to cross check for invalid inputs
		if err != nil || monthlyIncome.LessThan(decimal.Zero) {
			widget.NewLabel("Invalid monthly income input")
			return
		}

		// Calling functions to calculate for monthly contributions
		sssContributions := calculateSSSContributions(monthlyIncome)
		philhealthContributions := calculatePhilHealthContributions(monthlyIncome)
		pagibigContributions := calculatePagIbigContributions(monthlyIncome)
		totalContributions := decimal.Sum(sssContributions,
			philhealthContributions,
			pagibigContributions)

		// Calling functions to calculate for tax deductions
		taxableIncome := monthlyIncome.Sub(totalContributions)
		tax := calculateTax(taxableIncome)
		totalDeductions := totalContributions.Add(tax)
		netPayAfterDeductions := monthlyIncome.Sub(totalDeductions)

		// Create a display struct to hold the user input and calculated values
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

		/* Display the results of computation in Peso format
		with 2 digit precision for decimal points */

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

	/* Container for tax computations (i.e., taxable income and income tax) */
	taxContainer := container.NewVBox(
		widget.NewLabelWithStyle("Tax Computation",
			fyne.TextAlignLeading,
			fyne.TextStyle{Bold: true}),
		container.NewHBox(
			widget.NewLabel("Taxable Income\t\t"),
			taxableIncomeLabel,
		),
		container.NewHBox(
			widget.NewLabel("Income Tax\t\t\t"),
			taxLabel,
		),
		layout.NewSpacer())

	/* Container for monthly contributions computations (i.e., SSS, PagIbig and Philheath) */
	contribContainer := container.NewVBox(
		widget.NewLabelWithStyle("Monthly Contributions",
			fyne.TextAlignLeading,
			fyne.TextStyle{Bold: true}),
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

	/* Container for display of final computations on net pay and total deductions */
	finalComputations := container.NewVBox(
		widget.NewLabelWithStyle("Total Deductions",
			fyne.TextAlignLeading,
			fyne.TextStyle{Bold: true}),
		totalDeductionsLabel,
		widget.NewLabelWithStyle("Net Pay After Deductions",
			fyne.TextAlignLeading,
			fyne.TextStyle{Bold: true}),
		netPayAfterDeductionsLabel,
	)

	/* Container for display of prompt box for user's input on monthly contribution */
	content := container.New(layout.NewVBoxLayout(),
		container.NewVBox(
			widget.NewLabelWithStyle("Monthly Income",
				fyne.TextAlignLeading,
				fyne.TextStyle{Bold: true}),
			incomeEntry,
			calculateBtn,
			layout.NewSpacer(),
		),

		/* Resizing tax and contributions container through Grid Wrap Layout Manager */
		container.New(layout.NewGridWrapLayout(fyne.NewSize(300, 200)),
			contribContainer,
			taxContainer),
		container.New(layout.NewGridWrapLayout(fyne.NewSize(150, 150)),
			finalComputations),
	)

	// Create a new window for the desktop application
	myWindow := myApp.NewWindow("Tax Calculator")
	myApp.Settings().SetTheme(theme.LightTheme())
	myWindow.SetContent(content)
	myWindow.Resize(fyne.NewSize(600, 600))
	myWindow.SetFixedSize(true)
	myWindow.ShowAndRun()

}

// Define a function to calculate tax based on taxable income
func calculateTax(taxableIncome decimal.Decimal) decimal.Decimal {
	// Calculaed from the 2022 Tax Table

	var brackets = []decimal.Decimal{
		decimal.NewFromFloat(0),
		decimal.NewFromFloat(20833),
		decimal.NewFromFloat(33333),
		decimal.NewFromFloat(66667),
		decimal.NewFromFloat(166667),
		decimal.NewFromFloat(666667),
	}

	var rates = []decimal.Decimal{
		decimal.NewFromFloat(0.0),
		decimal.NewFromFloat(0.20),
		decimal.NewFromFloat(0.25),
		decimal.NewFromFloat(0.30),
		decimal.NewFromFloat(0.32),
		decimal.NewFromFloat(0.35),
	}

	var bases = []decimal.Decimal{
		decimal.NewFromFloat(0.0),
		decimal.NewFromFloat(0.0),
		decimal.NewFromFloat(2500.0),
		decimal.NewFromFloat(10833.33),
		decimal.NewFromFloat(40833.33),
		decimal.NewFromFloat(200833.33),
	}

	tax := decimal.Zero

	i := 1
	for i < 6 && taxableIncome.GreaterThanOrEqual(brackets[i]) {
		i++
	}
	i--

	tax = bases[i].Add((taxableIncome.Sub(brackets[i])).Mul(rates[i]))
	fmt.Println()

	return tax.Round(2)
}

func calculateSSSContributions(monthlyIncome decimal.Decimal) decimal.Decimal {
	/* Notice that based on the 2021-2022 SSS Table,
	the formula to get the gross contribution based on the monthly income is
	the nearest multiple of 500, except when it is lower than 3250 it is automatically 3000,
	and when it is greater than or equal to 24750 it is automatically 25000.
	This is then multiplied by 4.5% to get the employee's actual SSS contribution

	The following code is the implementation of this
	*/
	employeeRate := decimal.NewFromFloat(0.045)

	var gross decimal.Decimal
	if monthlyIncome.LessThan(decimal.NewFromInt(3250)) {
		gross = decimal.NewFromInt(3000)
	} else if monthlyIncome.GreaterThanOrEqual(decimal.NewFromInt(24750)) {
		gross = decimal.NewFromInt(25000)
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
	/* The 2021 Pag-Ibig contribution table takes the monthly income and
	multiplies it by 1% if it is less than or equal to 1500,
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
	/* The 2022 contribution rate for Philhealth is 4.0%
	which is split equally between the employee and employer.

	People with <= 10000 salary must contribute 200.
	People with >= 80000 salary must contribute 1600

	*/
	rate := decimal.NewFromFloat(0.02)
	min_threshold := decimal.NewFromFloat(10000)
	max_threshold := decimal.NewFromFloat(80000)

	if monthlyIncome.LessThan(min_threshold) {
		monthlyIncome = min_threshold
	} else if monthlyIncome.GreaterThan(max_threshold) {
		monthlyIncome = max_threshold
	}

	return monthlyIncome.Mul(rate)
}
