package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/Fanteria/EANBaker/app"
	"github.com/Fanteria/EANBaker/core"
)

func main() {
	var err error
	if len(os.Args) == 1 {
		err = app.RunGui()
	} else {
		generator, err := GetOpts()
		if err == nil {
			err = generator.GeneratePdf()
		}
	}
	if err != nil {
		log.Fatal(err)
	}
}

const USAGE string = `Usage:
  eanbaker [flags]

Description:
  This application reads a CSV file, extracts text and EAN code columns by their headers, and generates a PDF file containing barcodes. Each barcode in the PDF is accompanied by the corresponding text.

  If no flags are provided, the application starts in GUI mode.

Flags:
`

func GetOpts() (*core.Generator, error) {
	// Define flags
	generator := core.Generator{}
	flag.StringVar(&generator.CsvPath, "csv", "data.csv", "Path to the input data in CSV.")
	flag.StringVar(&generator.PdfPath, "pdf", "", "Path to the generated pdf file. If is not set, CSV file path with suffix changed to pdf is used.")
	flag.StringVar(&generator.TextHeader, "text-header", "Material Number", "Case insensitive header of column that will be used as text.")
	flag.StringVar(&generator.EanHeader, "ean-header", "ean", "Case insensitive header of column containing ean codes that will be used to generate barcode.")
	var comma_string string
	flag.StringVar(&comma_string, "csv-separator", ",", "CSV file column separator.")

	flag.Usage = func() {
		fmt.Print(USAGE)
		flag.PrintDefaults()
	}

	// Parse the flags
	flag.Parse()

	comma, err := core.CommaFromString(comma_string)
	if err != nil {
		return nil, err
	}
	generator.CsvComma = comma

	generator.UpdatePdfPath()

	// Check if opts are valid.
	err = generator.Validate()
	if err != nil {
		return nil, err
	}

	// Use the flag values
	return &generator, nil
}
