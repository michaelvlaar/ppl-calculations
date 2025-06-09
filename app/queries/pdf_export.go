package queries

import (
	"bytes"
	"context"
	"fmt"
	"github.com/michaelvlaar/ppl-calculations/domain/calculations"
	"github.com/michaelvlaar/ppl-calculations/domain/export"
	"github.com/michaelvlaar/ppl-calculations/domain/fuel"
	"github.com/phpdave11/gofpdf"
	"io"
	"io/fs"
	"strings"
	"time"
)

const (
	TitleFontSize       = 14
	TableHeaderFontSize = 12
)

type PdfExportHandler struct {
	calcService calculations.Service
	fontRegular []byte
	fontBold    []byte
	fontItalic  []byte
	version     string
}

type WeightAndBalanceItem struct {
	Name       string
	LeverArm   string
	Mass       string
	MassMoment string
}

type WeightAndBalanceState struct {
	Items        []WeightAndBalanceItem
	Total        WeightAndBalanceItem
	WithinLimits bool
}

type ExportData struct {
	CallSign  string
	Generated string
	Reference string

	FuelTaxi        string
	FuelTrip        string
	FuelAlternate   string
	FuelContingency string
	FuelReserve     string
	FuelTotal       string
	FuelExtra       string
	FuelExtraAbs    string
	FuelSufficient  bool

	ChartUrl string
	LdrUrl   string
	TdrUrl   string

	TakeOffRunRequired        string
	TakeOffDistanceRequired   string
	LandingDistanceRequired   string
	LandingGroundRollRequired string

	WeightAndBalanceTakeOff WeightAndBalanceState
	WeightAndBalanceLanding WeightAndBalanceState
}

func NewPdfExportHandler(version string, assets fs.FS, calcService calculations.Service) PdfExportHandler {
	robotoRegular, err := fs.ReadFile(assets, "assets/fonts/Roboto-Regular.ttf")
	if err != nil {
		panic(err)
	}

	robotoBold, err := fs.ReadFile(assets, "assets/fonts/Roboto-Bold.ttf")
	if err != nil {
		panic(err)
	}

	robotoItalic, err := fs.ReadFile(assets, "assets/fonts/Roboto-Italic.ttf")
	if err != nil {
		panic(err)
	}

	return PdfExportHandler{
		fontRegular: robotoRegular,
		fontBold:    robotoBold,
		fontItalic:  robotoItalic,
		calcService: calcService,
		version:     version,
	}
}

func parseNumber(number string) string {
	return strings.ReplaceAll(number, ".", ",")
}

