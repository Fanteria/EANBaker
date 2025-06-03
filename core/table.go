package core

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"

	"github.com/xuri/excelize/v2"
)

type Table [][]string

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
