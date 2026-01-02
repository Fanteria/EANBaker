package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
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
	TimesHeader  string `json:"times_header"`
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
func LoadGenerator(path string, log *slog.Logger) (*Generator, error) {
	file, err := os.Open(path)
	if err != nil {
		err = errors.Join(errors.New("Cannot load generator"), err)
		log.Error("Failed to open generator file", "err", err)
		return nil, err
	}
	defer file.Close()

	var g Generator
	if err := json.NewDecoder(file).Decode(&g); err != nil {
		log.Error("Failed to decode generator file", "err", err)
		return nil, err
	}
	log.Info("Generator loaded", "generator", g)
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

// Generator must be valid
func (g *Generator) GenerateFromTable(table Table, log *slog.Logger) error {
	records, err := RecordsFromTable(table, g.TextHeader, g.EanHeader, g.TimesHeader)
	if err != nil {
		log.Error("Failed to get records from table", "err", err)
		return err
	}
	log.Debug("Records in table", "records", records)
	pdf := NewPdf()
	pdf.AddPages(records, g.TimesEachEAN, log)
	err = pdf.Save(g.PdfPath)
	if err != nil {
		log.Error("Failed to save pdf file", "err", err)
		return err
	}
	return nil
}

func (g *Generator) Generate(filename string, content io.ReadSeeker, log *slog.Logger) error {
	log.Debug("Try to generate pdf", "filename", filename, "generator", *g)
	err := g.Validate()
	if err != nil {
		log.Error("Generator is invalid", "err", err)
		return err
	}
	log.Info("Generator is valid")

	var table Table
	switch strings.ToLower(filepath.Ext(filename)) {
	case "csv":
		table, err = TableFromCsv(content, rune(g.CsvComma))
	case "xlsx":
		table, err = TableFromExcel(content, 0)
	default:
		table, err = TableFromCsv(content, rune(g.CsvComma))
		if err != nil {
			content.Seek(0, io.SeekStart)
			table, err = TableFromExcel(content, 0)
		}
	}
	log.Debug("Table to generate pdf", "table", table)
	return g.GenerateFromTable(table, log)
}
