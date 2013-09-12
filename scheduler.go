package scheduler

import (
  _ "github.com/lib/pq"
  "string"
  "time"
  "database/sql"
  "code.google.com/p/go-uuid/uuid"
)

/* planning:
  so, internal library that returns objects for the app, handles db interaction.
  db info is in config.go
  one handler, you pass it a type.  Still unsure if this is the best architechture.
  handler returns a call to an appropriate CRUD method.
  doing this test-driven style.
  
  for now, an object will just return with id to other objects, I guess.

  create methods are looking for arrays of strings with other object ids if you already have
  which items you want in collections
  
  types: groups, users, schedules, spans.
*/

const (
  TypeSchedule = iota
  TypeUser     = iota
  TypeGroups   = iota
  TypeSpans    = iota
)

type Manager interface {
  
  Init() (nil, error)
  Create( otype int, args...interface{} ) (interface{}, error) //interface{} being used as a general type for data.  This may be wrong or bad style.
  Read(   otype int, args...interface{} ) (interface{}, error)
  Update( otype int, args...interface{} ) error
  Delete( otype int, ID string ) error
}

type Schedule struct {

  ID    string
  User  *User
  Dates *[]Span 
}

type User struct {
  
  ID        string 
  Username  string 
  Password  string 
  Salt      string 
  Schedules *[]Schedule
}

type Group struct {

  ID        string
  Schedules *[]Schedule
  Title     string
  Desc      string
}

type Span struct {

  ID    string
  Start string
  End   string
}

type GenManager struct {
  
  db *DB
}

func (GenManager gm) Init() nil {
  
  db, err := sql.Open( DriverName, ConnectionString )//constants from config.go
  if err != nil {
    return nil, err
  }
  gm.db = db
  return nil, nil
}

func CreateSchedule( db DB, user string, dates []string ) Schedule {

  if dates != nil {
    //seperate query?
  }
  id := uuid.New()
  if user != nil {
    res, err := db.Query("INSERT INTO schedules (ID, user) VALUES (" + id + ", " + user + ")" )
  } else {
    res, err := db.Query("INSERT INTO schedules (ID) VALUES (" + id + ")" )
  }

  if err != nil {
    return nil, err
  }

  userobj = User{ ID: user } //query for this data
  schedule = Schedule{ ID: uuid, User: &userobj }

  return schedule, nil

}

func CreateUser( db DB,  user string, pass string, salt string,
                 schedules []string, groups []string )   User {

}

func CreateGroup( db DB, users []string, schedules []string ) Group {

}

func CreateSpan( db DB,  start string, end string ) Span {
  
  if start == nil && end == nil {
    return nil, error.New("Span creation requires start and end date")
  }
  //check that start and end are proper timestamps?
   
  id := uuid.New();
  res, err := db.Query("INSERT INTO spans (ID, start, end) VALUES (" + id + ", " + start + ", " + end + ")" )

  if err != nil {
    return nil, err 
  }

  span = Span{ ID: id, Start: start, End: end }
  return span, nil

}


