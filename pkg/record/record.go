package record

import (
  "errors"
  "fmt"
  "encoding/csv"
  "os"
)

type Records map[string]string

func (r Records) Add(domain string, ipaddr string) error {
  for key, _ := range r {
    if key == domain {
      return errors.New("Domain record already exists")
    }
  }

  // Populate records map
  r[domain] = ipaddr
  return nil
}

func NewRecordsFromCSV(csvfilepath string) (Records, error) {

  f, err := os.Open(csvfilepath)
  if err != nil {
    return nil, errors.New(fmt.Sprintf("Error opening CSV file: %s", err))
  }
  defer f.Close()
  csvReader := csv.NewReader(f)
  
  data, err := csvReader.ReadAll()
  if err != nil {
    return nil, errors.New(fmt.Sprintf("Error reading CSV file: %s", err))
  }

  var records = make(Records)
  for i, line := range data {
    var domainFromCsv string
    var ipv4FromCsv string
    if i > 0 {
      for j, field := range line {
        if j == 0 {
          domainFromCsv = field
        } else if j== 1 {
          ipv4FromCsv = field
        }
      }

      // Create new record on the map
      records.Add(domainFromCsv, ipv4FromCsv)
    }
  }

  return records,nil
}
