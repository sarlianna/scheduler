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
  Read(   otype int, size int, offset int, id string) (interface{}, error)//id is optional. 
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
      var title, description string
      var users, schedules  []string
      switch args[0].(type) {
        case string:
          title = args[0].(string)
        default:
          return nil, errors.New("Invalid arguments to be passed to createGroup.")

      }
      switch args[1].(type) {
        case []string:
          description = args[1].(string)
        default:
          return nil, errors.New("Invalid arguments to be passed to createGroup.")

      }
      switch args[2].(type) {
        case []string:
          schedules = args[1].([]string)
        default:
          return nil, errors.New("Invalid arguments to be passed to createGroup.")

      }
      return createGroup( gm.db, title, description, schedules )

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

func (GenManager gm) Read( otype int, size int, offset int, id string) (interface{}, error) {

  if num < 0 || offset < 0 {
    return nil, errors.New("Invalid size or offset value.")
  }
  switch otype {
    case TypeSchedule:
      return readSchedule( gm.db, size, offset, id)
    case TypeUser:
      return readUser( gm.db, size, offset, id)
    case TypeGroup:
      return readGroup( gm.db, size, offset, id)
    case TypeSpan:
      return readSpan( gm.db, size, offset, id)
    default:
      return nil, errors.New("Invalid type passed to GenManager.Read.")
  }
}

func (GenManager gm) Update( otype int, args...interface{}) (interface{}, error) {


}

func (GenManager gm) Delete( otype int, id string) (error) {

  switch otype {
    case TypeSchedule:
      _, err := gm.db.Query("DELETE FROM schedules WHERE id=?", id)
    case TypeUser:
      _, err := gm.db.Query("DELETE FROM users WHERE id=?", id)
    case TypeGroup:
      _, err := gm.db.Query("DELETE FROM groups WHERE id=?", id)
    case TypeSpan:
      _, err := gm.db.Query("DELETE FROM spans WHERE id=?", id)
    default:
      return errors.New("Invalid type passed to GenManager.Delete.")
  }
  if err != nil {
    return err
  }

  return nil
}

func createSchedule( db *sql.DB, user string, dates []string ) (Schedule, error) {

  if dates != nil {
    //seperate query?
  }
  id := uuid.New()
  if strings.EqualFold(user,"") {
    _, err := db.Query("INSERT INTO schedules (schedule_id, user_id) VALUES (" + id + ", " + user + ")" )
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

  if schedules != nil {

  }
  if groups != nil {

  }

  id := uuid.New()

  _, err := db.Exec("INSERT INTO users (username, password, salt) VALUES ( ?, ?, ?)", user, pass, salt)
  if err != nil {
    return User{}, err
  }
  
  return User{ID: id, Username: user, Password: pass, Salt: salt}, nil //TODO: add Schedules
}

func createGroup( db *sql.DB, title string, description string, schedules []string ) (Group, error) {

  if schedules != nil {

  }

  id := uuid.New()

  _, err := db.Exec("INSERT INTO groups (group_id, title, description) VALUES ( ?, ?, ?)", id, title, description )
  if err != nil {
    return Group{}, err
  }

  return Group{ID: id, Title:title, Desc: description}, nil //TODO: add Schedules
}

func createSpan( db *sql.DB,  start time.Time, end time.Time ) (Span, error) {
  
  //check that start and end are proper timestamps?
   
  id := uuid.New();
  /*
  stmt, sterr := db.Prepare("INSERT INTO spans ( span_id , start_time , end_time ) VALUES ( ? , ? , ? )") //is this really good practice? (preparing just before
  if sterr != nil {
    fmt.Println("Prepare statement failed.")
    return Span{}, sterr
  }
  */
  fmt.Println("Values: " + id + start.String() + end.String() )
  //_, err := stmt.Exec(id, start, end)
  _, err := db.Exec("INSERT INTO spans (span_id , start_time , end_time) VALUES ('?' , ? , ?)", id, start, end )

  if err != nil {
    return Span{}, err 
  }

  return Span{ ID: id, Start: start, End: end }, nil
}

func readSchedule(db *sql.DB, num int, offset int, id string ) ([]Schedule, error) {

  if strings.EqualFold(id, "") {
    rows, err := db.Query("SELECT a.schedule_id, a.user_id, b.span_id, b.start_time, b.end_time FROM schedules a, spans b OFFSET "
                            + offset + " LIMIT " + num )

    ret := make([]Schedule)
    for i := 0; rows.Next(); i++ {
      //one loop isn't enough - we need to fill schedules, and for each schedule, fill every span.
    }


  } else {

  }
}

func readUser(db *sql.DB, num int, offset int, id string ) ([]User, error) {

  if strings.EqualFold(id, "") {
    query :="SELECT user.user_id, user.username, user.password, user.salt, " +
            "schedule.schedule_id, schedule.user_id, span.span_id, " +
            "span.start_time, span.end_time " +
            "FROM users user, schedules schedule, spans span " +
            "OFFSET " + offset + " LIMIT " + num
    rows, err := db.Query(query)

    ret := make( []User )
    for i:=0; rows.Next(); i++ {
      //
    }
  } else {

  }

}

func readGroup(db *sql.DB, num int, offset int, id string ) ([]Group, error) {

  if strings.EqualFold(id, "") {
    query := "SELECT group.group_id, group.title, group.description, " +
             "schedule.schedule_id, schedule.user_id, span.span_id, " +
             "span.start_time, span.end_time " +
             "FROM groups group, schedules schedule, spans span " +
             "OFFSET " + offset + " LIMIT " + num
    rows, err := db.Query(query)

    ret := make( []Group )
    for i:=0; rows.Next(); i++ {
      //
    }

  } else {

  }
 

}

func readSpan(db *sql.DB, num int, offset int, id string ) ([]Span, error) {

  if strings.EqualFold(id, "") {
    query := "SELECT span_id, start_time, end_time FROM spans OFFSET " + offset + " LIMIT " + num
    rows, err := db.Query(query)

    ret := make([]Span)
    for i:=0; rows.Next(); i++ {
      cur := Span{}
      scanerr := rows.Scan(&cur.ID, &cur.Start, &cur.End)
      if scanerr != nil {
        return nil, scanerr
      }
      ret = append(ret, cur)
    }

    return ret, nil
  
  } else {
    //lookup by id
    result := Span{}
    err := db.Query("SELECT span_id, start_time, end_time FROM spans WHERE span_id=" + id).Scan(&result.ID, &result.Start, &result.End)
    if err != nil {
      return nil, err
    }

    ret := []Span{ result }
    return ret, nil
  }
}
