package scheduler 

import (
  "fmt"
  "database/sql"
  _ "github.com/lib/pq"
)

//to make this executable, the best way is probably to make this package main, include scheduler, and use scheduler.DriverName etc.
//but I do want this distributed with the package...

func main() {
  
  db, err := sql.Open( DriverName, ConnectionString )
  if err != nil {
    fmt.Println( "Error connecting to db: ", err )
  }


  _, err = db.Exec( "CREATE TABLE users ( user_id UUID PRIMARY KEY, " +
                                      "username varchar(50)," +
                                      "password varchar(80)," +
                                      "salt varchar(80) )" )
  if err != nil {
    fmt.Println( "Error creating tables: ", err)
  }

  _, err = db.Exec( "CREATE TABLE groups ( group_id UUID PRIMARY KEY, " +
                                      "title varchar(120)," +
                                      "description varchar(200) )" )
  if err != nil {
    fmt.Println( "Error creating tables: ", err)
  }

  _, err = db.Exec( "CREATE TABLE schedules ( schedule_id UUID PRIMARY KEY, " +
                                      "user_id UUID REFERENCES users," +
                                      "group_id UUID REFERENCES groups )")
  if err != nil {
    fmt.Println( "Error creating tables: ", err)
  }

  _, err = db.Exec( "CREATE TABLE spans ( span_id UUID PRIMARY KEY, " +
                                      "start_time timestamp," +
                                      "end_time timestamp," +
                                      "schedule_id UUID REFERENCES schedules )")
  if err != nil {
    fmt.Println( "Error creating tables: ", err)
  } 

}
