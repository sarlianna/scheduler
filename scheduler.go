package scheduler

import (
  _ "github.com/lib/pq"
  "strings"
  "time"
  "database/sql"
  "code.google.com/p/go-uuid/uuid"
  "errors"
  "fmt"
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

//Using prepared sql statements stops SQL injection attacks.  Recommendable?

const (
  TypeSchedule = iota
  TypeUser
  TypeGroup
  TypeSpan
)

type Manager interface {
  
  Init() (error)
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
  Start time.Time
  End   time.Time
}

type GenManager struct {
  
  db *sql.DB
}

func (gm GenManager) Init() error { 
  
  db, err := sql.Open( DriverName, ConnectionString )//constants from config.go
  if err != nil {
    return err
  }
  gm.db = db
  return nil
}

func (gm GenManager) Create( otype int, args...interface{} ) (interface{}, error) {
  //this whole method is not very DRY, please rewrite when possible

  //alternatives:
  //instead of case-specific type switches and fills, I could do one general one here.
  //This would mean that the functions the arguments are passed to will also have to 
  //check the arguments to make sure they make sense.
  //For now I do a case specific type-fill and allow the private functions to trust
  //the data passed to them, as only this method should ever call them.
  //additional argument verification should be done here.

  //I think the first alternative might be better....
  
  //unsure that a type assertion is needed inside a type switch; not sure I'm doing this correctly.

  switch otype {

    case TypeSchedule:
      var user string 
      var dates []string
      switch args[0].(type) {
        case string:
          user = args[0].(string)
        default:
          return nil, errors.New("Invalid arguments to be passed to createSchedule.")
      }
      switch args[1].(type) {
        case []string:
          dates = args[1].([]string)
        default:
          return nil, errors.New("Invalid arguments to be passed to createSchedule.")
      }
      return createSchedule( gm.db, user, dates )

    case TypeUser:
      var username, pass, salt string
      var schedules, groups []string

      switch args[0].(type) {
        case string:
          username = args[0].(string)
        default:
          return nil, errors.New("Invalid arguments to be passed to createSchedule.")
      }
      switch args[1].(type) {
        case string:
          pass = args[1].(string)
        default:
          return nil, errors.New("Invalid arguments to be passed to createSchedule.")
      }
      switch args[2].(type) {
        case string:
          salt = args[2].(string)
        default:
          return nil, errors.New("Invalid arguments to be passed to createSchedule.")
      }
      switch args[3].(type) {
        case []string:
          schedules = args[3].([]string)
        default:
          return nil, errors.New("Invalid arguments to be passed to createSchedule.")
      }
      switch args[4].(type) {
        case []string:
          groups = args[4].([]string)
        default:
          return nil, errors.New("Invalid arguments to be passed to createSchedule.")
      }
      return createUser( gm.db, username, pass, salt, schedules, groups )

    case TypeGroup:
      var users, schedules  []string
      switch args[0].(type) {
        case string:
          users = args[0].([]string)
        default:
          return nil, errors.New("Invalid arguments to be passed to createGroup.")

      }
      switch args[1].(type) {
        case string:
          schedules = args[1].([]string)
        default:
          return nil, errors.New("Invalid arguments to be passed to createGroup.")

      }
      return createGroup( gm.db, users, schedules )

    case TypeSpan:
      var start, end time.Time
      switch args[0].(type) {
        case time.Time:
          start = args[0].(time.Time)
        default:
          return nil, errors.New("Invalid arguments to be passed to createSpan.")
      }
      switch args[1].(type) {
        case time.Time:
          end = args[1].(time.Time)
        default:
          return nil, errors.New("Invalid arguments to be passed to createSpan.")
      }
      return createSpan( gm.db, start, end)
    default:
      return nil, errors.New("Invalid type passed to GenManager.Create.")
  }
}

func createSchedule( db *sql.DB, user string, dates []string ) (Schedule, error) {

  if dates != nil {
    //seperate query?
  }
  id := uuid.New()
  var err error
  if strings.EqualFold(user,"") {
    _, err = db.Query("INSERT INTO schedules (schedule_id, user_id) VALUES (" + id + ", " + user + ")" )
  } else {
    _, err = db.Query("INSERT INTO schedules (schedule_id) VALUES (" + id + ")" )
  }

  if err != nil {
    return Schedule{}, err
  }

  userobj := User{ ID: user } //query for this data
  schedule := Schedule{ ID: id, User: &userobj }

  return schedule, nil

}

func createUser( db *sql.DB,  user string, pass string, salt string,
                 schedules []string, groups []string )   (User, error) {

  return User{}, nil
}

func createGroup( db *sql.DB, users []string, schedules []string ) (Group, error) {

  return Group{}, nil
}

func createSpan( db *sql.DB,  start time.Time, end time.Time ) (Span, error) {
  
  //check that start and end are proper timestamps?
   
  id := uuid.New();
  stmt, sterr := db.Prepare("INSERT INTO spans (span_id, start_time, end_time) VALUES (?, ?, ?)") //is this really good practice?
  if sterr != nil {
    fmt.Println("Prepare statement failed.")
    return Span{}, sterr
  }
  fmt.Println("Values: " + id + start.String() + end.String() )
  _, err := stmt.Exec(id, start, end)
  //_, err := db.Exec("INSERT INTO spans (span_id, start_time, end_time) VALUES ('?', ?, ?)", id  ,start, end )

  if err != nil {
    return Span{}, err 
  }

  return Span{ ID: id, Start: start, End: end }, nil
}

