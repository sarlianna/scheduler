package scheduler

import (
  "time"
  "database/sql"
  _ "github.com/lib/pq"
  "strconv"
)

//objects need to have a go library with exported CRUD methods.
//json api handled in another package/app

//things to ask colton about:
//pagination structure

//notes from talking with colton
//general manager if possible, but make sure whatever it is is extendable with actions, for the case of password hashing, etc.
//look through go stuff on 7th prime
//for read method, paginate it but don't filter it.  READ just returns all.  include a size and offset so people reading don't need
//  to access all.  

//more notes from colton (wed 21)
//extendsible actions - not structured this way in go, you need do nothing.
//sql storage with collections: many to one backwards lookup structure - the object with the collection has no indication of any of
//  the objects in the collection, but all the objects in the collection reference the main object.  You unite them with a join query.
//database is fine
//write tests!! use a test table with the database to check methods that use it.
//write migrations
//lookup go packages and new technologies, generally. (not reaally related to this project)
//

//TODO: 
// -fix update to take selective data
// -fix read.  Use interface{} if possible; multiple type assignment in that function may be messy.
// -declare common errors as constants so they can be easily checked against?

//type 'enum' for general manager
const (

  ScheduleType = iota
  UserType     = iota
  GroupType    = iota
  SpanType     = iota
)

type Manager interface {
  
  Init() error //initializes db connection

  Create( sType int ) (interface{}, error)
  Read( sType int, size int, offset int) ([]interface{}, error)//correct
  Update( sType int, id string, data interface{}) error
  Del( stype int, id string ) error
}

//base data types
type Schedule struct {

  ID string
  Dates []Span//many to one
  User *User 
}

type User struct {

  ID string
  Schedules *[]Schedule

  Username string
  Password string
  Salt     string
}

type Group struct {

  ID string
  Schedules []Schedule
  Title string
  Description string
}

type Span struct {

  ID string
  Start       Time
  End         Time
}


//manager and interface implementations

type SManager struct {

  db *DB
}

func (sm SManager) Init() error {
  
  db, err := sql.Open(DriverName, ConnectionString) //constants from config.go

  if err != nil {
    return err
  }

  sm.db = db
  return nil
}

func (sm SManager) Create( sType int, data interface{}  ) (interface{}, error) {
  //golang sql.Result interface gives you last inserted row's id.
  //should also only insert passed fields, like update.

  var query string
  switch sType {
    case ScheduleType:
      if data.dates != nil {
        //seperate query to update these?
      }
      if data.user != nil {
        query = "INSERT INTO schedules (user) VALUES (" + data.user + ")"
      } else {
        query = ""
      }
    case UserType:
      if data.schedules != nil {
        //do a seperate query to update each schedule
      }

      query += "INSERT INTO users ( "
      endquery := "("
      if data.username != nil {
        query += "username"
        endquery += data.username
      }
      if data.password != nil {
        if ! strings.HasSuffix(query, " ") {
          query += ", "
          endquery += ", "
        }
        query += "password"
        endquery += data.password 
      }
      if data.salt != nil {
        if ! strings.HasSuffix(query, " ") {
          query += ", "
          endquery += ", "
        }
        query += "salt=" + data.salt
        endquery += data.salt
      }
      endquery += ")"
      query += ") VALUES " + endquery

      if data.username != nil || data.password != nil || data.salt != nil {
        _, err := sm.db.Query(query)
      }

      //query = "INSERT INTO users (schedules, username, password, salt) VALUES (" + data.schedules + ", " + data.username + ", " +
      //          data.password + ", " + data.salt + ")"
    case GroupType:

      query += "INSERT INTO groups "
      endquery := "("
      if data.title != nil {
        query += "title" 
        endquery += data.title
      }
      if data.description != nil {
        if ! strings.HasSuffix(query, " ") {
          query += ", "
          endquery += ", "
        }
        query += "description" 
        endquery += data.description
      }
      endquery += ")"
      query += ") VALUES " + endquery

      if data.title != nil || data.description != nil {
        _, err := sm.db.Query(query)
      }

      //query = "INSERT INTO groups (schedules, title, descripton) VALUES (" + data.schedules + ", " + 
      //          data.title + ", " + data.descripton +")"
    case SpanType:

     query += "INSERT INTO spans ( "
     endquery := "("
     if data.start != nil {
        query += "start"
        endquery += data.start
      }
     if data.end != nil {
       if ! strings.HasSuffix(query, " ") {
         query += ", "
         endquery += ", "
       }
       query += "end"
       endquery += data.end
     }
     endquery += ")"
     query += ") VALUES " + endquery

     if data.start != nil || data.end != nil {
       _, err := sm.db.Query(query)
     }

      //query = "INSERT INTO spans (start, end) VALUES (" + data.start + ", " + data.end + ")"
    default:
      return nil, errors.new("invalid type")
  }

  _, err := sm.db.Query(query)
  if err != nil {
    return nil, err
  }

  //insert id into data
  return data
}

