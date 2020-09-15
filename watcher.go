package main

import (
	"log"
	"time"
)

var loc = time.FixedZone("", timezoneOffset)

func everyDay(s *Schedule, t *Timetable) {
	for {
		now := time.Now().Add(offset).In(loc)
		to := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 1, 0, now.Location())
		log.Printf("Waiting for %s", to.String())

		time.Sleep(to.Sub(time.Now().Add(offset)))
		completeThisDay(s, t)
	}
}

func completeThisDay(s *Schedule, t *Timetable) {
	now := time.Now().Add(offset).In(loc)
	log.Printf("Now is %s", now.String())
	weekday := now.Weekday()

	if weekday == time.Sunday || weekday == time.Saturday {
		log.Println("Today nothing to do :)")
	}

	for i := 0; i < 9; i++ {
		val, ok := s[weekday][i]
		if !ok {
			continue
		}

		// val - lesson, i - num

		if i == 0 {
			to := time.Date(now.Year(), now.Month(), now.Day(), t[0].hour, t[0].minute-15, 0, 0, loc)
			go doWhen(func() {
				alertAboutLesson(val, t, i)
			}, to)

			continue
		}

		to := time.Date(now.Year(), now.Month(), now.Day(), t[i-1].hour, t[i-1].minute+40, 0, 0, loc)

		h := i

		go doWhen(func() {
			alertAboutLesson(val, t, h)
		}, to)
	}
}

func doWhen(f func(), t time.Time) {
	log.Printf("Set function delay until %s", t.String())

	time.Sleep(t.Sub(time.Now().Add(offset)))
	log.Printf("Now is %s", time.Now().Add(offset).In(loc).String())
	f()
}

func alertAboutLesson(l Lesson, t *Timetable, n int) {
	msg := "Следующий урок – " + l.Name + "\r\n" +
		"Кабинет: " + l.ClassRoom + "\r\n" +
		"Учитель: " + l.TeacherName + "\r\n" +
		"Начало урока: " + t[n].toString()

	sendToAll(msg)
}
