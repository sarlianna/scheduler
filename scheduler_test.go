package scheduler

import (
  "testing"
  _ "github.com/lib/pq"
  "database/sql"
  "time"
  "strings"
)


func TestCreateUser (t *testing.T) {

  gm := GenManager{}
  tdb,err := sql.Open(DriverName, TestConnectionString)

  if err != nil {
    t.Errorf("DB initialization error:%s.\n", err)
  }

  gm.db = tdb
  test_user := "Bob"
  test_pass := "hello"
  test_salt := "haha"

  res, err := gm.Create( TypeUser, test_user, test_pass, test_salt )
  if err != nil {
    t.Error("Error returned: ", err)
  }
  
  user := res.(User) 
  if !strings.EqualFold(test_user, user.Username) || !strings.EqualFold(test_pass, user.Password) || 
    !strings.EqualFold(test_salt, user.Salt) {
      
      t.Errorf("Method value mismatch.\nExpected: %s, %s, %s.\nActual: %s, %s, %s.", test_user, test_pass, 
                test_salt, user.Username, user.Password, user.Salt)
   }

  var id, username, password, salt string
  row := tdb.QueryRow("SELECT (user_id, username, password, salt) FROM users WHERE user_id=" + user.ID)
  
  scanErr := row.Scan(&id, &username, &password, &salt)

  if scanErr != nil {
    if scanErr == sql.ErrNoRows {
      t.Error("No rows with user_id given by CreateUser were found in the database.")
    }
    t.Errorf("Row scan error:%s.\n", scanErr)
  }

  if !strings.EqualFold(id, user.ID) ||
      !strings.EqualFold(username,test_user) ||
      !strings.EqualFold(password,test_pass) ||
      !strings.EqualFold(salt,test_salt) {

       t.Errorf("DB value mismatch.\nExpected: %s, %s, %s, %s.\nActual: %s, %s, %s, %s.", user.ID, test_user, test_pass, 
                test_salt, id, username, password, salt)
  }

}

func TestCreateGroup (t *testing.T) {

  gm := GenManager{}
  tdb,err := sql.Open(DriverName, TestConnectionString)

  if err != nil {
    t.Errorf("DB initialization error:%s.\n", err)
  }

  gm.db = tdb
  test_title := "IMPORTANT GROUP"
  test_desc  := "HELLO there"

  res, err := gm.Create( TypeGroup, test_title, test_desc )
  if err != nil {
    
    t.Error("Error returned: ", err)
  }

  group := res.(Group)
  if !strings.EqualFold(test_title, group.Title) || !strings.EqualFold(test_desc, group.Desc) {

    t.Errorf("Method value mismatch. \nExpected: %s, %s.\nActual: %s, %s.", test_title, test_desc, group.Title, group.Desc)
  }

  var id, title, desc string
  row := tdb.QueryRow("SELECT (group_id, title, description) FROM groups WHERE group_id=" + group.ID)

  scanErr := row.Scan(&id, &title, &desc)
  if scanErr != nil {

    t.Error("Row scan error:%s.\n", scanErr)
  }

  if !strings.EqualFold(title, test_title) ||
      !strings.EqualFold(desc, test_desc) ||
      !strings.EqualFold(id, group.ID) {

        t.Errorf("DB value mismatch.\nExpected: %s, %s, %s.\nActual: %s, %s, %s.", group.ID, test_title, test_desc, id, title, desc)
     }


}

func TestCreateSchedule (t *testing.T) {
  
  gm := GenManager{}
  tdb,err := sql.Open(DriverName, TestConnectionString )//foregoing Init() in order to point it at test table

 
  if err != nil {
    t.Errorf("DB initialization error:%s.\n", err)
  }

  gm.db = tdb

}

func TestCreateSpan (t *testing.T) {

  gm := GenManager{}
  tdb,err := sql.Open( DriverName, TestConnectionString )

  if err != nil {
    t.Errorf("DB initialization error:%s.\n", err)
  }

  gm.db = tdb
  start, end := time.Now(), time.Now()
  
  res, err := gm.Create( TypeSpan, start, end )
  if err != nil {
    t.Error("Error returned: ", err)
  }

  span := res.(Span)
  if !span.Start.Equal(start) || !span.End.Equal(end) {
    t.Errorf("Method value mismatch.\nExpected:%s, %s. \nActual start:%s, end:%s.", start, end, span.Start, span.End,)
  }

  row := tdb.QueryRow("SELECT (span_id, start_time, end_time) from spans WHERE span_id=" + span.ID)

  var id string
  var outstart, outend time.Time
  scanErr := row.Scan(&id, &outstart, &outend)
  
  if scanErr != nil {
    if scanErr == sql.ErrNoRows {
      t.Error("No rows with span_id given by CreateSpan were found in the database.")
    }
    t.Errorf("Row scan error:%s.\n", scanErr)
  }

  if !start.Equal(outstart) || 
     !end.Equal(outend) || 
     !strings.EqualFold(id,span.ID) {
    t.Errorf("DB value mismatch.\nExpected: %s, %s, %s. \nActual: %s, %s, %s.",
                span.ID, start, end, id, outstart, outend )
  }

}
