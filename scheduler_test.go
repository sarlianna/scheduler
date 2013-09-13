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
  tdb,err := sql.Open("postgres", "user=postgres dbname=schedulertest")//foregoing Init() in order to point it at test table

  if err != nil {
    t.Log("DB initialization error:" + err + ".\nSkipping test.")
    t.Skip()
  }

  gm.db = tdb
  const start, end = time.Now().String, time.Now().String()
  
  res, err := gm.Create( TypeSpan, start, end )
  if err != nil {
    t.Error("Error returned: ", err)
  }

  span := res.(Span)
  if !strings.EqualFold(span.start,start) || !strings.EqualFold(span.end, end) {
    t.Errorf("Span start:%s, end:%s. \nExpected start:%s, end:%s.", span.start, span.end, start, end)
  }

  row, queryErr := tdb.QueryRow("SELECT (ID, start, end) from spans WHERE ID=" + span.ID)

  if queryErr != nil {
    t.Log("DB Query error:" + queryErr + ".\nSkipping test.")
    t.Skip()
  }

  var id,outstart,outend string
  scanErr := row.Scan(&id, &outstart, &outend)
  
  if scanErr != nil {
    if scanErr == sql.ErrNoRows {
      t.Error("No rows with ID given by CreateSpan were found in the database.")
    }
    t.Log("Row scan error:" + scanErr + ".\nSkipping test.")
    t.Skip()
  }

  if !strings.EqualFold(start, outstort) || 
     !strings.EqualFold(end,outend) || 
     !strings.EqualFold(id,span.ID) {
    t.Errorf("DB values - id:%s, start:%s, end:%s. \nExpected - id:%s, start: %s, end:%s.",
                id, outstart, outend, span.ID, start, end)
  }

}
