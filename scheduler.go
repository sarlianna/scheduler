package scheduler

import (
  _ "github.com/lib/pq"
  "strings"
  "time"
  "database/sql"
  "code.google.com/p/go-uuid/uuid"
  "errors"
  "fmt"
  "strconv"
)

/* better planning:
  
  think of it in terms of REST.  it should be designed with this in mind.

  Objects don't need references down at all ( unsure but go with it for now)
  To get children you would specify an additional parameter, i.e.
  api/user_id              will get a user with user_id, while
  api/user_id/schedules    will get the same user's schedules.

  when deleting, delete child objects with no other references also.
  (Seems natural to delete a user's schedules when deleting a user, but not schedules assigned to a group when deleting that group, etc.)

  may be changed in the future.

/* old planning:
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
  //Dates []*Span 
  //Group *Group
}

type User struct {
  
  ID        string 
  Username  string 
  Password  string 
  Salt      string 
  //Schedules []*Schedule
}

type Group struct {

  ID        string
  //Schedules []*Schedule
  Title     string
  Desc      string
}

type Span struct {

  ID       string
  Start    time.Time
  End      time.Time
  Schedule *Schedule
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
      var user, group string 
      switch args[0].(type) {
        case string:
          user = args[0].(string)
        default:
          return nil, errors.New("Invalid arguments to be passed to createSchedule.")
      }
      switch args[1].(type) {
        case string:
          group = args[1].(string)
        default:
          return nil, errors.New("Invalid arguments to be passed to createSchedule.")
      }
      return createSchedule( gm.db, user, group )

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
      var schedules  []string
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

func (gm GenManager) Read( otype int, size int, offset int, id string) (interface{}, error) {

  if size < 0 || offset < 0 {
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

func (gm GenManager) Update( otype int, args...interface{}) (interface{}, error) {

  return nil, nil

}

func (gm GenManager) Delete( otype int, id string) (error) {

  var err error
  switch otype {
    case TypeSchedule:
      _, err = gm.db.Query("DELETE FROM schedules WHERE id=?", id)
    case TypeUser:
      _, err = gm.db.Query("DELETE FROM users WHERE id=?", id)
    case TypeGroup:
      _, err = gm.db.Query("DELETE FROM groups WHERE id=?", id)
    case TypeSpan:
      _, err = gm.db.Query("DELETE FROM spans WHERE id=?", id)
    default:
      return errors.New("Invalid type passed to GenManager.Delete.")
  }
  if err != nil {
    return err
  }

  return nil
}

func createSchedule( db *sql.DB, user string, group string) (Schedule, error) {

  id := uuid.New()
  var query string
  if !strings.EqualFold(user,"") && !strings.EqualFold(group, "") {
    query = "INSERT INTO schedules (schedule_id, user_id, group_id) VALUES (" + id + ", " + user + ", " + group +")"
  } else if !strings.EqualFold(user, "") && strings.EqualFold(group, ""){
    query = "INSERT INTO schedules (schedule_id) VALUES (" + id + ")" 
  }

  _, err := db.Query(query)
  if err != nil {
    return Schedule{}, err
  }

  userobj := User{ ID: user } //query for this data
  schedule := Schedule{ ID: id, User: &userobj }

  return schedule, nil

}

func createUser( db *sql.DB,  user string, pass string, salt string,
                 schedules []string, groups []string )   (User, error) {

  id := uuid.New()

  _, err := db.Exec("INSERT INTO users (username, password, salt) VALUES ( ?, ?, ?)", user, pass, salt)
  if err != nil {
    return User{}, err
  }
  
  return User{ID: id, Username: user, Password: pass, Salt: salt}, nil 
}

func createGroup( db *sql.DB, title string, description string, schedules []string ) (Group, error) {

  id := uuid.New()

  _, err := db.Exec("INSERT INTO groups (group_id, title, description) VALUES ( ?, ?, ?)", id, title, description )
  if err != nil {
    return Group{}, err
  }

  return Group{ID: id, Title:title, Desc: description}, nil
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
  _, err := db.Exec("INSERT INTO spans (span_id , start_time , end_time) VALUES (? , ? , ?)", id, start, end )

  if err != nil {
    return Span{}, err 
  }

  return Span{ ID: id, Start: start, End: end }, nil
}

func readSchedule(db *sql.DB, num int, offset int, id string ) ([]Schedule, error) {

  if strings.EqualFold(id, "") {
    rows, err := db.Query("SELECT sch.schedule_id, sch.user_id, user.username, user.password, user.salt " + 
                            "FROM schedules sch, users user WHERE user.user_id = sch.user_id " + 
                            "OFFSET " + strconv.Itoa(offset) + " LIMIT " + strconv.Itoa(num) )
    if err != nil {
      return nil, err
    }


    ret := make([]Schedule, 1)
    for i := 0; rows.Next(); i++ {
      scanerr := rows.Scan( &ret[i].ID, &ret[i].User.ID, &ret[i].User.Username, &ret[i].User.Password, &ret[i].User.Salt )
      if scanerr != nil {
        return nil, scanerr
      }
    }
    
    return ret, nil

  } else {

    //get by id
    result := Schedule{}
    err := db.QueryRow("SELECT sche.schedule_id, sch.user_id, user.username, user.password, user.salt " + 
                        "FROM spans WHERE sch.schedule_id=" + id + 
                        " AND user.user_id = sch.user_id").Scan(&result.ID, &result.User.ID, &result.User.Username, 
                          &result.User.Password, &result.User.Salt )
    if err != nil {
      return nil, err
    }

    ret := []Schedule{ result }
    return ret, nil
  }
  

  return nil, nil
}

func readUser(db *sql.DB, num int, offset int, id string ) ([]User, error) {

  if strings.EqualFold(id, "") {
    rows, err := db.Query("SELECT user_id, username, password, salt, FROM users user" +
                          "OFFSET " + strconv.Itoa(offset) + " LIMIT " + strconv.Itoa(num))

    if err != nil {
      return nil, err
    }

    ret := make( []User, 1 )
    for i:=0; rows.Next(); i++ {
      scanerr := rows.Scan( &ret[i].ID, &ret[i].Username, &ret[i].Password, &ret[i].Salt )
      if scanerr != nil {
        return nil, scanerr
      }
    }

    return ret, nil
  } else {

    result := User{}
    err := db.QueryRow("SELECT user.user_id, user.username, user.password, user.salt, " +
            "FROM users user" +
            "OFFSET " + strconv.Itoa(offset) + " LIMIT " + strconv.Itoa(num)).Scan(&result.ID,
            &result.Password, &result.Salt)
    if err != nil {
      
      return nil, err
    }

    ret := []User { result }
    return ret, nil

  }

  return nil, nil
}

func readGroup(db *sql.DB, num int, offset int, id string ) ([]Group, error) {

  if strings.EqualFold(id, "") {
    rows, err := db.Query("SELECT group_id, title, description, FROM groups " +
                           "OFFSET " + strconv.Itoa(offset) + " LIMIT " + strconv.Itoa(num))
    if err != nil {
      return nil, err
    }

    ret := make( []Group, 1 )
    for i:=0; rows.Next(); i++ {
      scanerr := rows.Scan( &ret[i].ID, &ret[i].Title, &ret[i].Desc )
      if scanerr != nil {
        return nil, scanerr
      }
    }

    return ret, nil
  } else {

    result := Group {}
    err := db.QueryRow("SELECT group_id, title, description, FROM groups " +
                       "OFFSET " + strconv.Itoa(offset) + " LIMIT " + 
                       strconv.Itoa(num)).Scan(&result.ID, &result.Title, &result.Desc)
    if err != nil {
      return nil, err
    }

    ret := []Group{ result }
    return ret, nil
  }
 

  return nil, nil
}

func readSpan(db *sql.DB, num int, offset int, id string ) ([]Span, error) {

  if strings.EqualFold(id, "") {
    query := "SELECT span_id, start_time, end_time FROM spans OFFSET " + strconv.Itoa(offset) + " LIMIT " + strconv.Itoa(num)
    rows, err := db.Query(query)

    if err != nil {
      return nil, err
    }

    ret := make([]Span, 1)
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
    err := db.QueryRow("SELECT span_id, start_time, end_time FROM spans WHERE span_id=" + id).Scan(&result.ID, &result.Start, &result.End)
    if err != nil {
      return nil, err
    }

    ret := []Span{ result }
    return ret, nil
  }
}
