package core

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

type Generator struct {
	CsvPath    string `json:"csv_path"`
	PdfPath    string `json:"pdf_path"`
	TextHeader string `json:"text_header"`
	EanHeader  string `json:"ean_header"`
}

// Check if generator is valid.
func (g *Generator) Validate() error {
	{
		ext := filepath.Ext(g.CsvPath)
		if strings.ToLower(ext) != ".csv" {
			return errors.New("Error: Input file must have a .csv extension")
		}
	}
	{
		ext := filepath.Ext(g.PdfPath)
		if strings.ToLower(ext) != ".pdf" {
			return errors.New("Error: Input file must have a .pdf extension")
		}
	}
	return nil
}

// If pdf path is not set, update it to valid one.
func (g *Generator) UpdatePdfPath() {
	if g.PdfPath != "" {
		return
	}
	g.PdfPath = GeneratePdfPath(g.CsvPath)
}

func GeneratePdfPath(path string) string {
	ext := filepath.Ext(path)
	return strings.TrimSuffix(filepath.Base(path), ext) + ".pdf"
}

// Load generator from file.
func LoadGenerator(path string) (*Generator, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, errors.Join(errors.New("Cannot load generator"), err)
	}
	defer file.Close()

	var g Generator
	if err := json.NewDecoder(file).Decode(&g); err != nil {
		return nil, err
	}
	return &g, nil
}

// Save generator to file.
func (g *Generator) Save(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	return encoder.Encode(g)
}

// Generate pdf from actual instance of generator.
func (g *Generator) GeneratePdf() error {
	err := g.Validate()
	if err != nil {
		return err
	}

	//Open the CSV file
	file, err := os.Open(g.CsvPath)
	if err != nil {
		return err
	}
	defer file.Close()

	records, err := RecordsFromCsv(file, g.TextHeader, g.EanHeader)
	if err != nil {
		return err
	}
	pdf := NewPdf()
	pdf.AddPages(records)
	err = pdf.Save(g.PdfPath)
	if err != nil {
		return err
	}
	return nil
}
