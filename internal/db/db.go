package db

import (
  "database/sql"
  _ "github.com/lib/pq"
  "log"
)

var db *sql.DB

func Init() {
  //connect to database
  var err error
  db, err = sql.Open("postgres", "host=127.0.0.1 port=5432 user=postgres password=123 dbname=postgres sslmode=disable")
  if err != nil {
    log.Fatal("cannot open database", err)
  }

  query := `CREATE TABLE IF NOT EXISTS takes (
    id SERIAL PRIMARY KEY,
    msg_id BIGINT NOT NULL,
    chat_id BIGINT NOT NULL
  )`
  _, err = db.Exec(query)
  if err != nil {
    log.Fatal("cannot create table", err)
  }

  err = db.Ping()
  if err != nil {
    log.Fatal("cannot ping database", err)
  } 
}

func Add(msgID, chatID int64) {
  query := `INSERT INTO takes (msg_id, chat_id) VALUES ($1, $2)`
  _, err := db.Exec(query, msgID, chatID)
  if err != nil {
    log.Fatal("cannot add to table", err)
  }

  log.Println("added to table ", msgID, chatID)
}

func GetByMsgID(msgID int64) (int64, error) {
  var chatID int64
  log.Println(msgID)
  query := `SELECT chat_id FROM takes WHERE msg_id = $1`
  err := db.QueryRow(query, msgID).Scan(&chatID)
  if err != nil {
    log.Fatal("cannot get from table", err)
  }
  return chatID, err
}

func Close() {
  if db != nil {
    db.Close()
  }
}
