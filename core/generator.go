package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Comma rune

func CommaFromString(s string) (Comma, error) {
	if len(s) == 0 {
		return 0, nil
	}
	if len(s) != 1 {
		return 0, fmt.Errorf("expected a single character, got %q", s)
	}
	return Comma([]rune(s)[0]), nil
}

func (r Comma) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(r))
}

func (r *Comma) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	comma, err := CommaFromString(s)
	if err != nil {
		return err
	}
	if len(s) != 1 {
		return fmt.Errorf("expected a single character, got %q", s)
	}
	*r = comma
	return nil
}

type Generator struct {
	CsvPath      string `json:"csv_path"`
	PdfPath      string `json:"pdf_path"`
	CsvComma     Comma  `json:"csv_comma"`
	TextHeader   string `json:"text_header"`
	EanHeader    string `json:"ean_header"`
	TimesEachEAN uint   `json:"times_each_ean"`
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
	encoder.SetIndent("", "  ")
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

	csv, err := TableFromCsv(file, rune(g.CsvComma))
	if err != nil {
		return err
	}

	records, err := RecordsFromTable(csv, g.TextHeader, g.EanHeader)
	if err != nil {
		return err
	}
	pdf := NewPdf()
	pdf.AddPages(records, g.TimesEachEAN)
	err = pdf.Save(g.PdfPath)
	if err != nil {
		return err
	}
	return nil
}
