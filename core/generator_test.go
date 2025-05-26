package core

import (
	"testing"
)

func TestGenerator_Validate(t *testing.T) {
	tests := []struct {
		name    string // description of this test case
		gen     Generator
		wantErr bool
	}{
		{name: "Valid suffixes", gen: Generator{CsvPath: "a.csv", PdfPath: "a.pdf"}, wantErr: false},
		{name: "Invalid csv", gen: Generator{CsvPath: "a.txt", PdfPath: "a.pdf"}, wantErr: true},
		{name: "Invalid pdf", gen: Generator{CsvPath: "a.csv", PdfPath: "a.txt"}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := tt.gen.Validate()
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Validate() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Validate() succeeded unexpectedly")
			}
		})
	}
}

func TestGenerator_UpdatePdfPath(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		gen  Generator
		want string
	}{
		{name: "Simple csv name", gen: Generator{CsvPath: "data.csv", PdfPath: ""}, want: "data.pdf"},
		{
			name: "Pdf name already set",
			gen:  Generator{CsvPath: "data.csv", PdfPath: "different.pdf"},
			want: "different.pdf",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.gen.UpdatePdfPath()
			if tt.want != tt.gen.PdfPath {
				t.Errorf("GeneratePdfPath() = %v, want %v", tt.gen.PdfPath, tt.want)
			}
		})
	}
}

func TestGeneratePdfPath(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		path string
		want string
	}{
		{name: "Simple csv name", path: "data.csv", want: "data.pdf"},
		{name: "More dots", path: "data.some.text.csv", want: "data.some.text.pdf"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GeneratePdfPath(tt.path)
			if tt.want != got {
				t.Errorf("GeneratePdfPath() = %v, want %v", got, tt.want)
			}
		})
	}
}
