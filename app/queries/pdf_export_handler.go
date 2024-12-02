package queries

import (
	"bytes"
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"html/template"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"ppl-calculations/domain/calculations"
	"ppl-calculations/domain/fuel"
	"ppl-calculations/domain/state"
	"strings"
	"time"
)

type PdfExportHandler struct {
	template    bytes.Buffer
	calcService calculations.Service
	tmpFolder   string
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

func NewPdfExportHandler(assets fs.FS, calcService calculations.Service) PdfExportHandler {
	e, err := assets.Open("export.tex")
	if err != nil {
		logrus.WithError(err).Fatal("export template not found")
	}

	var exportBuf bytes.Buffer
	_, err = io.Copy(&exportBuf, e)
	if err != nil {
		logrus.WithError(err).Fatal("export template copy")
	}

	if err := e.Close(); err != nil {
		logrus.WithError(err).Fatal("export template closing")
	}

	return PdfExportHandler{
		calcService: calcService,
		tmpFolder:   os.Getenv("TMP_PATH"),
		template:    exportBuf,
	}
}

func parseNumber(number string) string {
	return strings.ReplaceAll(number, ".", ",")
}

func (h PdfExportHandler) Handle(ctx context.Context, stateService state.Service, reference string) (io.Reader, error) {
	s, err := stateService.State(ctx)
	if err != nil {
		return nil, err
	}

	if s.TripDuration == nil {
		return nil, ErrMissingFuelSheet
	}

	if s.CallSign == nil {
		return nil, ErrMissingFuelSheet
	}

	var f fuel.Fuel
	var takeOffWeightAndBalance *calculations.WeightBalance
	if s.MaxFuel != nil && *s.MaxFuel {
		takeOffWeightAndBalance, f, err = calculations.NewWeightAndBalanceMaxFuel(*s.CallSign, *s.Pilot, *s.PilotSeat, s.Passenger, s.PassengerSeat, *s.Baggage, *s.FuelType)
		if err != nil {
			return nil, err
		}
	} else {
		f = *s.Fuel
		takeOffWeightAndBalance, err = calculations.NewWeightAndBalance(*s.CallSign, *s.Pilot, *s.PilotSeat, s.Passenger, s.PassengerSeat, *s.Baggage, f)
		if err != nil {
			return nil, err
		}
	}

	fuelPlanning, err := calculations.NewFuelPlanning(*s.TripDuration, *s.AlternateDuration, f, *s.FuelVolumeType)
	if err != nil {
		return nil, err
	}

	landingWeightAndBalance, err := calculations.NewWeightAndBalance(*s.CallSign, *s.Pilot, *s.PilotSeat, s.Passenger, s.PassengerSeat, *s.Baggage, fuel.Subtract(f, fuelPlanning.Trip, fuelPlanning.Taxi))
	if err != nil {
		return nil, err
	}

	tdrChart, todRR, todDR, err := h.calcService.TakeOffDistance(*s.OutsideAirTemperature, *s.PressureAltitude, takeOffWeightAndBalance.Total.Mass, *s.Wind, calculations.ChartTypePNG)
	if err != nil {
		return nil, err
	}

	ldrChart, ldrDR, ldrGR, err := h.calcService.LandingDistance(*s.OutsideAirTemperature, *s.PressureAltitude, landingWeightAndBalance.Total.Mass, *s.Wind, calculations.ChartTypePNG)
	if err != nil {
		return nil, err
	}

	wbChart, err := h.calcService.WeightAndBalance(*s.CallSign, *takeOffWeightAndBalance.Total, *landingWeightAndBalance.Total, takeOffWeightAndBalance.WithinLimits, calculations.ChartTypePNG)
	if err != nil {
		return nil, err
	}

	performance := &calculations.Performance{
		TakeOffRunRequired:        todRR,
		TakeOffDistanceRequired:   todDR,
		LandingDistanceRequired:   ldrDR,
		LandingGroundRollRequired: ldrGR,
	}

	tdrBytes, err := io.ReadAll(tdrChart)
	if err != nil {
		return nil, err
	}

	ldrBytes, err := io.ReadAll(ldrChart)
	if err != nil {
		return nil, err
	}

	wbBytes, err := io.ReadAll(wbChart)
	if err != nil {
		return nil, err
	}

	tempDir, err := os.MkdirTemp(h.tmpFolder, "download.*")
	if err != nil {
		return nil, err
	}
	defer func() {
		err := os.RemoveAll(tempDir)
		if err != nil {
			logrus.WithError(err).Error("removing temporary directory")
		}
	}()

	err = os.WriteFile(path.Join(tempDir, "tdr.png"), tdrBytes, 0644)
	if err != nil {
		return nil, err
	}

	err = os.WriteFile(path.Join(tempDir, "ldr.png"), ldrBytes, 0644)
	if err != nil {
		return nil, err
	}

	err = os.WriteFile(path.Join(tempDir, "wb.png"), wbBytes, 0644)
	if err != nil {
		return nil, err
	}

	tmpl, err := template.New("pdfTemplate").Parse(h.template.String())
	if err != nil {
		return nil, err
	}

	exData := ExportData{}

	exData.CallSign = s.CallSign.String()
	exData.Generated = time.Now().Format("15:04:05 02-01-2006")
	exData.Reference = reference

	exData.FuelSufficient = fuelPlanning.Sufficient
	exData.FuelTaxi = parseNumber(fuelPlanning.Taxi.Volume.String(fuelPlanning.VolumeType))
	exData.FuelTrip = parseNumber(fuelPlanning.Trip.Volume.String(fuelPlanning.VolumeType))
	exData.FuelAlternate = parseNumber(fuelPlanning.Alternate.Volume.String(fuelPlanning.VolumeType))
	exData.FuelContingency = parseNumber(fuelPlanning.Contingency.Volume.String(fuelPlanning.VolumeType))
	exData.FuelReserve = parseNumber(fuelPlanning.Reserve.Volume.String(fuelPlanning.VolumeType))
	exData.FuelTotal = parseNumber(fuelPlanning.Total.Volume.String(fuelPlanning.VolumeType))
	exData.FuelExtra = parseNumber(fuelPlanning.Extra.Volume.String(fuelPlanning.VolumeType))
	exData.FuelExtraAbs = parseNumber(fuelPlanning.Extra.Volume.Abs().String(fuelPlanning.VolumeType))

	wbState := WeightAndBalanceState{}
	for _, i := range takeOffWeightAndBalance.Moments {
		m := parseNumber(fmt.Sprintf("%.2f", i.Mass.Kilo()))
		if strings.HasPrefix(i.Name, "Fuel") {
			m = fmt.Sprintf("(%s) %s", exData.FuelTotal, m)
		}

		wbState.Items = append(wbState.Items, WeightAndBalanceItem{
			Name:       parseNumber(i.Name),
			LeverArm:   parseNumber(fmt.Sprintf("%.4f", i.Arm)),
			Mass:       m,
			MassMoment: parseNumber(fmt.Sprintf("%.2f", i.KGM())),
		})
	}

	wbState.Total = WeightAndBalanceItem{
		Name:       parseNumber(takeOffWeightAndBalance.Total.Name),
		LeverArm:   parseNumber(fmt.Sprintf("%.4f", takeOffWeightAndBalance.Total.Arm)),
		Mass:       parseNumber(fmt.Sprintf("%.2f", takeOffWeightAndBalance.Total.Mass.Kilo())),
		MassMoment: parseNumber(fmt.Sprintf("%.2f", takeOffWeightAndBalance.Total.KGM())),
	}

	wbState.WithinLimits = takeOffWeightAndBalance.WithinLimits

	wbLandingState := WeightAndBalanceState{}

	for _, i := range landingWeightAndBalance.Moments {
		m := parseNumber(fmt.Sprintf("%.2f", i.Mass.Kilo()))
		if strings.HasPrefix(i.Name, "Fuel") {
			m = fmt.Sprintf("(%s) %s", parseNumber(fuelPlanning.Total.Volume.Subtract(fuelPlanning.Trip.Volume).String(fuelPlanning.VolumeType)), m)
		}

		wbLandingState.Items = append(wbLandingState.Items, WeightAndBalanceItem{
			Name:       parseNumber(i.Name),
			LeverArm:   parseNumber(fmt.Sprintf("%.4f", i.Arm)),
			Mass:       m,
			MassMoment: parseNumber(fmt.Sprintf("%.2f", i.KGM())),
		})
	}

	wbLandingState.Total = WeightAndBalanceItem{
		Name:       parseNumber(landingWeightAndBalance.Total.Name),
		LeverArm:   parseNumber(fmt.Sprintf("%.4f", landingWeightAndBalance.Total.Arm)),
		Mass:       parseNumber(fmt.Sprintf("%.2f", landingWeightAndBalance.Total.Mass.Kilo())),
		MassMoment: parseNumber(fmt.Sprintf("%.2f", landingWeightAndBalance.Total.KGM())),
	}

	wbLandingState.WithinLimits = takeOffWeightAndBalance.WithinLimits

	exData.TakeOffDistanceRequired = fmt.Sprintf("%.0f", performance.TakeOffDistanceRequired)
	exData.TakeOffRunRequired = fmt.Sprintf("%.0f", performance.TakeOffRunRequired)
	exData.LandingDistanceRequired = fmt.Sprintf("%.0f", performance.LandingDistanceRequired)
	exData.LandingGroundRollRequired = fmt.Sprintf("%.0f", performance.LandingGroundRollRequired)

	exData.WeightAndBalanceTakeOff = wbState
	exData.WeightAndBalanceLanding = wbLandingState

	var output bytes.Buffer
	err = tmpl.Execute(&output, exData)
	if err != nil {
		return nil, err
	}

	err = os.WriteFile(path.Join(tempDir, "export.tex"), output.Bytes(), 0644)
	if err != nil {
		return nil, err
	}

	cmd := exec.Command("xelatex", "-halt-on-error", "-interaction=nonstopmode", "export.tex")
	cmd.Dir = tempDir

	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	err = cmd.Run()
	if err != nil {
		fmt.Println(err)
		fmt.Println(stdoutBuf.String())
		fmt.Println(stderrBuf.String())
		return nil, err
	}

	pdfPath := filepath.Join(tempDir, "export.pdf")

	pdfData, err := os.ReadFile(pdfPath)
	if err != nil {
		return nil, err
	}

	return bytes.NewBufferString(string(pdfData)), nil
}
