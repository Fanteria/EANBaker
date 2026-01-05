package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/Fanteria/EANBaker/app"
	"github.com/Fanteria/EANBaker/core"
	"github.com/Fanteria/EANBaker/values"
)

func main() {
	err := func() error {
		logger, err := core.MultiLoggerFromEnv()
		if err != nil {
			return err
		}
		if len(os.Args) == 1 {
			return app.RunGui(logger)
		} else {
			generator, err := GetOpts()
			if err == nil {
				//Open the CSV file
				file, err := os.Open(generator.CsvPath)
				if err != nil {
					return err
				}
				defer file.Close()
				return generator.Generate(generator.CsvPath, file, logger.Logger)
			}
			return nil
		}
	}()
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

// Parses command-line flags and creates a configured Generator instance.
// Defines and parses all CLI options including CSV path, PDF path, headers, and barcode repetition.
// Validates the configuration before returning the generator.
func GetOpts() (*core.Generator, error) {
	// Define flags
	generator := core.Generator{}
	flag.StringVar(&generator.CsvPath, "csv", "data.csv", "Path to the input data in CSV.")
	flag.StringVar(&generator.PdfPath, "pdf", "", "Path to the generated pdf file. If is not set, CSV file path with suffix changed to pdf is used.")
	flag.StringVar(&generator.TextHeader, "text-header", "Material Number", "Case insensitive header of column that will be used as text.")
	flag.StringVar(&generator.EanHeader, "ean-header", "ean", "Case insensitive header of column containing ean codes that will be used to generate barcode.")
	flag.StringVar(&generator.TimesHeader, "times-header", "", `Name of the column that specifies how many times each EAN code should be generated. If the column is empty, each EAN code is generated only once.
If the column contains a number, the EAN code is generated that many times. Rows are processed line by line, so identical EANs appear consecutively.`)
	flag.UintVar(&generator.TimesEachEAN, "times-each-ean", 1, "Number of times each EAN code will be printed in the output PDF.")
	comma_string := flag.String("csv-separator", ",", "CSV file column separator.")
	print_version := flag.Bool("version", false, "Print version information and exit")

	flag.Usage = func() {
		fmt.Print(USAGE)
		flag.PrintDefaults()
	}

	// Parse the flags
	flag.Parse()

	if *print_version {
		fmt.Println("Version:", values.Version)
		fmt.Println("Commit:", values.Commit)
		fmt.Println("Build date:", values.Date)
		os.Exit(0)
	}

	comma, err := core.CommaFromString(*comma_string)
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
