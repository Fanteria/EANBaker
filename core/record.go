package core

import (
	"errors"
	"fmt"
	"image/png"
	"os"
	"strconv"
	"strings"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/ean"
)

type Record struct {
	Text  string
	Ean   string
	Times int
}

// Extracts Record structures from a 2D string table using column headers.
// Finds the specified text and EAN columns (case-insensitive) and creates records for each row.
// Skips rows with empty EAN values. Returns an error if headers are not found or table is empty.
func RecordsFromTable(table [][]string, text string, ean string, times string) ([]Record, error) {
	if text == "" {
		return nil, errors.New("Text column header cannot be empty")
	}
	if ean == "" {
		return nil, errors.New("Ean column header cannot be empty")
	}
	if len(table) == 0 {
		return nil, errors.New("Table with data cannot be empty")
	}

	// Find headers
	text_index := -1
	ean_index := -1
	times_index := -1
	text_lower := strings.ToLower(text)
	ean_lower := strings.ToLower(ean)
	times_lower := strings.ToLower(times)
	for i, item := range table[0] {
		switch strings.ToLower(item) {
		case text_lower:
			text_index = i
		case ean_lower:
			ean_index = i
		case times_lower:
			times_index = i
		}
	}
	// Check if headers was found
	if text_index == -1 {
		return nil, fmt.Errorf("Cannot find text header '%s'", text)
	}
	if ean_index == -1 {
		return nil, fmt.Errorf("Cannot find ean header '%s'", ean)
	}
	if times_index == -1 && strings.TrimSpace(times) != "" {
		return nil, fmt.Errorf("Cannot find times header '%s'", times)
	}

	// Print each record
	ret := []Record{}
	for _, csv_line := range table[1:] {
		if csv_line[ean_index] != "" {
			times_value := 1
			if times_index != -1 {
				times_str := strings.TrimSpace(csv_line[times_index])
				value_float, err := strconv.ParseFloat(times_str, 0)
				if err != nil {
					value_int, err := strconv.ParseInt(times_str, 10, 0)
					if err != nil {
						return nil, errors.Join(errors.New("Times column contain non numeric string"), err)
					}
					times_value = int(value_int)
				} else {
					times_value = int(value_float)
				}
			}
			ret = append(ret, Record{
					Text: csv_line[text_index],
					Ean: csv_line[ean_index],
					Times: times_value,
			})
		}
	}
	return ret, nil
}

// Creates a PNG barcode image file for the record's EAN code.
// Generates an EAN barcode, scales it to 200x200 pixels, and saves it to the specified path.
func (r *Record) GenerateBarcode(path string) error {
	// Create the barcode
	qrCode, err := ean.Encode(r.Ean)
	if err != nil {
		return err
	}

	// Scale the barcode to 200x200 pixels
	ean, err := barcode.Scale(qrCode, 200, 200)
	if err != nil {
		return err
	}

	// create the output file
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	// encode the barcode as png
	return png.Encode(file, ean)
}
