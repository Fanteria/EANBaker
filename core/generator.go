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

// Converts a string to a Comma rune type.
// Returns 0 for empty strings, or the first rune for single-character strings.
// Returns an error for multi-character strings.
func CommaFromString(s string) (Comma, error) {
	if len(s) == 0 {
		return 0, nil
	}
	if len(s) != 1 {
		return 0, fmt.Errorf("expected a single character, got %q", s)
	}
	return Comma([]rune(s)[0]), nil
}

// Implements the json.Marshaler interface for Comma.
// Converts the Comma rune to a JSON string representation.
func (r Comma) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(r))
}

// Implements the json.Unmarshaler interface for Comma.
// Converts a JSON string to a Comma rune, validating it's a single character.
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

// Validate checks if the generator configuration is valid.
// Verifies that CSV input file has .csv extension and PDF
// output file has .pdf extension.
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

// Sets the PDF output path if it's not already configured.
// Generates a PDF path based on the CSV input path by changing the extension.
func (g *Generator) UpdatePdfPath() {
	if g.PdfPath != "" {
		return
	}
	g.PdfPath = GeneratePdfPath(g.CsvPath)
}

// Creates a PDF filename from an input file path.
// Extracts the base filename and replaces the extension with .pdf.
func GeneratePdfPath(path string) string {
	ext := filepath.Ext(path)
	return strings.TrimSuffix(filepath.Base(path), ext) + ".pdf"
}

// Reads and deserializes a generator configuration from a JSON file.
// Returns a pointer to the loaded Generator or an error if loading fails.
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

// Save serializes and writes the generator configuration to a JSON file.
// Creates the file with proper indentation for readability.
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

// Processes the CSV file and generates a PDF with barcodes.
// Validates configuration, reads the CSV file, extracts records, and creates a PDF
// with the specified number of repetitions for each barcode.
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
