package scheduler 

import (
  "fmt"
  "database/sql"
  _ "github.com/lib/pq"
)

func main() {
  
  db, err := sql.Open( DriverName, ConnectionString )
  if err != nil {
    fmt.Println( "Error connecting to db: ", err )
  }


  _, err = db.Exec( "CREATE TABLE users ( ID UUID PRIMARY KEY, " +
                                      "username varchar," +
                                      "password varchar," +
                                      "salt varchar )" )
  if err != nil {
    fmt.Println( "Error creating tables: ", err)
  }

  _, err = db.Exec( "CREATE TABLE groups ( ID UUID PRIMARY KEY, " +
                                      "title varchar," +
                                      "description varchar )" )
  if err != nil {
    fmt.Println( "Error creating tables: ", err)
  }

  _, err = db.Exec( "CREATE TABLE schedules ( ID UUID PRIMARY KEY, " +
                                      "user UUID REFERENCES users(ID)," +
                                      "group UUID REFERENCES groups(ID) )")
  if err != nil {
    fmt.Println( "Error creating tables: ", err)
  }

  _, err = db.Exec( "CREATE TABLE spans ( ID UUID PRIMARY KEY, " +
                                      "start timestamp," +
                                      "end timestamp," +
                                      "schedule UUID REFERENCES schedules(ID) )")
  if err != nil {
    fmt.Println( "Error creating tables: ", err)
  } 

}
