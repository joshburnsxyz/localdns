package record

type Records map[string]string

func (r *Records) Add(domain string, ipaddr string) error {
  // TODO: check if domain is already present and 
  // return an error if so.
  r[domain] = ipaddr
  return nil
}
