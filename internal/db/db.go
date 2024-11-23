package db

import (
  "database/sql"
  _ "github.com/lib/pq"
  "log"
)

var db *sql.DB

func Init(connStr string) {
  //connect to database
  var err error
  db, err = sql.Open("postgres", connStr)
  if err != nil {
    log.Fatal("[DB]: [FATAL]: cannot open database")
  }
  
  //create takes table
  query := `CREATE TABLE IF NOT EXISTS takes (
    id SERIAL PRIMARY KEY,
    msg_id BIGINT NOT NULL,
    chat_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL
  )`
  _, err = db.Exec(query)
  if err != nil {
    log.Fatal("[DB]: [FATAL]: cannot create takes table")
  }

  log.Println("[DB]: [INFO]: takes table created")

  //create users table 
  query = `CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    user_name VARCHAR(50) NOT NULL,
    banned BOOLEAN DEFAULT FALSE
  )`
  _, err = db.Exec(query)
  if err != nil {
    log.Fatal("[DB]: [FATAL]: cannot create users table")
  }

  log.Println("[DB]: [INFO]: users table created")

  err = db.Ping()
  if err != nil {
    log.Fatal("[DB]: [FATAL]: cannot ping database")
  }

  log.Println("[DB]: [INFO]: connected to database")
}

func Add(msgID, chatID, usrID int64, usrName string) {
  //add to takes
  query := `INSERT INTO takes (msg_id, chat_id, user_id) VALUES ($1, $2, $3)`
  _, err := db.Exec(query, msgID, chatID, usrID)
  if err != nil {
    log.Println("[DB]: [ERROR]: cannot add to table")
  }

  log.Println("[DB]: [INFO]: added to table takes ", msgID, chatID, usrID)

  //add to users

  //check if user exists
  if checkIfUserExists(usrID) {
    return
  }

  query = `INSERT INTO users (user_id, user_name) VALUES ($1, $2)`
  _, err = db.Exec(query, usrID, usrName)
  if err != nil {
    log.Println("[DB]: [ERROR]: cannot add to table")
  }

  log.Println("[DB]: [INFO]: added to table users ", usrID, usrName)
}

func checkIfUserExists(usrID int64) bool {
  var dbUserID int64
  query := `SELECT user_id FROM users WHERE user_id = $1`
  err := db.QueryRow(query, usrID).Scan(&dbUserID)

  if err == sql.ErrNoRows {
    return false
  } else if err != nil {
    return true
  } else {
    log.Println("[DB]: [ERROR]: cannot make query when checking if user exists")
  }

  return false
}

func GetChatIDByMsgID(msgID int64) (int64) {
  var chatID int64
  query := `SELECT chat_id FROM takes WHERE msg_id = $1`
  err := db.QueryRow(query, msgID).Scan(&chatID)
  if err != nil {
    log.Println("[DB]: [ERROR]: cannot get chatID from table takes, chatID = ", chatID)
  }

  log.Println("[DB]: [INFO]: got chatID from table takes, chatID = ", chatID)
  return chatID
}

func GetUsrIDByMsgID(msgID int64) (int64) {
  var usrID int64
  query := `SELECT user_id FROM takes WHERE msg_id = $1`
  err := db.QueryRow(query, msgID).Scan(&usrID)
  if err != nil {
    log.Println("[DB]: [ERROR]: cannot get usrID from table users, usrID = ", usrID)
  }

  log.Println("[DB]: [INFO]: got usrID from table users, usrID = ", usrID)
  return usrID
}

func Ban(usrID int64) {
  query := `UPDATE users SET banned = true WHERE user_id = $1`
  _, err := db.Exec(query, usrID)
  if err != nil {
    log.Println("[DB]: [ERROR]: cannot find user in users for ban, banned user = ", usrID)
  }

  log.Println("[DB]: [INFO]: banned user ", usrID)
}

func UnBan(usrID int64) {
  query := `UPDATE users SET banned = false WHERE user_id = $1`
  _, err := db.Exec(query, usrID)
  if err != nil {
    log.Println("[DB]: [ERROR]: find user in users for unban, unbanned user = ", usrID)
  }

  log.Println("[DB]: [INFO]: unbanned user", usrID)
}

func CheckBan(usrID int64) bool {
  var banned bool
  query := `SELECT banned FROM users WHERE user_id = $1`
  err := db.QueryRow(query, usrID).Scan(&banned)
  if err != nil {
    log.Println("[DB]: [ERROR]: cannot check user from table users, usrID = ", usrID)
  }

  log.Println("[DB]: [INFO]: checked user from table users, usrID = ", usrID, ", banned = ", banned)
  return banned
}

func Close() {
  if db != nil {
    db.Close()
  }

  log.Println("[DB]: [INFO]: closed database")
}
