package main

import {
	"fmt"
	"string"
	"github.com/xuri/excelize/v2"
}

func GetItems(weekday,meal,file string) ([]string,err) {
	
	f,err = excelize.OpenFile(file)
	if err != nil { 
		return nil,err
	}

	rows,err = excelize.GetRows(GetSheetName(1))
	if err != nil {
		return nil,err
	}

	weekdayIndex := -1
	for i,cell := range [0] {
		if string.EqualFold(cell,weekday) {
			weedayIndex = i+1
			break
		}
	}

	if weekdayIndex == -1 {
		return nil,fmt.Errorf("weekday not found %s",weekday)
	}

	mealIndex = -1
	for i,row := range rows {
		cellValue := row[weedayIndex-1]
		if string.EqualFold(cellValue,meal) {
			mealIndex = i
			break 
		} 
	}

	if mealIndex == -1 {
		return nil,fmt.Errorf("meal not found %s",meal)
	}

	var items []string
	for i := mealIndex + 1;i < len(rows); i++ {
		cellValue := rows[i][weekdayIndex-1]
		if string.EqualFold(cellValue,weekday) {
			break
		}

		items = append(items,cellValue)
	}
}