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

  //create users table 
  query = `CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    user_name VARCHAR(50) NOT NULL,
    banned BOOLEAN DEFAULT FALSE
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

func Add(msgID, chatID, usrID int64, usrName string) {
  //add to takes
  query := `INSERT INTO takes (msg_id, chat_id) VALUES ($1, $2)`
  _, err := db.Exec(query, msgID, chatID)
  if err != nil {
    log.Fatal("cannot add to table", err)
  }

  //add to users 
  query = `INSERT INTO users (user_id, user_name) VALUES ($1, $2)`
  _, err = db.Exec(query, usrID, usrName)
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

func Ban(usrID int64) {
  query := `UPDATE users SET banned = true WHERE user_id = $1`
  _, err := db.Exec(query, usrID)
  if err != nil {
    log.Fatal("cannot ban user", err)
  }

  log.Println("banned user", usrID)
}

func UnBan(usrID int64) {
  query := `UPDATE users SET banned = false WHERE user_id = $1`
  _, err := db.Exec(query, usrID)
  if err != nil {
    log.Fatal("cannot unban user", err)
  }

  log.Println("unbanned user", usrID)
}

func CheckBan(usrID int64) bool {
  var banned bool
  query := `SELECT banned FROM users WHERE user_id = $1`
  err := db.QueryRow(query, usrID).Scan(&banned)
  if err != nil {
    log.Fatal("cannot get from table", err)
  }

  return banned
}

func Close() {
  if db != nil {
    db.Close()
  }
}