func (h PdfExportHandler) Handle(_ context.Context, e export.Export) (io.Reader, error) {
	takeOffWeightAndBalance, err := calculations.NewWeightAndBalance(e.CallSign, e.Pilot, e.PilotSeat, e.Passenger, e.PassengerSeat, e.Baggage, e.Fuel)
	if err != nil {
		return nil, err
	}

	fuelPlanning, err := calculations.NewFuelPlanning(e.TripDuration, e.AlternateDuration, e.Fuel, e.Fuel.Volume.Type)
	if err != nil {
		return nil, err
	}

	landingWeightAndBalance, err := calculations.NewWeightAndBalance(e.CallSign, e.Pilot, e.PilotSeat, e.Passenger, e.PassengerSeat, e.Baggage, fuel.Subtract(e.Fuel, fuelPlanning.Trip, fuelPlanning.Taxi))
	if err != nil {
		return nil, err
	}

	tdrChart, todRR, todDR, err := h.calcService.TakeOffDistance(e.OutsideAirTemperature, e.PressureAltitude, takeOffWeightAndBalance.Total.Mass, e.Wind, calculations.ChartTypePNG)
	if err != nil {
		return nil, err
	}

	ldrChart, ldrDR, ldrGR, err := h.calcService.LandingDistance(e.OutsideAirTemperature, e.PressureAltitude, landingWeightAndBalance.Total.Mass, e.Wind, calculations.ChartTypePNG)
	if err != nil {
		return nil, err
	}

	wbChart, err := h.calcService.WeightAndBalance(e.CallSign, *takeOffWeightAndBalance.Total, *landingWeightAndBalance.Total, takeOffWeightAndBalance.WithinLimits, calculations.ChartTypePNG)
	if err != nil {
		return nil, err
	}

	performance := &calculations.Performance{
		TakeOffRunRequired:        todRR,
		TakeOffDistanceRequired:   todDR,
		LandingDistanceRequired:   ldrDR,
		LandingGroundRollRequired: ldrGR,
	}

	pdf := gofpdf.New("P", "mm", "A4", "")

	wbImageOptions := gofpdf.ImageOptions{ImageType: "PNG", ReadDpi: true}
	tdrImageOptions := gofpdf.ImageOptions{ImageType: "PNG", ReadDpi: true}
	ldrImageOptions := gofpdf.ImageOptions{ImageType: "PNG", ReadDpi: true}

	pdf.RegisterImageOptionsReader("wb", wbImageOptions, wbChart)
	pdf.RegisterImageOptionsReader("tdr", tdrImageOptions, tdrChart)
	pdf.RegisterImageOptionsReader("ldr", ldrImageOptions, ldrChart)

	pdf.AddUTF8FontFromBytes("Roboto", "", h.fontRegular)
	pdf.AddUTF8FontFromBytes("Roboto", "B", h.fontBold)
	pdf.AddUTF8FontFromBytes("Roboto", "I", h.fontItalic)

	pdf.SetHeaderFuncMode(func() {
		pdf.SetFillColor(78, 65, 246)
		pdf.Rect(0, 0, 210, 15, "F")
		pdf.SetY(5)
		pdf.SetFont("Roboto", "B", 10)
		pdf.SetTextColor(255, 255, 255)
		pdf.SetX(25)
		pdf.CellFormat(160/3, 5, e.CallSign.String(), "", 0, "L", false, 0, "")
		pdf.CellFormat(160/3, 5, e.Name.String(), "", 0, "C", false, 0, "")
		pdf.CellFormat(160/3, 5, time.Now().Format("15:04 02-01-2006"), "", 0, "R", false, 0, "")
	}, true)

	pdf.SetTextColor(0, 0, 0)
	pdf.AliasNbPages("{nb}")
	pdf.SetFooterFunc(func() {
		pdf.SetY(-12)
		pdf.SetFillColor(240, 240, 240)
		pdf.Rect(0, 282, 210, 15, "F")
		pdf.SetFont("Roboto", "I", 10)
		pdf.SetX(25)
		pdf.SetFillColor(240, 240, 240)
		pdf.CellFormat(100, 10, fmt.Sprintf("Gegenereerd door github.com/michaelvlaar/ppl-calculations (%s)", h.version), "", 0, "L", false, 0, "https://github.com/michaelvlaar/ppl-calculations")
		pdf.CellFormat(60, 10, "Pagina "+fmt.Sprint(pdf.PageNo())+" van {nb}", "", 0, "R", false, 0, "")
	})

	pdf.SetMargins(20, 20, 20)
	pdf.AddPage()
	pdf.SetFont("Roboto", "B", TitleFontSize)
	pdf.CellFormat(0, 10, "Gewicht en Balans", "", 1, "C", false, 0, "")

	if !takeOffWeightAndBalance.WithinLimits {
		pdf.SetFillColor(200, 0, 0)
		pdf.SetTextColor(255, 255, 255)
		pdf.SetFont("Roboto", "B", 10)
		pdf.SetX(25)
		pdf.MultiCell(160, 6, "De huidige gewichts- en balansberekening geeft aan dat de belading van het vliegtuig buiten de toegestane limieten valt. Controleer en herbereken de gewichts- en balansverdeling zorgvuldig om te voldoen aan de veiligheidsvoorschriften.", "1", "C", true)
		pdf.SetTextColor(0, 0, 0)
	}

	pdf.Ln(5)

	pdf.ImageOptions("wb", 55, pdf.GetY(), 100, 0, false, wbImageOptions, 0, "")
	pdf.Ln(105)

	pdf.SetFont("Roboto", "B", TitleFontSize)
	pdf.CellFormat(0, 10, "Take-off", "", 1, "C", false, 0, "")
	pdf.Ln(1)

	headers := []string{"NAME", "ARM [M]", "MASS [KG]", "MASS MOMENT [KG M]"}
	colWidths := []float64{50, 30, 30, 50}

	pdf.SetFont("Roboto", "B", TableHeaderFontSize)
	pdf.SetFillColor(170, 170, 170)
	pdf.SetX((210 - (colWidths[0] + colWidths[1] + colWidths[2] + colWidths[3])) / 2)
	for i, h := range headers {
		align := "C"
		if i == 0 {
			align = "L"
		}
		pdf.CellFormat(colWidths[i], 6, h, "1", 0, align, true, 0, "")
	}
	pdf.Ln(-1)
	pdf.SetFont("Roboto", "", TableHeaderFontSize)
	for _, i := range takeOffWeightAndBalance.Moments {
		pdf.SetX((210 - (colWidths[0] + colWidths[1] + colWidths[2] + colWidths[3])) / 2)
		mass := strings.ReplaceAll(fmt.Sprintf("%.2f", i.Mass.Kilo()), ".", ",")
		if i.Name == "Fuel" {
			mass = fmt.Sprintf("(%s) %s", parseNumber(fuelPlanning.Total.Volume.Subtract(fuelPlanning.Trip.Volume).String(fuelPlanning.VolumeType)), mass)
		}
		pdf.CellFormat(colWidths[0], 6, strings.ReplaceAll(i.Name, ".", ","), "1", 0, "", false, 0, "")
		pdf.CellFormat(colWidths[1], 6, strings.ReplaceAll(fmt.Sprintf("%.4f", i.Arm), ".", ","), "1", 0, "C", false, 0, "")
		pdf.CellFormat(colWidths[2], 6, strings.ReplaceAll(mass, ".", ","), "1", 0, "C", false, 0, "")
		pdf.CellFormat(colWidths[3], 6, strings.ReplaceAll(fmt.Sprintf("%.2f", i.KGM()), ".", ","), "1", 0, "C", false, 0, "")
		pdf.Ln(-1)
	}

	total := takeOffWeightAndBalance.Total
	pdf.SetFont("Roboto", "B", TableHeaderFontSize)
	pdf.SetFillColor(170, 170, 170)
	pdf.SetX((210 - (colWidths[0] + colWidths[1] + colWidths[2] + colWidths[3])) / 2)
	pdf.CellFormat(colWidths[0], 6, "TOTAL", "1", 0, "", true, 0, "")
	pdf.CellFormat(colWidths[1], 6, strings.ReplaceAll(fmt.Sprintf("%.4f", total.Arm), ".", ","), "1", 0, "C", true, 0, "")
	pdf.CellFormat(colWidths[2], 6, strings.ReplaceAll(fmt.Sprintf("%.2f", total.Mass.Kilo()), ".", ","), "1", 0, "C", true, 0, "")
	pdf.CellFormat(colWidths[3], 6, strings.ReplaceAll(fmt.Sprintf("%.2f", total.KGM()), ".", ","), "1", 0, "C", true, 0, "")
	pdf.Ln(10)

	pdf.SetFont("Roboto", "B", TitleFontSize)
	pdf.CellFormat(0, 10, "Landing", "", 1, "C", false, 0, "")
	pdf.Ln(1)

	pdf.SetFont("Roboto", "B", TableHeaderFontSize)
	pdf.SetFillColor(170, 170, 170)
	pdf.SetX((210 - (colWidths[0] + colWidths[1] + colWidths[2] + colWidths[3])) / 2)
	for i, h := range headers {
		align := "C"
		if i == 0 {
			align = "L"
		}
		pdf.CellFormat(colWidths[i], 6, h, "1", 0, align, true, 0, "")
	}
	pdf.Ln(-1)
	pdf.SetFont("Roboto", "", TableHeaderFontSize)

	for _, i := range landingWeightAndBalance.Moments {
		pdf.SetX((210 - (colWidths[0] + colWidths[1] + colWidths[2] + colWidths[3])) / 2)
		mass := strings.ReplaceAll(fmt.Sprintf("%.2f", i.Mass.Kilo()), ".", ",")
		if i.Name == "Fuel" {
			mass = fmt.Sprintf("(%s) %s", parseNumber(fuelPlanning.Total.Volume.Subtract(fuelPlanning.Trip.Volume).String(fuelPlanning.VolumeType)), mass)
		}
		pdf.CellFormat(colWidths[0], 6, strings.ReplaceAll(i.Name, ".", ","), "1", 0, "", false, 0, "")
		pdf.CellFormat(colWidths[1], 6, strings.ReplaceAll(fmt.Sprintf("%.4f", i.Arm), ".", ","), "1", 0, "C", false, 0, "")
		pdf.CellFormat(colWidths[2], 6, strings.ReplaceAll(mass, ".", ","), "1", 0, "C", false, 0, "")
		pdf.CellFormat(colWidths[3], 6, strings.ReplaceAll(fmt.Sprintf("%.2f", i.KGM()), ".", ","), "1", 0, "C", false, 0, "")
		pdf.Ln(-1)
	}

	totalLanding := landingWeightAndBalance.Total
	pdf.SetFont("Roboto", "B", TableHeaderFontSize)
	pdf.SetFillColor(170, 170, 170)
	pdf.SetX((210 - (colWidths[0] + colWidths[1] + colWidths[2] + colWidths[3])) / 2)
	pdf.CellFormat(colWidths[0], 6, "TOTAL", "1", 0, "", true, 0, "")
	pdf.CellFormat(colWidths[1], 6, strings.ReplaceAll(fmt.Sprintf("%.4f", totalLanding.Arm), ".", ","), "1", 0, "C", true, 0, "")
	pdf.CellFormat(colWidths[2], 6, strings.ReplaceAll(fmt.Sprintf("%.2f", totalLanding.Mass.Kilo()), ".", ","), "1", 0, "C", true, 0, "")
	pdf.CellFormat(colWidths[3], 6, strings.ReplaceAll(fmt.Sprintf("%.2f", totalLanding.KGM()), ".", ","), "1", 0, "C", true, 0, "")

	pdf.AddPage()
	pdf.SetFont("Roboto", "B", TitleFontSize)
	pdf.CellFormat(0, TableHeaderFontSize, "Brandstofplanning", "", 1, "C", false, 0, "")
	pdf.Ln(1)

	if !fuelPlanning.Sufficient {
		pdf.SetFont("Roboto", "B", 10)
		pdf.SetFillColor(200, 0, 0)
		pdf.SetTextColor(255, 255, 255)
		pdf.SetX(25)
		pdf.MultiCell(160, 6, "De huidige brandstofvoorraad van "+parseNumber(fuelPlanning.Total.Volume.String(fuelPlanning.VolumeType))+" is onvoldoende om de geplande vlucht veilig uit te voeren. Er moet minimaal "+parseNumber(fuelPlanning.Extra.Volume.Abs().String(fuelPlanning.VolumeType))+" extra brandstof worden bijgetankt om te voldoen aan de veiligheidsvoorschriften.", "1", "C", true)
		pdf.SetTextColor(0, 0, 0)
		pdf.Ln(5)
	}

	fuelRows := []struct {
		label string
		value string
	}{
		{"Taxi Brandstof", parseNumber(fuelPlanning.Taxi.Volume.String(fuelPlanning.VolumeType))},
		{"Reisbrandstof (17L/H)", parseNumber(fuelPlanning.Trip.Volume.String(fuelPlanning.VolumeType))},
		{"Onvoorziene brandstof (10%)", parseNumber(fuelPlanning.Contingency.Volume.String(fuelPlanning.VolumeType))},
		{"Brandstof alternatieve luchthaven", parseNumber(fuelPlanning.Alternate.Volume.String(fuelPlanning.VolumeType))},
		{"Eindreservebrandstof (45 minuten)", parseNumber(fuelPlanning.Reserve.Volume.String(fuelPlanning.VolumeType))},
		{"Extra brandstof", parseNumber(fuelPlanning.Extra.Volume.String(fuelPlanning.VolumeType))},
	}

	pdf.SetFont("Roboto", "B", TableHeaderFontSize)
	pdf.SetFillColor(170, 170, 170)
	pdf.SetX((210 - 160) / 2)
	pdf.CellFormat(120, 6, "Brandstofcategorie", "1", 0, "L", true, 0, "")
	pdf.CellFormat(40, 6, "Hoeveelheid", "1", 0, "C", true, 0, "")
	pdf.Ln(-1)
	pdf.SetFont("Roboto", "", TableHeaderFontSize)
	for _, row := range fuelRows {
		pdf.SetX((210 - 160) / 2)
		pdf.CellFormat(120, 6, row.label, "1", 0, "L", false, 0, "")
		pdf.CellFormat(40, 6, row.value, "1", 0, "C", false, 0, "")
		pdf.Ln(-1)
	}
	pdf.SetFillColor(170, 170, 170)
	pdf.SetFont("Roboto", "B", TableHeaderFontSize)
	pdf.SetX((210 - 160) / 2)
	pdf.CellFormat(120, 6, "Totaal", "1", 0, "L", true, 0, "")
	pdf.CellFormat(40, 6, parseNumber(fuelPlanning.Total.Volume.String(fuelPlanning.VolumeType)), "1", 0, "C", true, 0, "")

	pdf.Ln(10)

	pdf.SetFont("Roboto", "B", TitleFontSize)
	pdf.CellFormat(0, 10, "Prestaties", "", 1, "C", false, 0, "")
	pdf.Ln(2)

	if !takeOffWeightAndBalance.WithinLimits {
		pdf.SetFillColor(200, 0, 0)
		pdf.SetTextColor(255, 255, 255)
		pdf.SetFont("Roboto", "B", 10)
		pdf.SetX(25)
		pdf.MultiCell(160, 4, "De prestaties kunnen niet worden berekend omdat de huidige gewichts- en balansberekening aangeeft dat de belading van het vliegtuig buiten de toegestane limieten valt. Controleer en herbereken de gewichts- en balansverdeling zorgvuldig om te voldoen aan de veiligheidsvoorschriften", "1", "C", true)
		pdf.SetTextColor(0, 0, 0)
	} else {
		pdf.SetFont("Roboto", "B", TableHeaderFontSize)
		pdf.SetFillColor(170, 170, 170)
		pdf.SetX((210 - 160) / 2)
		pdf.CellFormat(120, 6, "Name", "1", 0, "L", true, 0, "")
		pdf.CellFormat(40, 6, "Distance [m]", "1", 0, "C", true, 0, "")
		pdf.Ln(-1)
		pdf.SetFont("Roboto", "", TableHeaderFontSize)
		table := []struct{ Label, Value string }{
			{"Take-off Run Required (Ground Roll)", fmt.Sprintf("%s", strings.ReplaceAll(fmt.Sprintf("%.0f", performance.TakeOffRunRequired), ".", ","))},
			{"Take-off Distance Required", fmt.Sprintf("%s", strings.ReplaceAll(fmt.Sprintf("%.0f", performance.TakeOffDistanceRequired), ".", ","))},
			{"Landing Distance Required", fmt.Sprintf("%s", strings.ReplaceAll(fmt.Sprintf("%.0f", performance.LandingDistanceRequired), ".", ","))},
			{"Landing Ground Roll Required", fmt.Sprintf("%s", strings.ReplaceAll(fmt.Sprintf("%.0f", performance.LandingGroundRollRequired), ".", ","))},
		}
		for _, row := range table {
			pdf.SetX((210 - 160) / 2)
			pdf.CellFormat(120, 6, row.Label, "1", 0, "L", false, 0, "")
			pdf.CellFormat(40, 6, row.Value, "1", 0, "C", false, 0, "")
			pdf.Ln(-1)
		}
		pdf.AddPage()
		pdf.Ln(10)
		pdf.ImageOptions("tdr", 25, pdf.GetY(), 160, 0, false, tdrImageOptions, 0, "")
		pdf.Ln(130)
		pdf.ImageOptions("ldr", 25, pdf.GetY(), 160, 0, false, ldrImageOptions, 0, "")
	}
	var buf bytes.Buffer
	err = pdf.Output(&buf)
	if err != nil {
		return nil, err
	}

	return &buf, nil
}
