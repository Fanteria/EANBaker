package core

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRecordsFromTable(t *testing.T) {
	tests := []struct {
		name        string
		table       Table
		textHeader  string
		eanHeader   string
		timesHeader string
		wantLen     int
		wantErr     bool
	}{
		{
			name: "Valid table with data",
			table: Table{
				{"Text", "EAN"},
				{"Product A", "1234567890123"},
				{"Product B", "9876543210987"},
			},
			textHeader:  "Text",
			eanHeader:   "EAN",
			timesHeader: "",
			wantLen:     2,
			wantErr:     false,
		},
		{
			name: "Case insensitive headers",
			table: Table{
				{"TEXT", "ean"},
				{"Product", "1234567890123"},
			},
			textHeader:  "text",
			eanHeader:   "EAN",
			timesHeader: "",
			wantLen:     1,
			wantErr:     false,
		},
		{
			name: "Skip empty EAN rows",
			table: Table{
				{"Text", "EAN"},
				{"Product A", "1234567890123"},
				{"Product B", ""},
				{"Product C", "9876543210987"},
			},
			textHeader:  "Text",
			eanHeader:   "EAN",
			timesHeader: "",
			wantLen:     2,
			wantErr:     false,
		},
		{
			name:        "Empty table",
			table:       Table{},
			textHeader:  "Text",
			eanHeader:   "EAN",
			timesHeader: "",
			wantLen:     0,
			wantErr:     true,
		},
		{
			name: "Missing text header",
			table: Table{
				{"Name", "EAN"},
				{"Product", "1234567890123"},
			},
			textHeader:  "Text",
			eanHeader:   "EAN",
			timesHeader: "",
			wantLen:     0,
			wantErr:     true,
		},
		{
			name: "Missing EAN header",
			table: Table{
				{"Text", "Code"},
				{"Product", "1234567890123"},
			},
			textHeader:  "Text",
			eanHeader:   "EAN",
			timesHeader: "",
			wantLen:     0,
			wantErr:     true,
		},
		{
			name:        "Empty text header",
			table:       Table{{"Text", "EAN"}},
			textHeader:  "",
			eanHeader:   "EAN",
			timesHeader: "",
			wantLen:     0,
			wantErr:     true,
		},
		{
			name:        "Empty EAN header",
			table:       Table{{"Text", "EAN"}},
			textHeader:  "Text",
			eanHeader:   "",
			timesHeader: "",
			wantLen:     0,
			wantErr:     true,
		},
		{
			name: "With times column",
			table: Table{
				{"Text", "EAN", "Times"},
				{"Product A", "1234567890123", "2"},
				{"Product B", "9876543210987", "5"},
			},
			textHeader:  "Text",
			eanHeader:   "EAN",
			timesHeader: "Times",
			wantLen:     2,
			wantErr:     false,
		},
		{
			name: "Times column not found when specified",
			table: Table{
				{"Text", "EAN"},
				{"Product", "1234567890123"},
			},
			textHeader:  "Text",
			eanHeader:   "EAN",
			timesHeader: "Count",
			wantLen:     0,
			wantErr:     true,
		},
		{
			name: "Header only table",
			table: Table{
				{"Text", "EAN"},
			},
			textHeader:  "Text",
			eanHeader:   "EAN",
			timesHeader: "",
			wantLen:     0,
			wantErr:     false,
		},
		{
			name: "Multiple columns",
			table: Table{
				{"ID", "Text", "EAN", "Price", "Stock"},
				{"1", "Product A", "1234567890123", "9.99", "100"},
				{"2", "Product B", "9876543210987", "19.99", "50"},
			},
			textHeader:  "Text",
			eanHeader:   "EAN",
			timesHeader: "",
			wantLen:     2,
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			records, err := RecordsFromTable(tt.table, tt.textHeader, tt.eanHeader, tt.timesHeader)
			if (err != nil) != tt.wantErr {
				t.Errorf("RecordsFromTable() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && len(records) != tt.wantLen {
				t.Errorf("RecordsFromTable() returned %d records, want %d", len(records), tt.wantLen)
			}
		})
	}
}

func TestRecordsFromTable_RecordContent(t *testing.T) {
	table := Table{
		{"Text", "EAN"},
		{"Product A", "1234567890123"},
	}

	records, err := RecordsFromTable(table, "Text", "EAN", "")
	if err != nil {
		t.Fatalf("RecordsFromTable() failed: %v", err)
	}

	if len(records) != 1 {
		t.Fatalf("Expected 1 record, got %d", len(records))
	}

	if records[0].Text != "Product A" {
		t.Errorf("Record Text = %v, want %v", records[0].Text, "Product A")
	}
	if records[0].Ean != "1234567890123" {
		t.Errorf("Record Ean = %v, want %v", records[0].Ean, "1234567890123")
	}
}

func TestRecordsFromTable_TimesColumn(t *testing.T) {
	tests := []struct {
		name       string
		table      Table
		wantTimes  int
		wantErr    bool
	}{
		{
			name: "Integer times",
			table: Table{
				{"Text", "EAN", "Times"},
				{"Product", "1234567890123", "5"},
			},
			wantTimes: 5,
			wantErr:   false,
		},
		{
			name: "Float times",
			table: Table{
				{"Text", "EAN", "Times"},
				{"Product", "1234567890123", "3.0"},
			},
			wantTimes: 3,
			wantErr:   false,
		},
		{
			name: "Float times truncated",
			table: Table{
				{"Text", "EAN", "Times"},
				{"Product", "1234567890123", "3.7"},
			},
			wantTimes: 3,
			wantErr:   false,
		},
		{
			name: "Times with spaces",
			table: Table{
				{"Text", "EAN", "Times"},
				{"Product", "1234567890123", "  4  "},
			},
			wantTimes: 4,
			wantErr:   false,
		},
		{
			name: "Invalid times string",
			table: Table{
				{"Text", "EAN", "Times"},
				{"Product", "1234567890123", "abc"},
			},
			wantTimes: 0,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			records, err := RecordsFromTable(tt.table, "Text", "EAN", "Times")
			if (err != nil) != tt.wantErr {
				t.Errorf("RecordsFromTable() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && len(records) > 0 && records[0].Times != tt.wantTimes {
				t.Errorf("Record.Times = %v, want %v", records[0].Times, tt.wantTimes)
			}
		})
	}
}

