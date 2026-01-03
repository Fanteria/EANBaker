package core

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"testing"
)

func TestNewPdf(t *testing.T) {
	pdf := NewPdf()

	if pdf.pdf == nil {
		t.Error("NewPdf() returned Pdf with nil internal pdf")
	}
}

func TestNewPdf_InitialState(t *testing.T) {
	pdf := NewPdf()

	// New PDF should have no pages
	if pdf.pdf.PageCount() != 0 {
		t.Errorf("New PDF should have 0 pages, got %d", pdf.pdf.PageCount())
	}
}

func TestPdf_AddPages_SingleRecord(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "pdf-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	pdf := NewPdf()
	records := []Record{
		{Text: "Product A", Ean: "5901234123457", Times: 0},
	}

	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	err = pdf.AddPages(records, 1, log)
	if err != nil {
		t.Errorf("AddPages() failed: %v", err)
	}

	// Should have 1 page
	if pdf.pdf.PageCount() != 1 {
		t.Errorf("PageCount() = %d, want 1", pdf.pdf.PageCount())
	}
}

func TestPdf_AddPages_MultipleRecords(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "pdf-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	pdf := NewPdf()
	records := []Record{
		{Text: "Product A", Ean: "5901234123457", Times: 0},
		{Text: "Product B", Ean: "4006381333931", Times: 0},
		{Text: "Product C", Ean: "96385074", Times: 0}, // EAN-8
	}

	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	err = pdf.AddPages(records, 1, log)
	if err != nil {
		t.Errorf("AddPages() failed: %v", err)
	}

	if pdf.pdf.PageCount() != 3 {
		t.Errorf("PageCount() = %d, want 3", pdf.pdf.PageCount())
	}
}

func TestPdf_AddPages_TimesMultiplier(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "pdf-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	pdf := NewPdf()
	records := []Record{
		{Text: "Product", Ean: "5901234123457", Times: 0},
	}

	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	err = pdf.AddPages(records, 5, log) // 5 copies
	if err != nil {
		t.Errorf("AddPages() failed: %v", err)
	}

	// Should have 5 pages (1 record * 5 times)
	if pdf.pdf.PageCount() != 5 {
		t.Errorf("PageCount() = %d, want 5", pdf.pdf.PageCount())
	}
}

func TestPdf_AddPages_RecordTimes(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "pdf-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	pdf := NewPdf()
	records := []Record{
		{Text: "Product", Ean: "5901234123457", Times: 3}, // Record has Times=3
	}

	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	err = pdf.AddPages(records, 1, log) // Base times = 1
	if err != nil {
		t.Errorf("AddPages() failed: %v", err)
	}

	// Should have 4 pages (base 1 + record times 3)
	if pdf.pdf.PageCount() != 4 {
		t.Errorf("PageCount() = %d, want 4", pdf.pdf.PageCount())
	}
}

func TestPdf_AddPages_CombinedTimes(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "pdf-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	pdf := NewPdf()
	records := []Record{
		{Text: "Product A", Ean: "5901234123457", Times: 2},
		{Text: "Product B", Ean: "4006381333931", Times: 0},
	}

	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	err = pdf.AddPages(records, 2, log) // Base times = 2
	if err != nil {
		t.Errorf("AddPages() failed: %v", err)
	}

	// Product A: 2 (base) + 2 (record) = 4 pages
	// Product B: 2 (base) + 0 (record) = 2 pages
	// Total: 6 pages
	if pdf.pdf.PageCount() != 6 {
		t.Errorf("PageCount() = %d, want 6", pdf.pdf.PageCount())
	}
}

func TestPdf_AddPages_ZeroTimes(t *testing.T) {
	pdf := NewPdf()
	records := []Record{
		{Text: "Product", Ean: "5901234123457", Times: 0},
	}

	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	err := pdf.AddPages(records, 0, log) // times = 0 should error
	if err == nil {
		t.Error("AddPages() should fail when times = 0")
	}
}

func TestPdf_AddPages_EmptyRecords(t *testing.T) {
	pdf := NewPdf()
	records := []Record{}

	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	err := pdf.AddPages(records, 1, log)
	if err != nil {
		t.Errorf("AddPages() should succeed with empty records: %v", err)
	}

	if pdf.pdf.PageCount() != 0 {
		t.Errorf("PageCount() = %d, want 0 for empty records", pdf.pdf.PageCount())
	}
}

func TestPdf_AddPages_InvalidEAN(t *testing.T) {
	pdf := NewPdf()
	records := []Record{
		{Text: "Product", Ean: "invalid-ean", Times: 0},
	}

	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	err := pdf.AddPages(records, 1, log)
	if err == nil {
		t.Error("AddPages() should fail for invalid EAN")
	}
}

