package route53

import (
	"encoding/xml"
	"net/http"
	"net/url"
)

// ZonesIterator последовательно перебирает зоны
type ZonesIterator iterator

// IterateZones возвращает объект ZonesIterator для последовательного перебора зон
func (c *Conn) IterateZones() *ZonesIterator {
	return (*ZonesIterator)(newIterator(&zonesFetcher{Conn: c}))
}

// Next переходит к следующей зоне и возвращает false если список зон закончился или возникал ошибка
func (z *ZonesIterator) Zone() *Zone { rec, _ := (*iterator)(z).Value().(*Zone); return rec }

// Zone возвращает текущую зону
func (z *ZonesIterator) Next() bool { return (*iterator)(z).Next() }

// Error возвращает ошибку, если она произошла при переборе зон или nil, если всё было в порядке
func (z *ZonesIterator) Error() error { return (*iterator)(z).Error() }

// Zones возвращает список всех зон
func (c *Conn) Zones() (zones []*Zone, err error) {
	iter := c.IterateZones()
	for iter.Next() {
		zones = append(zones, iter.Zone())
	}
	err = iter.Error()
	return
}

func (c *Conn) FindZone(name string) (zone *Zone, err error) {
	iter := c.IterateZones()
	for iter.Next() {
		if iter.Zone().Name == name {
			zone = iter.Zone()
			return
		}
	}
	err = iter.Error()
	return
}

////////////////////////

type zonesFetcher struct {
	*Conn
	nextMarker string
}

type listHostedZonesResponse struct {
	Marker      string
	NextMarker  string
	IsTruncated bool
	MaxItems    int
	Zones       []*Zone `xml:"HostedZones>HostedZone"`
}

func (z *zonesFetcher) itFetch() (vals []interface{}, isLast bool, eerr error) {
	u, _ := url.Parse(APIRoot + "/hostedzone")
	q := make(url.Values)
	setIf(q, "marker", z.nextMarker)
	u.RawQuery = q.Encode()

	req, _ := http.NewRequest("GET", u.String(), nil)
	z.Sign(req)

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

	rsp := &listHostedZonesResponse{}
	xml.NewDecoder(resp.Body).Decode(rsp)

	vals = make([]interface{}, len(rsp.Zones))
	for i, v := range rsp.Zones {
		v.conn = z.Conn
		vals[i] = v
	}

	isLast = !rsp.IsTruncated

	return
}
