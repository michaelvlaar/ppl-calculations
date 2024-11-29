package queries

import (
	"bytes"
	"context"
	"fmt"
	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/sirupsen/logrus"
	"io"
	"ppl-calculations/domain/state"
)

type PdfExportHandler struct {
}

func NewPdfExportHandler() PdfExportHandler {
	return PdfExportHandler{}
}

func (h PdfExportHandler) Handle(_ context.Context, stateService state.Service) (io.Reader, error) {
	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		logrus.WithError(err).Error("generating pdf generator")
		return nil, err
	}

	// Add a new page with HTML content
	page := wkhtmltopdf.NewPageReader(bytes.NewReader([]byte(`
		
	`)))

	pdfg.PageSize.Set(wkhtmltopdf.PageSizeA4)
	pdfg.Orientation.Set(wkhtmltopdf.OrientationPortrait)
	pdfg.AddPage(page)

	pdfg.Orientation.Set(wkhtmltopdf.OrientationLandscape)
	pdfg.AddPage(page)

	var pdfBuffer bytes.Buffer
	pdfg.SetOutput(&pdfBuffer)

	err = pdfg.Create()
	if err != nil {
		fmt.Println("Error generating PDF:", err)
		return nil, nil
	}

	return &pdfBuffer, nil
}
