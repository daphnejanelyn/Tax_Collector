package main

import (
	"fmt"
	//"image/color"

	//"image/color"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// Define a struct to hold user input
type TaxInputs struct {
	MonthlyIncome           float64
	TaxableIncome           float64
	Tax                     float64
	NetPayAfterTax          float64
	SSSContributions        float64
	PhilHealthContributions float64
	PagIbigContributions    float64
	TotalContributions      float64
	TotalDeductions         float64
	NetPayAfterDeductions   float64
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
	netPayAfterTaxLabel := widget.NewLabel("")
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

		// Convert input to float64
		monthlyIncome, err := strconv.ParseFloat(incomeStr, 64)
		if err != nil {
			widget.NewLabel("Invalid monthly income input")
			return
		}

		// Calculate the tax, net pay after tax, monthly contributions, total deductions and net pay after deductions
		taxableIncome := monthlyIncome * 12 // INSERT RIGHT FORMULA
		tax := calculateTax(taxableIncome)
		netPayAfterTax := monthlyIncome - (tax / 12) // INSERT RIGHT FORMULA
		sssContributions := calculateSSSContributions(monthlyIncome)
		philhealthContributions := calculatePhilHealthContributions(monthlyIncome)
		pagibigContributions := calculatePagIbigContributions(monthlyIncome)
		totalContributions := sssContributions + philhealthContributions + pagibigContributions
		totalDeductions := tax/12 + totalContributions // INSERT RIGHT FORMULA
		netPayAfterDeductions := monthlyIncome - totalDeductions

		// Create a TaxInputs struct to hold the user input and calculated values
		inputs := TaxInputs{
			MonthlyIncome:           monthlyIncome,
			TaxableIncome:           taxableIncome,
			Tax:                     tax,
			NetPayAfterTax:          netPayAfterTax,
			SSSContributions:        sssContributions,
			PhilHealthContributions: philhealthContributions,
			PagIbigContributions:    pagibigContributions,
			TotalContributions:      totalContributions,
			TotalDeductions:         totalDeductions,
			NetPayAfterDeductions:   netPayAfterDeductions,
		}

		// Display the results
		taxLabel.SetText(fmt.Sprintf("Income Tax: Php %.2f", inputs.Tax))
		taxableIncomeLabel.SetText(fmt.Sprintf("Taxable Income: Php %.2f", inputs.TaxableIncome))
		netPayAfterTaxLabel.SetText(fmt.Sprintf("Net Pay After Tax: Php %.2f", inputs.NetPayAfterTax))
		sssContributionsLabel.SetText(fmt.Sprintf("Php %.2f", inputs.SSSContributions))
		philhealthContributionsLabel.SetText(fmt.Sprintf("Php %.2f", inputs.PhilHealthContributions))
		pagibigContributionsLabel.SetText(fmt.Sprintf("Php %.2f", inputs.PagIbigContributions))
		totalContributionsLabel.SetText(fmt.Sprintf("Php %.2f", inputs.TotalContributions))
		totalDeductionsLabel.SetText(fmt.Sprintf("Php %.2f", inputs.TotalDeductions))
		netPayAfterDeductionsLabel.SetText(fmt.Sprintf("Php %.2f", inputs.NetPayAfterDeductions))
	})
	//calculateBtn.SetBackgroundColor(color.RGBA{R: 211, G: 211, B: 211, A: 211})

	// Create the layout
	// Create the layout
	taxContainer := container.NewVBox(
		widget.NewLabel("Tax Computation"),
		taxableIncomeLabel,
		taxLabel,
		netPayAfterTaxLabel,
		layout.NewSpacer())
	contribContainer := container.NewVBox(
		widget.NewLabel("Monthly Contributions"),
		container.NewHBox(
			widget.NewLabel("SSS Contribution"),
			sssContributionsLabel,
		),
		container.NewHBox(
			widget.NewLabel("Philheath Contribution"),
			philhealthContributionsLabel,
		),
		container.NewHBox(
			widget.NewLabel("PagIbig Contribution"),
			pagibigContributionsLabel,
		),
		container.NewHBox(
			widget.NewLabel("Total Contribution"),
			totalContributionsLabel,
		))

	content := container.New(layout.NewVBoxLayout(),
		container.NewVBox(
			widget.NewLabel("Monthly Income"),
			incomeEntry,
			calculateBtn,
			layout.NewSpacer(),
		),
		container.New(layout.NewGridWrapLayout(fyne.NewSize(300, 300)), contribContainer, taxContainer),
		container.NewVBox(
			widget.NewLabel("Total Deductions"),
			totalDeductionsLabel,
			widget.NewLabel("Net Pay After Deductions Label"),
			netPayAfterDeductionsLabel,
		),
	)
	// Create the window
	myWindow := myApp.NewWindow("Tax Calculator")
	myApp.Settings().SetTheme(theme.LightTheme())
	myWindow.SetContent(content)
	myWindow.Resize(fyne.NewSize(800, 600))
	myWindow.ShowAndRun()

}

// Define a function to calculate tax based on taxable income
func calculateTax(taxableIncome float64) float64 {
	var tax float64

	if taxableIncome <= 250000 {
		tax = 0
	} else if taxableIncome <= 400000 {
		tax = (taxableIncome - 250000) * 0.20
	} else if taxableIncome <= 800000 {
		tax = 30000 + (taxableIncome-400000)*0.25
	} else if taxableIncome <= 2000000 {
		tax = 130000 + (taxableIncome-800000)*0.30
	} else if taxableIncome <= 8000000 {
		tax = 490000 + (taxableIncome-2000000)*0.32
	} else {
		tax = 2410000 + (taxableIncome-8000000)*0.35
	}

	return tax
}

func calculateSSSContributions(taxableIncome float64) float64 {
	// insert code
	return 0
}

func calculatePagIbigContributions(taxableIncome float64) float64 {
	// insert code
	return 0
}

func calculatePhilHealthContributions(taxableIncome float64) float64 {
	// insert code
	return 0
}

func calculateNetPayAfterDeductions(taxableIncome float64) float64 {
	// insert code
	return 0
}
