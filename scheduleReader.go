package main

import (
	"github.com/tealeg/xlsx"
	"log"
	"regexp"
	"strconv"
	"time"
)

type Lesson struct {
	Name        string
	ClassRoom   string
	TeacherName string
}

type Schedule [7]map[int]Lesson

func (s Schedule) String() string {
	out := "*---Понедельник---*"
	for i := 0; i < 9; i++ {
		l, ok := s[time.Monday][i]
		if !ok {
			continue
		}
		out += "\r\n" + strconv.Itoa(i) + ". " + l.Name + " в " + l.ClassRoom
	}
	out += "\r\n\r\n*---Вторник---*"
	for i := 0; i < 9; i++ {
		l, ok := s[time.Tuesday][i]
		if !ok {
			continue
		}
		out += "\r\n" + strconv.Itoa(i) + ". " + l.Name + " в " + l.ClassRoom
	}
	out += "\r\n\r\n*---Среда---*"
	for i := 0; i < 9; i++ {
		l, ok := s[time.Wednesday][i]
		if !ok {
			continue
		}
		out += "\r\n" + strconv.Itoa(i) + ". " + l.Name + " в " + l.ClassRoom
	}
	out += "\r\n\r\n*---Четверг---*"
	for i := 0; i < 9; i++ {
		l, ok := s[time.Thursday][i]
		if !ok {
			continue
		}
		out += "\r\n" + strconv.Itoa(i) + ". " + l.Name + " в " + l.ClassRoom
	}
	out += "\r\n\r\n*---Пятница---*"
	for i := 0; i < 9; i++ {
		l, ok := s[time.Friday][i]
		if !ok {
			continue
		}
		out += "\r\n" + strconv.Itoa(i) + ". " + l.Name + " в " + l.ClassRoom
	}

	return out
}

func GetSchedule() *Schedule {
	log.Print("Opening and checking \"schedule.xlsx\"...")
	wb, err := xlsx.OpenFile(workDir + "schedule.xlsx")
	if err != nil {
		log.Fatal(err)
	}
	if len(wb.Sheets) == 0 {
		log.Fatal("Workbook doesn't content any sheet")
	}

	sheet := wb.Sheets[0]
	if !checkSheet(sheet) {
		log.Fatal("Sheet has unknown format")
	}

	return getOurSchedule(sheet)
}

func checkSheet(s *xlsx.Sheet) bool {
	// Проверяем первую строку: на ней должны быть названия классов с третьего столбца
	row, err := s.Row(0)
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; ; i++ {
		matched, err := regexp.MatchString(`^\d\d.\d|Дни|Уроки$`, row.GetCell(i).Value)
		if err != nil {
			log.Fatal(err)
		}

		if i == 2 && !matched {
			log.Print("Test 1 failed")
			return false
		}

		if !matched {
			log.Printf("Test 1 passed, columns:%d", i)
			break
		}
	}

	// Проверяем колонку с номерами уроков
	for i := 1; i < 10; i++ {
		row, err := s.Row(i)
		if err != nil {
			log.Fatal(err)
		}

		matched, err := regexp.MatchString(`^[0-8]$`, row.GetCell(1).Value)
		if err != nil {
			log.Fatal(err)
		}

		if !matched {
			log.Print("Test 2 failed")
			return false
		}
	}
	log.Print("Test 2 passed")
	return true
}

func getOurSchedule(s *xlsx.Sheet) *Schedule {
	// Id столбца нашего
	var colId int
	sch := new(Schedule)
	sch[time.Monday] = make(map[int]Lesson)
	sch[time.Tuesday] = make(map[int]Lesson)
	sch[time.Wednesday] = make(map[int]Lesson)
	sch[time.Thursday] = make(map[int]Lesson)
	sch[time.Friday] = make(map[int]Lesson)

	row, _ := s.Row(0)
	for i := 2; i < 1000; i++ {
		if row.GetCell(i).Value == className {
			colId = i
			break
		}
		if i == 999 {
			log.Fatal("Class' column not found")
		}
	}

	// Сканируем дни недели
	for startRow := 1; startRow <= 41; startRow += 10 {
		for i := startRow; i < startRow+9; i++ {
			cell, _ := s.Cell(i, colId)
			if cell.Value == "" {
				continue
			}

			lesson := parseLesson(cell.Value)
			switch i / 10 {
			case 0:
				sch[time.Monday][i-startRow] = lesson
			case 1:
				sch[time.Tuesday][i-startRow] = lesson
			case 2:
				sch[time.Wednesday][i-startRow] = lesson
			case 3:
				sch[time.Thursday][i-startRow] = lesson
			case 4:
				sch[time.Friday][i-startRow] = lesson
			}
		}
	}

	return sch
}

func parseLesson(s string) Lesson {
	re := regexp.MustCompile(`^'?(.*)(?:\r\n|\r|\n)(.*)(?:\r\n|\r|\n)(.*)$`)
	result := re.FindAllStringSubmatch(s, -1)

	if len(result) == 0 || len(result[0]) < 4 {
		log.Fatalf("Не удалось распарсить урок: %s", s)
	}

	resultL := result[0]
	return Lesson{
		Name:        resultL[1],
		ClassRoom:   resultL[3],
		TeacherName: resultL[2],
	}
}
