package core

import (
	"errors"
	"fmt"
	"image/png"
	"os"
	"strings"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/ean"
)

type Record struct {
	Text string
	Ean  string
}

func RecordsFromTable(table [][]string, text string, ean string) ([]Record, error) {
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
	text_lower := strings.ToLower(text)
	ean_lower := strings.ToLower(ean)
	for i, item := range table[0] {
		if strings.ToLower(item) == text_lower {
			text_index = i
		} else if strings.ToLower(item) == ean_lower {
			ean_index = i
		}
	}
	// Check if headers was found
	if text_index == -1 {
		return nil, errors.New(fmt.Sprintf("Cannot find text header '%s'", text))
	}
	if ean_index == -1 {
		return nil, errors.New(fmt.Sprintf("Cannot find ean header '%s'", ean))
	}

	// Print each record
	ret := []Record{}
	for _, csv_line := range table[1:] {
		if csv_line[ean_index] != "" {
			ret = append(ret, Record{Text: csv_line[text_index], Ean: csv_line[ean_index]})
		}
	}
	return ret, nil
}

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
