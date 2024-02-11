package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/xuri/excelize/v2"
)

func GetItems(weekday, meal, file string) ([]string, error) {

	f, err := excelize.OpenFile(file)
	if err != nil {
		return nil, err
	}

	rows, err := f.GetRows("Sheet1")
	if err != nil {
		return nil, err
	}

	weekdayIndex := -1
	for i, cell := range rows[0] {
		if strings.EqualFold(cell, weekday) {
			weekdayIndex = i + 1
			break
		}
	}

	if weekdayIndex == -1 {
		return nil, fmt.Errorf("weekday not found %s", weekday)
	}

	mealIndex := -1
	for i, row := range rows {
		cellValue := row[weekdayIndex-1]
		if strings.EqualFold(cellValue, meal) {
			mealIndex = i
			break
		}
	}

	if mealIndex == -1 {
		return nil, fmt.Errorf("meal not found %s", meal)
	}

	var items []string
	for i := mealIndex + 1; i < len(rows); i++ {
		cellValue := rows[i][weekdayIndex-1]
		if strings.EqualFold(cellValue, weekday) {
			break
		}

		items = append(items, cellValue)
	}

	return items, nil
}

func GetItemCount(day, meal, file string) (int, error) {
	items, err := GetItems(day, meal, file)
	if err != nil {
		return 0, err
	}

	return len(items), nil
}

func IsItemInMeal(day, meal, item, file string) (bool, error) {
	items, err := GetItems(day, meal, file)
	if err != nil {
		return false, err
	}

	for _, i := range items {
		if strings.EqualFold(i, item) {
			return true, nil
		}
	}

	return false, nil
}

type Meal struct {
	Day   string
	Date  string
	Meal  string
	Items []string
}

var meals []Meal

func (m *Meal) PrintDetails() {
	fmt.Printf("Day: %s\n", m.Day)
	fmt.Printf("Date: %s\n", m.Date)
	fmt.Printf("Meal: %s\n", m.Meal)
	fmt.Println("Items:")
	for _, item := range m.Items {
		fmt.Printf("- %s\n", item)
	}
	fmt.Println()
}

func JsonConv(file string) {

	menu := make(map[string]map[string][]string)

	f, err := excelize.OpenFile(file)
	if err != nil {
		fmt.Println("error in opening excel file:", err)
		return
	}

	cols, err := f.GetCols("Sheet1")
	if err != nil {
		fmt.Println("error in getting columns:", err)
		return
	}

	for _, col := range cols {
		menu[col[0]] = make(map[string][]string)
		j := 3
		for j < len(col) {
			if strings.EqualFold(col[j], "Lunch") || strings.EqualFold(col[j], col[0]) {
				j++
				break
			}
			menu[col[0]]["Breakfast"] = append(menu[col[0]]["Breakfast"], col[j])
			j++
		}
		for j < len(col) {
			if strings.EqualFold(col[j], "Dinner") || strings.EqualFold(col[j], col[0]) {
				j++
				break
			}
			menu[col[0]]["Lunch"] = append(menu[col[0]]["Lunch"], col[j])
			j++
		}
		for j < len(col) {
			if strings.EqualFold(col[j], "") {
				break
			}
			menu[col[0]]["Dinner"] = append(menu[col[0]]["Dinner"], col[j])
			j++
		}
	}

	for _, col := range cols {
		day := col[0]
		date := col[1]
		breakfast := Meal{
			Day:   day,
			Date:  date,
			Meal:  "Breakfast",
			Items: menu[day]["Breakfast"],
		}
		lunch := Meal{
			Day:   day,
			Date:  date,
			Meal:  "Lunch",
			Items: menu[day]["Lunch"],
		}
		dinner := Meal{
			Day:   day,
			Date:  date,
			Meal:  "Dinner",
			Items: menu[day]["Dinner"],
		}

		meals = append(meals, breakfast)
		meals = append(meals, lunch)
		meals = append(meals, dinner)
	}

	jsonFile, err := os.Create("menu.json")
	if err != nil {
		fmt.Println("error in creation of file:", err)
		return
	}
	defer jsonFile.Close()

	encoder := json.NewEncoder(jsonFile)
	encoder.SetIndent("", "    ") // Optional: Set indentation for readability
	if err := encoder.Encode(menu); err != nil {
		fmt.Println("error in encoding :", err)
		return
	}

	fmt.Println("Proceess succeded")

	for _, meal := range meals {
		meal.PrintDetails()
	}

}

func main() {
	filePath := "Sample-Menu.xlsx"
	// weekday := "Monday"
	// meal := "Dinner"

	JsonConv(filePath)

}