func (sm SManager) Read( sType int, size int, offset int) ([]interface{}, error) {
  //input check
  if size < 0 { //max bounds check ever necessary?
    return nil, errors.New("invalid size value")
  }
  if offset < 0 {
    return nil, errors.New("invalid offset value")
  }

  var query string
  s := strconv.Itoa(size)
  o := strconv.Itoa(offset)
  switch sType {
    case ScheduleType:
      query = "SELECT a.ID, a.user, b.ID, b.start, b.end FROM schedules a, spans b OFFSET " + o + " LIMIT " + s
      rows, err := sm.db.Query(query)
      var ret []Schedule

      for i:= 0; rows.Next(); i++ {
        errc := rows.Scan(&ret[i].ID, &ret[i].user, &ret[i].dates[i].ID, &ret[i].dates[i].start, &ret[i].dates[i].end)
        if errc != nil {
          return nil, errc
        }
      }
      return ret, nil

    case UserType:
      query = "SELECT a.ID, a.username, a.password, a.salt, " +
              "b.ID, b.user, c.ID, c.start, c.end " +
              "FROM users a, schedules b, spans c " + 
              "OFFSET " + o + " LIMIT " + s
      rows, err := sm.db.Query(query)
      var ret []User

      for i:= 0; rows.Next(); i++ {
        errc := rows.Scan(&ret[i].ID, &ret[i].username, &ret[i].password, &ret[i].salt, 
                          &ret[i].schedules[i].ID, &ret[i].schedules[i].user, &ret[i].schedules[i].dates[i].ID,
                          &ret[i].schedules[i].dates[i].start, &ret[i].schedules[i].dates[i].end)
        if errc != nil {
          return nil, errc
        }
      }
      return ret, nil
      
    case GroupType:
      query = "SELECT a.ID, a.title, a.description " + 
              "b.ID, b.user, c.ID, c.start, c.end " +
              "FROM groups a, schedules b, spans c " + 
              "OFFSET " + o + " LIMIT " + s
      rows, err := sm.db.Query(query)
      var ret []Group

      for i:= 0; rows.Next(); i++ {
        errc := rows.Scan(&ret[i].ID, &ret[i].title, &ret[i].description, &ret[i].schedules[i].ID,
                          &ret[i].schedules[i].user, &ret[i].schedules[i].dates[i].ID,
                          &ret[i].schedules[i].dates[i].start, &ret[i].schedules[i].dates[i].end) 
        if errc != nil {
          return nil, errc
        }
      }
      return ret, nil

    case SpanType:
      query = "SELECT ID, start, end FROM spans OFFSET " + o + " LIMIT " + s
      rows, err := sm.db.Query(query)
      var ret []Span

      for i:= 0; rows.Next(); i++ {
        errc := rows.Scan(&ret[i].ID, &ret[i].start, &ret[i].end)
        if errc != nil {
          return nil, errc
        }
      }
      return ret, nil

    default:
      return nil,errors.new("invalid type")
  }

  return nil, errors.new("execution reached an unintended portion")
}

func (sm SManager) Update(sType int, id string, data interface{}) error {
  //TODO: finish many-one updates

  query := "UPDATE "
  switch sType {
    case ScheduleType:
      if data.dates != nil {
        //if dates is set, do a seperate query to update each of the dates.
      }
      if data.user != nil {
        _, err := sm.db.Query("UPDATE schedules SET user=? WHERE id=?", data.user, id)
      }

    case UserType: 
      if data.schedules != nil {
        //do a seperate query to update each schedule
      }

      query += "users SET "
      if data.username != nil {
        query += "username=" + data.username
      }
      if data.password != nil {
        if ! strings.HasSuffix(query, " ") {
          query += ","
        }
        query += " password=" + data.password
      }
      if data.salt != nil {
        if ! strings.HasSuffix(query, " ") {
          query += ","
        }
        query += " salt=" + data.salt
      }
      query += " WHERE id=" + id

      if data.username != nil || data.password != nil || data.salt != nil {
        _, err := sm.db.Query(query)
      }

    case GroupType:
      if data.schedules != nil {
        //do queries to each schedule and update
      }

      query += "groups SET "
      if data.title != nil {
        query += "title=" + data.title
      }
      if data.description != nil {
        if ! strings.HasSuffix(query, " ") {
          query += ", "
        }
        query += "description=" + data.description
      }
      query += "WHERE id=" + id

      if data.title != nil || data.description != nil {
        _, err := sm.db.Query(query)
      }

    case SpanType:
      query += "span SET "
      if data.start != nil {
        query += "start=" + data.start
      }
      if data.end != nil {
        if ! strings.HasSuffix(query, " ") {
          query += ", "
        }
        query += "end=" + data.end
      }
      query += "WHERE id=" + id

      if data.start != nil || data.end != nil {
        _, err := sm.db.Query(query)
      }

    default:
      return errors.new("invalid type")
  }
  if err != nil {
    return err
  }

  return nil
}

func (sm SManager) Del( sType int, id string ) error {
  
  //the switch statement could be taken out if the types were enumerated to string equivalents. (not really an enum, I guess)
  var err error;
  switch sType {
    case ScheduleType:
      _, err = sm.db.Query("DELETE FROM schedules WHERE id=?", id)
    case UserType:
      _, err = sm.db.Query("DELETE FROM users WHERE id=?", id)
    case GroupType:
      _, err = sm.db.Query("DELETE FROM groups WHERE id=?", id)
    case SpanType:
      _, err = sm.db.Query("DELETE FROM spans WHERE id=?", id)
    default:
      return errors.New("invalid type")
  }
  if err != nil {
    return err
  }

  return nil
}

