package core

import (
	"log"
	"os"
	"path/filepath"

	"codeberg.org/go-pdf/fpdf"
)

type Pdf struct {
	pdf *fpdf.Fpdf
}

// Create new pdf instance
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

// Add records to pdf
func (p *Pdf) AddPages(records []Record) error {
	// Create temporary directory
	dir, err := os.MkdirTemp("", "generate-barcodes-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)

	// Add records to pdf
	for _, record := range records {
		barcode_path := filepath.Join(dir, record.Ean+".png")
		err := record.GenerateBarcode(barcode_path)
		if err != nil {
			log.Fatal(err)
		}
		p.addPage(record, barcode_path)
	}

	return nil
}

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

	// Footer ean in text
	p.pdf.SetFooterFuncLpi(func(lastPage bool) {
		p.pdf.SetXY(0, 0)
		p.pdf.CellFormat(30.0, 14.0, record.Ean, "", 0, "CB", false, 0, "")
	})
	return nil
}

// Save instance to new pdf file.
func (p *Pdf) Save(path string) error {
	err := p.pdf.OutputFileAndClose(path)
	if err != nil {
		return err
	}
	return nil
}
