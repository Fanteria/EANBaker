package core

import (
	"errors"
	"log/slog"
	"os"
	"path/filepath"

	"codeberg.org/go-pdf/fpdf"
)

type Pdf struct {
	pdf *fpdf.Fpdf
}

// Creates and configures a new PDF document for barcode generation.
// Sets up landscape orientation with custom dimensions (15x30mm) suitable for barcode labels.
// Disables auto page breaks and removes top margin for optimal barcode layout.
func NewPdf() Pdf {
	// Create pdf
	pdf := fpdf.NewCustom(&fpdf.InitType{
		OrientationStr: "L",
		UnitStr:        "mm",
		Size: fpdf.SizeType{
			Wd: 15,
			Ht: 30,
		},
	})
	pdf.SetTopMargin(0)
	pdf.SetAutoPageBreak(false, 0)
	return Pdf{pdf: pdf}
}

// Adds barcode pages to the PDF for each record.
// Creates temporary barcode images and adds the specified number of pages per record.
// Each page contains the record text, barcode image, and EAN number.
func (p *Pdf) AddPages(records []Record, times uint, log *slog.Logger) error {
	if times == 0 {
		const ERR_MSG string = "Bar code must be added at lease once time."
		log.Error(ERR_MSG)
		return errors.New(ERR_MSG)
	}
	// Create temporary directory
	dir, err := os.MkdirTemp("", "generate-barcodes-*")
	if err != nil {
		log.Error("Failed to create temporary dir", "err", err)
		return err
	}
	defer os.RemoveAll(dir)

	// Add records to pdf
	for _, record := range records {
		barcode_path := filepath.Join(dir, record.Ean+".png")
		err := record.GenerateBarcode(barcode_path)
		if err != nil {
			log.Error("Failed to generate barcode", "err", err)
			return err
		}
		for i := 0; i < int(times)+record.Times; i++ {
			log.Debug("Add page", "record", record, "barcode", barcode_path)
			p.addPage(record, barcode_path)
		}
	}
	log.Info("Pages added", "count", p.pdf.PageCount())

	return nil
}

// Adds a single barcode page to the PDF document.
// Layouts the record text at the top, barcode image
// in the center, and EAN number at the bottom.
func (p *Pdf) addPage(record Record, image string) error {
	p.pdf.AddPage()
	p.pdf.SetFont("Arial", "", 4)

	// Top text
	p.pdf.SetXY(1, 1)
	p.pdf.MultiCell(27.0, 1.6, record.Text, "", "L", false)

	// Center image
	imageWidth := 27.0
	imageHeight := 7.0
	p.pdf.ImageOptions(image, 1.5, 5.0, imageWidth, imageHeight, false, fpdf.ImageOptions{}, 0, "")

	// Footer EAN in text
	p.pdf.SetFooterFuncLpi(func(lastPage bool) {
		p.pdf.SetXY(0, 0)
		p.pdf.CellFormat(30.0, 14.0, record.Ean, "", 0, "CB", false, 0, "")
	})
	return nil
}

// Save writes the PDF document to the specified file path and closes it.
// Returns an error if the file cannot be created or written.
func (p *Pdf) Save(path string) error {
	err := p.pdf.OutputFileAndClose(path)
	if err != nil {
		return err
	}
	return nil
}
