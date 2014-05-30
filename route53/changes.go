package route53

import (
	"bytes"
	"encoding/xml"
	"net/http"
	"net/url"
)

func (z *Zone) NewBatch() *ChangeBatch { return &ChangeBatch{zone: z} }

func (c *ChangeBatch) SetComment(comment string) *ChangeBatch { c.Comment = comment; return c }

func (c *ChangeBatch) CreateRecord(r *Record) *ChangeBatch {
	c.Changes = append(c.Changes, &Change{Action: "CREATE", Record: r})
	return c
}

func (c *ChangeBatch) DeleteRecord(r *Record) *ChangeBatch {
	c.Changes = append(c.Changes, &Change{Action: "DELETE", Record: r})
	return c
}

func (c *ChangeBatch) UpsertRecord(r *Record) *ChangeBatch {
	c.Changes = append(c.Changes, &Change{Action: "UPSERT", Record: r})
	return c
}

func (c *ChangeBatch) Apply() (*ChangeInfo, error) {
	dat, _ := xml.Marshal(
		&struct {
			XMLName     xml.Name `xml:"https://route53.amazonaws.com/doc/2013-04-01/ ChangeResourceRecordSetsRequest"`
			ChangeBatch *ChangeBatch
		}{ChangeBatch: c},
	)

	u, _ := url.Parse(APIRoot + c.zone.Id + "/rrset")
	req, _ := http.NewRequest("POST", u.String(), bytes.NewReader(dat))
	c.zone.conn.Sign(req)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := checkError(resp); err != nil {
		return nil, err
	}

	rsp := &struct{ ChangeInfo *ChangeInfo }{}
	xml.NewDecoder(resp.Body).Decode(rsp)

	rsp.ChangeInfo.zone = c.zone

	return rsp.ChangeInfo, nil
}

func (c *ChangeInfo) Update() (*ChangeInfo, error) {
	u, _ := url.Parse(APIRoot + c.Id)
	req, _ := http.NewRequest("GET", u.String(), nil)
	c.zone.conn.Sign(req)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := checkError(resp); err != nil {
		return nil, err
	}

	rsp := &struct{ ChangeInfo *ChangeInfo }{}
	xml.NewDecoder(resp.Body).Decode(rsp)

	rsp.ChangeInfo.zone = c.zone

	return rsp.ChangeInfo, nil
}
