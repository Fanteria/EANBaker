package core

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"

	"github.com/xuri/excelize/v2"
)

type Table [][]string

// Reads CSV data from an io.Reader and returns it as a 2D string table.
// Uses the specified comma rune as the field separator. If comma is 0,
// uses the default separator.
func TableFromCsv(r io.Reader, comma rune) (Table, error) {
	if r == nil {
		return nil, errors.New("Reader is <nil>")
	}

	// Create a new CSV reader
	reader := csv.NewReader(r)
	fmt.Printf("Comma: '%d'", comma)
	if comma != 0 {
		reader.Comma = comma
	}

	// Read all records
	csv_data, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	return Table(csv_data), nil
}

// Reads Excel data from an io.Reader and returns the first sheet as a 2D string table.
// Opens the Excel file and extracts all rows from the first available sheet.
func TableFromExcel(r io.Reader, sheet int) (Table, error) {
	exel, err := excelize.OpenReader(r)
	if err != nil {
		return nil, err
	}

	sheets := exel.GetSheetList()
	if len(sheets) == 0 {
		return nil, errors.New("Excel containing 0 sheets.")
	}
	return exel.GetRows(sheets[0])
}
