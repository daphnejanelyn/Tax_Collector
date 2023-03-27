package main

import (
	"fmt"
	"strconv"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

// Define a struct to hold user input
type TaxInputs struct {
	Income        float64
	Deductions    float64
	TaxableIncome float64
	Tax           float64
}

func main() {
	// Create a new application
	myApp := app.New()

	// Create the input widgets
	incomeEntry := widget.NewEntry()
	incomeEntry.SetPlaceHolder("Enter your income")

	deductionsEntry := widget.NewEntry()
	deductionsEntry.SetPlaceHolder("Enter your deductions")

	// Create the calculate button
	calculateButton := widget.NewButton("Calculate", func() {
		// Parse user input
		incomeStr := incomeEntry.Text
		deductionsStr := deductionsEntry.Text

		// Convert input to float64
		income, err := strconv.ParseFloat(incomeStr, 64)
		if err != nil {
			widget.NewLabel("Invalid income input")
			return
		}

		deductions, err := strconv.ParseFloat(deductionsStr, 64)
		if err != nil {
			widget.NewLabel("Invalid deductions input")
			return
		}

		taxLabel := widget.NewLabel("")
		// Calculate the taxable income and tax
		taxableIncome := income - deductions
		tax := calculateTax(taxableIncome)

		// Create a TaxInputs struct to hold the user input and calculated tax
		inputs := TaxInputs{
			Income:        income,
			Deductions:    deductions,
			TaxableIncome: taxableIncome,
			Tax:           tax,
		}

		// Display the results
		taxString := fmt.Sprintf("Tax: Php %.2f", inputs.Tax)
		taxLabel.SetText(taxString)
	})

	// Create the layout
	content := container.New(layout.NewVBoxLayout(),
		widget.NewLabel("Income"),
		incomeEntry,
		widget.NewLabel("Deductions"),
		deductionsEntry,
		calculateButton,
	)

	// Create the window
	myWindow := myApp.NewWindow("Tax Calculator")
	myWindow.SetContent(content)
	myWindow.ShowAndRun()
}

// Define a function to calculate tax based on taxable income
/*func calculateTax(taxableIncome float64) float64 {
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
}*/

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