func TestPdf_Save(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "pdf-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	pdf := NewPdf()
	records := []Record{
		{Text: "Product", Ean: "5901234123457", Times: 0},
	}

	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	err = pdf.AddPages(records, 1, log)
	if err != nil {
		t.Fatalf("AddPages() failed: %v", err)
	}

	pdfPath := filepath.Join(tmpDir, "output.pdf")
	err = pdf.Save(pdfPath)
	if err != nil {
		t.Errorf("Save() failed: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(pdfPath); os.IsNotExist(err) {
		t.Error("PDF file was not created")
	}

	// Verify file is not empty
	info, err := os.Stat(pdfPath)
	if err != nil {
		t.Errorf("Failed to stat PDF file: %v", err)
	}
	if info.Size() == 0 {
		t.Error("PDF file is empty")
	}
}

func TestPdf_Save_InvalidPath(t *testing.T) {
	pdf := NewPdf()
	records := []Record{
		{Text: "Product", Ean: "5901234123457", Times: 0},
	}

	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	_ = pdf.AddPages(records, 1, log)

	err := pdf.Save("/nonexistent/directory/output.pdf")
	if err == nil {
		t.Error("Save() should fail for invalid path")
	}
}

func TestPdf_Save_EmptyPdf(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "pdf-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	pdf := NewPdf()
	// No pages added

	pdfPath := filepath.Join(tmpDir, "empty.pdf")
	err = pdf.Save(pdfPath)
	// Empty PDF might still save (library dependent)
	// Just verify no crash
	_ = err
}

func TestPdf_Save_OverwriteExisting(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "pdf-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	pdfPath := filepath.Join(tmpDir, "output.pdf")

	// Create first PDF
	pdf1 := NewPdf()
	records1 := []Record{
		{Text: "Product 1", Ean: "5901234123457", Times: 0},
	}
	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	_ = pdf1.AddPages(records1, 1, log)
	err = pdf1.Save(pdfPath)
	if err != nil {
		t.Fatalf("First Save() failed: %v", err)
	}

	firstSize, _ := os.Stat(pdfPath)

	// Create second PDF and overwrite
	pdf2 := NewPdf()
	records2 := []Record{
		{Text: "Product 1", Ean: "5901234123457", Times: 0},
		{Text: "Product 2", Ean: "4006381333931", Times: 0},
	}
	_ = pdf2.AddPages(records2, 1, log)
	err = pdf2.Save(pdfPath)
	if err != nil {
		t.Fatalf("Second Save() failed: %v", err)
	}

	secondSize, _ := os.Stat(pdfPath)

	// Second file should be larger (more content)
	if secondSize.Size() <= firstSize.Size() {
		t.Error("Overwritten PDF should be larger")
	}
}

func TestPdf_FullWorkflow(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "pdf-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create PDF
	pdf := NewPdf()

	// Add various records
	records := []Record{
		{Text: "Small Product", Ean: "96385074", Times: 0},       // EAN-8
		{Text: "Regular Product", Ean: "5901234123457", Times: 1}, // EAN-13 with extra
		{Text: "Multi-line\nProduct Name", Ean: "4006381333931", Times: 0}, // Multi-line text
	}

	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	err = pdf.AddPages(records, 2, log)
	if err != nil {
		t.Fatalf("AddPages() failed: %v", err)
	}

	// Save
	pdfPath := filepath.Join(tmpDir, "workflow.pdf")
	err = pdf.Save(pdfPath)
	if err != nil {
		t.Fatalf("Save() failed: %v", err)
	}

	// Verify
	info, err := os.Stat(pdfPath)
	if err != nil {
		t.Fatalf("Failed to stat PDF: %v", err)
	}
	if info.Size() == 0 {
		t.Error("PDF file is empty")
	}

	// Expected pages:
	// Record 1: 2 (base) + 0 = 2
	// Record 2: 2 (base) + 1 = 3
	// Record 3: 2 (base) + 0 = 2
	// Total: 7
	expectedPages := 7
	if pdf.pdf.PageCount() != expectedPages {
		t.Errorf("PageCount() = %d, want %d", pdf.pdf.PageCount(), expectedPages)
	}
}

func TestPdf_LargeNumberOfPages(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "pdf-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	pdf := NewPdf()
	
	// Create many records
	records := make([]Record, 50)
	for i := 0; i < 50; i++ {
		records[i] = Record{
			Text:  "Product",
			Ean:   "5901234123457",
			Times: 0,
		}
	}

	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	err = pdf.AddPages(records, 1, log)
	if err != nil {
		t.Fatalf("AddPages() failed: %v", err)
	}

	pdfPath := filepath.Join(tmpDir, "large.pdf")
	err = pdf.Save(pdfPath)
	if err != nil {
		t.Fatalf("Save() failed: %v", err)
	}

	if pdf.pdf.PageCount() != 50 {
		t.Errorf("PageCount() = %d, want 50", pdf.pdf.PageCount())
	}
}

func TestPdf_SpecialCharactersInText(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "pdf-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	pdf := NewPdf()
	records := []Record{
		{Text: "Příliš žluťoučký kůň", Ean: "5901234123457", Times: 0}, // Czech characters
		{Text: "Product™ ®", Ean: "4006381333931", Times: 0},           // Symbols
		{Text: "Café & Co.", Ean: "96385074", Times: 0},                // Ampersand
	}

	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	err = pdf.AddPages(records, 1, log)
	if err != nil {
		t.Errorf("AddPages() failed with special characters: %v", err)
	}

	pdfPath := filepath.Join(tmpDir, "special.pdf")
	err = pdf.Save(pdfPath)
	if err != nil {
		t.Errorf("Save() failed: %v", err)
	}
}
