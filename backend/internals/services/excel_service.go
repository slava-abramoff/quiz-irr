package services

import (
	"context"
	"fmt"
	"io"
	"quiz-irr/internals/storage/models"
	"sort"

	"github.com/xuri/excelize/v2"
)

type excelService struct{}

func NewExcelService() *excelService {
	return &excelService{}
}

func (e *excelService) MakeResults(
	ctx context.Context,
	writer io.Writer,
	results []models.TestResult,
) error {
	f := excelize.NewFile()
	defer func() {
		_ = f.Close()
	}()

	sheet := "Results"
	if err := f.SetSheetName("Sheet1", sheet); err != nil {
		return fmt.Errorf("excel set sheet name: %w", err)
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].TotalScore > results[j].TotalScore
	})

	headers := []string{
		"ФИО",
		"Почта",
		"Учреждение",
		"Время",
		"В срок",
		"Итого",
	}

	for i, h := range headers {
		cell, err := excelize.CoordinatesToCellName(i+1, 1)
		if err != nil {
			return fmt.Errorf("excel header cell name: %w", err)
		}
		if err := f.SetCellValue(sheet, cell, h); err != nil {
			return fmt.Errorf("excel set header: %w", err)
		}
	}

	headerStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
	})
	if err != nil {
		return fmt.Errorf("excel header style: %w", err)
	}

	lastHeaderCell, err := excelize.CoordinatesToCellName(len(headers), 1)
	if err != nil {
		return fmt.Errorf("excel last header cell: %w", err)
	}
	if err := f.SetCellStyle(sheet, "A1", lastHeaderCell, headerStyle); err != nil {
		return fmt.Errorf("excel header cell style: %w", err)
	}

	style := "mm:ss"
	timeStyle, err := f.NewStyle(&excelize.Style{
		CustomNumFmt: &style,
	})
	if err != nil {
		return fmt.Errorf("excel time style: %w", err)
	}

	for i, r := range results {
		row := i + 2

		values := []interface{}{
			r.FullName,
			r.Email,
			r.Org,
			float64(r.Duration) / 86400,
			r.IsOnTime,
			r.TotalScore,
		}

		for j, v := range values {
			cell, err := excelize.CoordinatesToCellName(j+1, row)
			if err != nil {
				return fmt.Errorf("excel data cell name: %w", err)
			}
			if err := f.SetCellValue(sheet, cell, v); err != nil {
				return fmt.Errorf("excel set cell: %w", err)
			}
		}

		timeCell, err := excelize.CoordinatesToCellName(4, row)
		if err != nil {
			return fmt.Errorf("excel time cell name: %w", err)
		}
		if err := f.SetCellStyle(sheet, timeCell, timeCell, timeStyle); err != nil {
			return fmt.Errorf("excel time cell style: %w", err)
		}
	}

	for i := range headers {
		col, err := excelize.ColumnNumberToName(i + 1)
		if err != nil {
			return fmt.Errorf("excel column name: %w", err)
		}
		if err := f.SetColWidth(sheet, col, col, 22); err != nil {
			return fmt.Errorf("excel column width: %w", err)
		}
	}

	scoreCol := "F"
	lastRow := len(results) + 1

	if lastRow >= 2 {
		if err := f.SetConditionalFormat(sheet,
			fmt.Sprintf("%s2:%s%d", scoreCol, scoreCol, lastRow),
			[]excelize.ConditionalFormatOptions{
				{
					Type:     "3_color_scale",
					Criteria: "=",
					MinType:  "min",
					MidType:  "percentile",
					MidValue: "50",
					MaxType:  "max",
					MinColor: "#F8696B",
					MidColor: "#FFEB84",
					MaxColor: "#63BE7B",
				},
			},
		); err != nil {
			return fmt.Errorf("excel conditional format (scale): %w", err)
		}

		topStyleID, err := f.NewStyle(&excelize.Style{
			Font: &excelize.Font{
				Bold:  true,
				Color: "#000000",
			},
			Fill: excelize.Fill{
				Type:    "pattern",
				Color:   []string{"#FFD700"},
				Pattern: 1,
			},
		})
		if err != nil {
			return fmt.Errorf("excel top style: %w", err)
		}
		if err := f.SetConditionalFormat(sheet,
			fmt.Sprintf("%s2:%s%d", scoreCol, scoreCol, lastRow),
			[]excelize.ConditionalFormatOptions{
				{
					Type:     "top",
					Criteria: "=",
					Format:   &topStyleID,
					Value:    "3",
				},
			},
		); err != nil {
			return fmt.Errorf("excel conditional format (top): %w", err)
		}
	}

	endCol, err := excelize.ColumnNumberToName(len(headers))
	if err != nil {
		return fmt.Errorf("excel filter column name: %w", err)
	}
	filterRef := fmt.Sprintf("A1:%s%d", endCol, lastRow)
	if err := f.AutoFilter(sheet, filterRef, nil); err != nil {
		return fmt.Errorf("excel autofilter: %w", err)
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	if err := f.Write(writer); err != nil {
		return fmt.Errorf("excel write: %w", err)
	}

	return nil
}

func (e *excelService) MakeRaws() {}
