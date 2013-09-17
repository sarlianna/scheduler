package scheduler

import (
  "testing"
  _ "github.com/lib/pq"
  "database/sql"
  "time"
  "strings"
)

//TODO: add database checks. Make sure these are in a test schema/db.
//i don't think skipping tests is the right way to handle it in the case of errors unrelated to the function
//(db errors).  Unsure what to do besides though.

func TestCreateSchedule (t *testing.T) {
  
}

func TestCreateSpan (t *testing.T) {

  gm := GenManager{}
  tdb,err := sql.Open("postgres", "user=postgres dbname=schedulertest sslmode=disable")//foregoing Init() in order to point it at test table

  if err != nil {
    t.Errorf("DB initialization error:%s.\n", err)
  }

  gm.db = tdb
  start, end := time.Now(), time.Now()
  
  res, err := gm.Create( TypeSpan, start, end )
  if err != nil {
    t.Error("Error returned: ", err)
  }

  span := Span(res.(Span))
  if !span.Start.Equal(start) || !span.End.Equal(end) {
    t.Errorf("Span start:%s, end:%s. \nExpected start:%s, end:%s.", span.Start, span.End, start, end)
  }

  row := tdb.QueryRow("SELECT (span_id, start_time, end_time) from spans WHERE ID=" + span.ID)

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
    t.Errorf("DB values - id:%s, start:%s, end:%s. \nExpected - id:%s, start: %s, end:%s.",
                id, outstart, outend, span.ID, start, end)
  }

}