func TestRecordsFromTable_DefaultTimes(t *testing.T) {
	table := Table{
		{"Text", "EAN"},
		{"Product", "1234567890123"},
	}

	records, err := RecordsFromTable(table, "Text", "EAN", "")
	if err != nil {
		t.Fatalf("RecordsFromTable() failed: %v", err)
	}

	if len(records) != 1 {
		t.Fatalf("Expected 1 record, got %d", len(records))
	}

	// Default times should be 1
	if records[0].Times != 1 {
		t.Errorf("Record.Times = %v, want 1 (default)", records[0].Times)
	}
}

func TestRecordsFromTable_EmptyTimesHeader(t *testing.T) {
	// When times header is empty string (whitespace only), it should not search for times column
	table := Table{
		{"Text", "EAN"},
		{"Product", "1234567890123"},
	}

	records, err := RecordsFromTable(table, "Text", "EAN", "   ")
	if err != nil {
		t.Fatalf("RecordsFromTable() failed: %v", err)
	}

	if len(records) != 1 {
		t.Fatalf("Expected 1 record, got %d", len(records))
	}
}

func TestRecord_GenerateBarcode_ValidEAN13(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "barcode-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	record := Record{
		Text: "Test Product",
		Ean:  "5901234123457", // Valid EAN-13
	}

	barcodePath := filepath.Join(tmpDir, "barcode.png")
	err = record.GenerateBarcode(barcodePath)
	if err != nil {
		t.Errorf("GenerateBarcode() failed: %v", err)
	}

	// Check file was created
	if _, err := os.Stat(barcodePath); os.IsNotExist(err) {
		t.Error("Barcode file was not created")
	}

	// Check file is not empty
	info, err := os.Stat(barcodePath)
	if err != nil {
		t.Errorf("Failed to stat barcode file: %v", err)
	}
	if info.Size() == 0 {
		t.Error("Barcode file is empty")
	}
}

func TestRecord_GenerateBarcode_ValidEAN8(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "barcode-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	record := Record{
		Text: "Test Product",
		Ean:  "96385074", // Valid EAN-8
	}

	barcodePath := filepath.Join(tmpDir, "barcode.png")
	err = record.GenerateBarcode(barcodePath)
	if err != nil {
		t.Errorf("GenerateBarcode() failed for EAN-8: %v", err)
	}

	if _, err := os.Stat(barcodePath); os.IsNotExist(err) {
		t.Error("Barcode file was not created")
	}
}

func TestRecord_GenerateBarcode_InvalidEAN(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "barcode-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	tests := []struct {
		name string
		ean  string
	}{
		{name: "Too short", ean: "123"},
		{name: "Too long", ean: "12345678901234567890"},
		{name: "Invalid characters", ean: "123456789ABCD"},
		{name: "Empty", ean: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			record := Record{
				Text: "Test",
				Ean:  tt.ean,
			}

			barcodePath := filepath.Join(tmpDir, tt.name+".png")
			err := record.GenerateBarcode(barcodePath)
			if err == nil {
				t.Errorf("GenerateBarcode() should fail for invalid EAN: %s", tt.ean)
			}
		})
	}
}

func TestRecord_GenerateBarcode_InvalidPath(t *testing.T) {
	record := Record{
		Text: "Test Product",
		Ean:  "5901234123457",
	}

	err := record.GenerateBarcode("/nonexistent/directory/barcode.png")
	if err == nil {
		t.Error("GenerateBarcode() should fail for invalid path")
	}
}

func TestRecord_GenerateBarcode_MultipleRecords(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "barcode-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	records := []Record{
		{Text: "Product 1", Ean: "5901234123457"},
		{Text: "Product 2", Ean: "4006381333931"},
		{Text: "Product 3", Ean: "96385074"}, // EAN-8
	}

	for i, record := range records {
		barcodePath := filepath.Join(tmpDir, record.Ean+".png")
		err := record.GenerateBarcode(barcodePath)
		if err != nil {
			t.Errorf("GenerateBarcode() failed for record %d: %v", i, err)
		}

		if _, err := os.Stat(barcodePath); os.IsNotExist(err) {
			t.Errorf("Barcode file for record %d was not created", i)
		}
	}
}

func TestRecord_Fields(t *testing.T) {
	record := Record{
		Text:  "Test Product",
		Ean:   "1234567890123",
		Times: 5,
	}

	if record.Text != "Test Product" {
		t.Errorf("Record.Text = %v, want %v", record.Text, "Test Product")
	}
	if record.Ean != "1234567890123" {
		t.Errorf("Record.Ean = %v, want %v", record.Ean, "1234567890123")
	}
	if record.Times != 5 {
		t.Errorf("Record.Times = %v, want %v", record.Times, 5)
	}
}

func TestRecord_ZeroValue(t *testing.T) {
	var record Record

	if record.Text != "" {
		t.Errorf("Zero Record.Text = %v, want empty string", record.Text)
	}
	if record.Ean != "" {
		t.Errorf("Zero Record.Ean = %v, want empty string", record.Ean)
	}
	if record.Times != 0 {
		t.Errorf("Zero Record.Times = %v, want 0", record.Times)
	}
}
