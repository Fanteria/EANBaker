package core

import (
	"bytes"
	"strings"
	"testing"
)

func TestTableFromCsv(t *testing.T) {
	tests := []struct {
		name    string
		csv     string
		comma   rune
		wantLen int
		wantErr bool
	}{
		{
			name:    "Simple CSV",
			csv:     "a,b,c\n1,2,3\n4,5,6",
			comma:   ',',
			wantLen: 3,
			wantErr: false,
		},
		{
			name:    "Semicolon separator",
			csv:     "a;b;c\n1;2;3",
			comma:   ';',
			wantLen: 2,
			wantErr: false,
		},
		{
			name:    "Tab separator",
			csv:     "a\tb\tc\n1\t2\t3",
			comma:   '\t',
			wantLen: 2,
			wantErr: false,
		},
		{
			name:    "Default comma (0)",
			csv:     "a,b,c\n1,2,3",
			comma:   0,
			wantLen: 2,
			wantErr: false,
		},
		{
			name:    "Empty CSV",
			csv:     "",
			comma:   ',',
			wantLen: 0,
			wantErr: false,
		},
		{
			name:    "Header only",
			csv:     "a,b,c",
			comma:   ',',
			wantLen: 1,
			wantErr: false,
		},
		{
			name:    "Quoted fields",
			csv:     `"a","b","c"` + "\n" + `"1","2","3"`,
			comma:   ',',
			wantLen: 2,
			wantErr: false,
		},
		{
			name:    "Quoted fields with commas",
			csv:     `"a,1","b,2","c,3"` + "\n" + `"1","2","3"`,
			comma:   ',',
			wantLen: 2,
			wantErr: false,
		},
		{
			name:    "Pipe separator",
			csv:     "a|b|c\n1|2|3",
			comma:   '|',
			wantLen: 2,
			wantErr: false,
		},
		{
			name:    "Single column",
			csv:     "a\n1\n2\n3",
			comma:   ',',
			wantLen: 4,
			wantErr: false,
		},
		{
			name:    "Windows line endings",
			csv:     "a,b,c\r\n1,2,3\r\n4,5,6",
			comma:   ',',
			wantLen: 3,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.csv)
			table, err := TableFromCsv(reader, tt.comma)
			if (err != nil) != tt.wantErr {
				t.Errorf("TableFromCsv() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && len(table) != tt.wantLen {
				t.Errorf("TableFromCsv() returned %d rows, want %d", len(table), tt.wantLen)
			}
		})
	}
}

func TestTableFromCsv_NilReader(t *testing.T) {
	_, err := TableFromCsv(nil, ',')
	if err == nil {
		t.Error("TableFromCsv() should fail for nil reader")
	}
}

func TestTableFromCsv_Content(t *testing.T) {
	csv := "Name,EAN,Price\nProduct A,1234567890123,9.99\nProduct B,9876543210987,19.99"
	reader := strings.NewReader(csv)

	table, err := TableFromCsv(reader, ',')
	if err != nil {
		t.Fatalf("TableFromCsv() failed: %v", err)
	}

	// Check header
	if len(table) < 1 {
		t.Fatal("Table has no rows")
	}
	expectedHeader := []string{"Name", "EAN", "Price"}
	for i, col := range expectedHeader {
		if table[0][i] != col {
			t.Errorf("Header[%d] = %v, want %v", i, table[0][i], col)
		}
	}

	// Check first data row
	if len(table) < 2 {
		t.Fatal("Table has no data rows")
	}
	if table[1][0] != "Product A" {
		t.Errorf("Row 1, Col 0 = %v, want %v", table[1][0], "Product A")
	}
	if table[1][1] != "1234567890123" {
		t.Errorf("Row 1, Col 1 = %v, want %v", table[1][1], "1234567890123")
	}
	if table[1][2] != "9.99" {
		t.Errorf("Row 1, Col 2 = %v, want %v", table[1][2], "9.99")
	}
}

func TestTableFromCsv_MalformedCSV(t *testing.T) {
	// Inconsistent number of fields
	csv := "a,b,c\n1,2"
	reader := strings.NewReader(csv)

	_, err := TableFromCsv(reader, ',')
	if err == nil {
		t.Error("TableFromCsv() should fail for malformed CSV")
	}
}

func TestTableFromCsv_SpecialCharacters(t *testing.T) {
	tests := []struct {
		name  string
		csv   string
		want  string
	}{
		{
			name: "UTF-8 characters",
			csv:  "Name\nPříliš žluťoučký kůň",
			want: "Příliš žluťoučký kůň",
		},
		{
			name: "Emoji",
			csv:  "Name\n🎉 Party",
			want: "🎉 Party",
		},
		{
			name: "Newline in quoted field",
			csv:  "Name\n\"Line 1\nLine 2\"",
			want: "Line 1\nLine 2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.csv)
			table, err := TableFromCsv(reader, ',')
			if err != nil {
				t.Fatalf("TableFromCsv() failed: %v", err)
			}
			if len(table) < 2 {
				t.Fatal("Table has insufficient rows")
			}
			if table[1][0] != tt.want {
				t.Errorf("Cell content = %v, want %v", table[1][0], tt.want)
			}
		})
	}
}

