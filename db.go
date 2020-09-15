package main

import (
	"github.com/mxk/go-sqlite/sqlite3"
	"io"
	"log"
	"strconv"
)

var db *sqlite3.Conn

func connectToDB() {
	db, err = sqlite3.Open(workDir + "database.db")
	if err != nil {
		log.Fatal(err)
	}

	err = db.Exec("CREATE TABLE IF NOT EXISTS chats (id INTEGER PRIMARY KEY AUTOINCREMENT, chatID INTEGER UNIQUE)")
	if err != nil {
		log.Fatal(err)
	}
}

func isChatEnabled(id int64) bool {
	_, err := db.Query("SELECT chatID FROM chats WHERE chatId = " + strconv.FormatInt(id, 10))
	if err == io.EOF {
		return false
	}
	if err != nil {
		log.Fatal(err)
	}

	return true
}

func disableChat(id int64) {
	log.Printf("Disabling chat %d", id)
	err := db.Exec("DELETE FROM chats WHERE chatId = " + strconv.FormatInt(id, 10))
	if err != nil {
		log.Fatal(err)
	}
}

func enableChat(id int64) {
	log.Printf("Enabling chat %d", id)
	err := db.Exec("INSERT INTO chats(chatID) VALUES (" + strconv.FormatInt(id, 10) + ")")
	if err != nil {
		log.Fatal(err)
	}
}

func getEnabledChats() []int64 {
	var chats []int64

	for s, err := db.Query("SELECT chatID FROM chats"); err == nil; err = s.Next() {
		var chatID int64

		err := s.Scan(&chatID)
		if err != nil {
			log.Fatal(err)
		}

		chats = append(chats, chatID)
	}

	return chats
}
