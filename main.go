package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"strconv"
	"time"
)

//const configPath = "config.json"
const configPath = "/etc/mytgdsch/config.json"

//var offset = time.Until(time.Date(2020, 9, 14, 8, 54, 55, 0, time.FixedZone("UTC+3", 3*60*60)))
//var offset = time.Until(time.Date(2020, 9, 17, 8, -1, 40, 0, time.FixedZone("UTC+3", 3*60*60)))
var offset = time.Duration(0)

var (
	botToken  string
	botPort   string
	botUrl    string
	botPath   string
	listen    string
	workDir   string
	className string
)

type Timetable [9]TimetableItem

type TimetableItem struct {
	hour   int
	minute int
}

func (t TimetableItem) toString() string {
	out := ""
	if t.hour < 10 {
		out += "0" + strconv.Itoa(t.hour)
	} else {
		out += strconv.Itoa(t.hour)
	}
	out += ":"
	if t.minute < 10 {
		out += "0" + strconv.Itoa(t.minute)
	} else {
		out += strconv.Itoa(t.minute)
	}

	return out
}

func (t TimetableItem) addMinutes(a int) TimetableItem {
	n := new(TimetableItem)
	n.minute = (t.minute + a) % 60
	n.hour = (t.hour + (t.minute+a)/60) % 24

	return *n
}

func main() {
	openConfig()
	s, t := GetSchedule(), getTimetable()
	connectToDB()
	initBotApi()

	completeThisDay(s, t)
	everyDay(s, t)
}

func getTimetable() *Timetable {
	t := Timetable{
		TimetableItem{
			hour:   8,
			minute: 15,
		},
		TimetableItem{
			hour:   9,
			minute: 00,
		},
		TimetableItem{
			hour:   9,
			minute: 50,
		},
		TimetableItem{
			hour:   10,
			minute: 45,
		},
		TimetableItem{
			hour:   11,
			minute: 40,
		},
		TimetableItem{
			hour:   12,
			minute: 40,
		},
		TimetableItem{
			hour:   13,
			minute: 40,
		},
		TimetableItem{
			hour:   14,
			minute: 40,
		},
		TimetableItem{
			hour:   15,
			minute: 30,
		},
	}

	return &t
}

func openConfig() {
	f, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatal(err)
	}

	value := struct {
		Token          string
		Port           int
		Listen         string
		WorkDir        string
		Url            string
		Path           string
		TimezoneOffset int
		ClassName      string
	}{}

	err = json.Unmarshal(f, &value)
	if err != nil {
		log.Fatal(err)
	}

	botToken = value.Token
	botPort = strconv.Itoa(value.Port)
	botUrl = value.Url + "/"
	botPath = value.Path + "/"
	workDir = value.WorkDir
	listen = value.Listen
	loc = time.FixedZone("", value.TimezoneOffset)
	className = value.ClassName
}
