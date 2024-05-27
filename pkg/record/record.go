package record

import (
  "errors"
)

type Records map[string]string

func (r *Records) Add(domain string, ipaddr string) error {
  for key, val := range r {
    if key == domain {
      return errors.New("Domain record already exists")
    }
  }

  // Populate records map
  r[domain] = ipaddr
  return nil
}
