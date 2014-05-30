package route53

import (
	"encoding/xml"
	"net/http"
	"net/url"
)

// RecordsIterator последовательно перебирает записи в зоне
type RecordsIterator iterator

func (z *Zone) IterateRecords() *RecordsIterator {
	return (*RecordsIterator)(newIterator(&zoneFetcher{Zone: z}))
}

func (r *RecordsIterator) Record() *Record { rec, _ := (*iterator)(r).Value().(*Record); return rec }
func (r *RecordsIterator) Next() bool      { return (*iterator)(r).Next() }
func (r *RecordsIterator) Error() error    { return (*iterator)(r).Error() }

func (z *Zone) Records() (records []*Record, err error) {
	iter := z.IterateRecords()
	for iter.Next() {
		records = append(records, iter.Record())
	}
	err = iter.Error()
	return
}

func (z *Zone) FindRecord(name, rType, ident string) (rec *Record, err error) {
	iter := z.IterateRecords()
	for iter.Next() {
		r := iter.Record()
		if r.Name == name && r.Type == rType && r.SetIdentifier == ident {
			rec = r
			return
		}
	}
	err = iter.Error()
	return
}

func (r *Record) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	r1 := &encRecord{
		Name:          r.Name,
		Type:          r.Type,
		TTL:           r.TTL,
		HealthCheckId: r.HealthCheckId,
		SetIdentifier: r.SetIdentifier,
		Weight:        r.Weight,
		AliasTarget:   r.AliasTarget,
		Region:        r.Region,
		Failover:      r.Failover,
	}
	r1.ResourceRecords = make([]encResourceRecord, len(r.Values))
	for i, v := range r.Values {
		r1.ResourceRecords[i].Value = v
	}
	return e.EncodeElement(r1, start)
}

////////////////

type zoneFetcher struct {
	*Zone
	nextName  string
	nextType  string
	nextIdent string
}

type zoneListResponse struct {
	Records              []*Record `xml:"ResourceRecordSets>ResourceRecordSet"`
	NextRecordName       string
	NextRecordType       string
	NextRecordIdentifier string
	IsTruncated          bool
	MaxItems             int
}

func (z *zoneFetcher) itFetch() (vals []interface{}, isLast bool, eerr error) {
	u, _ := url.Parse(APIRoot + z.Id + "/rrset")
	q := make(url.Values)
	setIf(q, "name", z.nextName)
	setIf(q, "type", z.nextType)
	setIf(q, "identifier", z.nextIdent)
	u.RawQuery = q.Encode()

	req, _ := http.NewRequest("GET", u.String(), nil)
	z.conn.Sign(req)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		eerr = err
		return
	}
	defer resp.Body.Close()

	if err := checkError(resp); err != nil {
		eerr = err
		return
	}

	rsp := &zoneListResponse{}
	xml.NewDecoder(resp.Body).Decode(rsp)

	vals = make([]interface{}, len(rsp.Records))
	for i, v := range rsp.Records {
		v.Name = unescape(v.Name)
		vals[i] = v
	}

	isLast = !rsp.IsTruncated

	return
}