func TestTableFromCsv_EmptyFields(t *testing.T) {
	csv := "a,b,c\n,,\n1,,3"
	reader := strings.NewReader(csv)

	table, err := TableFromCsv(reader, ',')
	if err != nil {
		t.Fatalf("TableFromCsv() failed: %v", err)
	}

	if len(table) != 3 {
		t.Fatalf("Expected 3 rows, got %d", len(table))
	}

	// Check empty row
	for i := 0; i < 3; i++ {
		if table[1][i] != "" {
			t.Errorf("Row 1, Col %d = %v, want empty string", i, table[1][i])
		}
	}

	// Check mixed row
	if table[2][0] != "1" {
		t.Errorf("Row 2, Col 0 = %v, want %v", table[2][0], "1")
	}
	if table[2][1] != "" {
		t.Errorf("Row 2, Col 1 = %v, want empty string", table[2][1])
	}
	if table[2][2] != "3" {
		t.Errorf("Row 2, Col 2 = %v, want %v", table[2][2], "3")
	}
}

func TestTableFromCsv_LargeFile(t *testing.T) {
	// Generate a large CSV
	var sb strings.Builder
	sb.WriteString("ID,Name,Value\n")
	for i := 0; i < 1000; i++ {
		sb.WriteString("1,Test,100\n")
	}

	reader := strings.NewReader(sb.String())
	table, err := TableFromCsv(reader, ',')
	if err != nil {
		t.Fatalf("TableFromCsv() failed: %v", err)
	}

	// 1 header + 1000 data rows
	if len(table) != 1001 {
		t.Errorf("Expected 1001 rows, got %d", len(table))
	}
}

func TestTableFromExcel_InvalidData(t *testing.T) {
	// Not a valid Excel file
	reader := strings.NewReader("not an excel file")
	_, err := TableFromExcel(reader, 0)
	if err == nil {
		t.Error("TableFromExcel() should fail for invalid data")
	}
}

func TestTableFromExcel_EmptyReader(t *testing.T) {
	reader := bytes.NewReader([]byte{})
	_, err := TableFromExcel(reader, 0)
	if err == nil {
		t.Error("TableFromExcel() should fail for empty reader")
	}
}

func TestTable_Type(t *testing.T) {
	table := Table{
		{"a", "b", "c"},
		{"1", "2", "3"},
	}

	// Table should be a [][]string
	var data [][]string = table
	if len(data) != 2 {
		t.Errorf("Table length = %d, want 2", len(data))
	}
}

func TestTable_EmptyTable(t *testing.T) {
	table := Table{}
	if len(table) != 0 {
		t.Errorf("Empty table length = %d, want 0", len(table))
	}
}

func TestTable_NilTable(t *testing.T) {
	var table Table
	if table != nil {
		t.Error("Nil table should be nil")
	}
}

func TestTable_Modification(t *testing.T) {
	table := Table{
		{"a", "b"},
		{"1", "2"},
	}

	// Modify table
	table[0][0] = "modified"
	if table[0][0] != "modified" {
		t.Error("Table modification failed")
	}

	// Append row
	table = append(table, []string{"3", "4"})
	if len(table) != 3 {
		t.Errorf("Table append failed, length = %d, want 3", len(table))
	}
}

func TestTableFromCsv_RealWorldData(t *testing.T) {
	// Simulate real-world CSV data with EAN codes
	csv := `Material Number,EAN,Stock,Price
Product Alpha,5901234123457,100,9.99
Product Beta,4006381333931,50,19.99
Product Gamma,96385074,200,4.99`

	reader := strings.NewReader(csv)
	table, err := TableFromCsv(reader, ',')
	if err != nil {
		t.Fatalf("TableFromCsv() failed: %v", err)
	}

	// Verify structure
	if len(table) != 4 {
		t.Errorf("Expected 4 rows, got %d", len(table))
	}

	// Verify headers
	expectedHeaders := []string{"Material Number", "EAN", "Stock", "Price"}
	for i, h := range expectedHeaders {
		if table[0][i] != h {
			t.Errorf("Header[%d] = %v, want %v", i, table[0][i], h)
		}
	}

	// Verify EAN values
	expectedEANs := []string{"5901234123457", "4006381333931", "96385074"}
	for i, ean := range expectedEANs {
		if table[i+1][1] != ean {
			t.Errorf("EAN[%d] = %v, want %v", i, table[i+1][1], ean)
		}
	}
}

func TestTableFromCsv_SemicolonSeparatedRealWorld(t *testing.T) {
	// European-style CSV with semicolon separator
	csv := `Name;EAN;Price
"Produkt A";5901234123457;9,99
"Produkt B";4006381333931;19,99`

	reader := strings.NewReader(csv)
	table, err := TableFromCsv(reader, ';')
	if err != nil {
		t.Fatalf("TableFromCsv() failed: %v", err)
	}

	if len(table) != 3 {
		t.Errorf("Expected 3 rows, got %d", len(table))
	}

	// Check quoted field
	if table[1][0] != "Produkt A" {
		t.Errorf("Name = %v, want %v", table[1][0], "Produkt A")
	}

	// Check price with European decimal
	if table[1][2] != "9,99" {
		t.Errorf("Price = %v, want %v", table[1][2], "9,99")
	}
}
