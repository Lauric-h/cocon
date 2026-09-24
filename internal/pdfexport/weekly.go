package pdfexport

import (
	"bytes"
	"cocon/internal/handlers/view"
	"fmt"

	"codeberg.org/go-pdf/fpdf"
)

const (
	marginX    = 15.0
	colGap     = 8.0
	leftWidth  = 110.0
	rightWidth = 62.0 // 180 (largeur utile A4) - leftWidth - colGap
)

func BuildWeeklySummary(weekRangeLabel string, days []view.WeeklySummaryDay, tasks []view.TaskItem) ([]byte, error) {
	pdf := fpdf.New("P", "mm", "A4", "")
	tr := pdf.UnicodeTranslatorFromDescriptor("")
	pdf.SetMargins(marginX, marginX, marginX)
	pdf.AddPage()

	pdf.SetFont("Helvetica", "B", 16)
	pdf.CellFormat(0, 10, tr("Semaine du "+weekRangeLabel), "", 1, "C", false, 0, "")

	leftX := marginX
	rightX := marginX + leftWidth + colGap
	y := pdf.GetY() + 6

	leftY := writeDaysColumn(pdf, tr, leftX, y, days)
	rightY := writeTasksColumn(pdf, tr, rightX, y, tasks)

	_ = leftY
	_ = rightY

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// writeDaysColumn écrit la colonne des jours à partir de (x, y), et renvoie le y final atteint.
func writeDaysColumn(pdf *fpdf.Fpdf, tr func(string) string, x, y float64, days []view.WeeklySummaryDay) float64 {
	for _, day := range days {
		pdf.SetXY(x, y)
		pdf.SetFont("Helvetica", "B", 11)
		pdf.CellFormat(leftWidth, 6, tr(day.DayName+" "+day.DateLabel), "", 0, "L", false, 0, "")
		y += 6

		for _, ev := range day.Events {
			pdf.SetXY(x, y)
			pdf.SetFont("Helvetica", "B", 9)
			pdf.SetTextColor(196, 96, 58) // terracotta
			pdf.CellFormat(leftWidth, 5, tr(ev.Title+" @ "+ev.TimeLabel), "", 0, "L", false, 0, "")
			pdf.SetTextColor(0, 0, 0)
			y += 5
		}

		pdf.SetFont("Helvetica", "", 10)
		for _, line := range []string{
			"Matin : " + orDash(day.Breakfast),
			"Midi : " + orDash(day.Lunch),
			"Soir : " + orDash(day.Dinner),
		} {
			pdf.SetXY(x, y)
			pdf.CellFormat(leftWidth, 5, tr(line), "", 0, "L", false, 0, "")
			y += 5
		}
		y += 4 // espace entre deux jours
	}
	return y
}

// writeTasksColumn écrit la colonne des tâches à partir de (x, y), et renvoie le y final atteint.
func writeTasksColumn(pdf *fpdf.Fpdf, tr func(string) string, x, y float64, tasks []view.TaskItem) float64 {
	pdf.SetXY(x, y)
	pdf.SetFont("Helvetica", "B", 12)
	pdf.CellFormat(rightWidth, 7, tr("Tâches"), "", 0, "L", false, 0, "")
	y += 8

	if len(tasks) == 0 {
		pdf.SetXY(x, y)
		pdf.SetFont("Helvetica", "", 9)
		pdf.CellFormat(rightWidth, 5, tr("Aucune tâche planifiée."), "", 0, "L", false, 0, "")
		return y + 5
	}

	pdf.SetFont("Helvetica", "", 9)
	for _, t := range tasks {
		status := "[ ]"
		if t.Done {
			status = "[x]"
		}
		label := fmt.Sprintf("%s %s", status, t.Title)
		pdf.SetXY(x, y)
		pdf.CellFormat(rightWidth, 5, tr(label), "", 0, "L", false, 0, "")
		y += 4

		pdf.SetXY(x, y)
		pdf.SetFont("Helvetica", "", 8)
		pdf.SetTextColor(122, 110, 96)
		pdf.CellFormat(rightWidth, 4, tr(t.DueDateLabel+" · "+t.CategoryLabel), "", 0, "L", false, 0, "")
		pdf.SetTextColor(0, 0, 0)
		pdf.SetFont("Helvetica", "", 9)
		y += 6
	}
	return y
}

func orDash(s string) string {
	if s == "" {
		return "-"
	}
	return s
}
