package route53

import (
	"bytes"
	"encoding/xml"
	"net/http"
	"net/url"
)

func (c *Conn) CreateZone(name, callerRef, comment string) (*HostedZoneCreated, error) {
	dat, _ := xml.Marshal(
		&struct {
			XMLName   xml.Name `xml:"https://route53.amazonaws.com/doc/2013-04-01/ CreateHostedZoneRequest"`
			Name      string
			CallerRef string `xml:"CallerReference"`
			Comment   string `xml:"HostedZoneConfig>Comment"`
		}{Name: name, CallerRef: callerRef, Comment: comment},
	)

	u, _ := url.Parse(APIRoot + "/hostedzone")
	req, _ := http.NewRequest("POST", u.String(), bytes.NewReader(dat))
	c.Sign(req)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := checkError(resp); err != nil {
		return nil, err
	}

	rsp := &HostedZoneCreated{}
	xml.NewDecoder(resp.Body).Decode(rsp)

	rsp.Zone.conn = c

	return rsp, nil
}

func (z *Zone) Delete() (*ChangeInfo, error) {
	u, _ := url.Parse(APIRoot + z.Id)
	req, _ := http.NewRequest("DELETE", u.String(), nil)
	z.conn.Sign(req)

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

	return rsp.ChangeInfo, nil
}

// LoadDetails загружает дополнительные данные о зоне
func (z *Zone) Details() (*HostedZoneDetails, error) {
	u, _ := url.Parse(APIRoot + z.Id)
	req, _ := http.NewRequest("GET", u.String(), nil)
	z.conn.Sign(req)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := checkError(resp); err != nil {
		return nil, err
	}

	rsp := &HostedZoneDetails{}
	xml.NewDecoder(resp.Body).Decode(rsp)
	rsp.Zone.conn = z.conn

	return rsp, nil
}

func (c *Conn) ZoneById(id string) *Zone { return &Zone{Id: "/hostedzone/" + id, conn: c} }
