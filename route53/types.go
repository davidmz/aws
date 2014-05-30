package route53

import "time"

type Zone struct {
	conn *Conn

	Id        string
	Name      string
	CallerRef string `xml:"CallerReference"`
	Comment   string `xml:"Config>Comment"`
	RSCount   int    `xml:"ResourceRecordSetCount"`
}

type ChangeInfo struct {
	zone *Zone

	Id          string
	Status      string
	SubmittedAt time.Time
}

type DelegationSet struct {
	NameServers []string `xml:"NameServers>NameServer"`
}

type HostedZoneDetails struct {
	Zone          *Zone `xml:"HostedZone"`
	DelegationSet *DelegationSet
}

type HostedZoneCreated struct {
	Zone          *Zone `xml:"HostedZone"`
	ChangeInfo    *ChangeInfo
	DelegationSet *DelegationSet
}

type Record struct {
	Name          string
	Type          string
	TTL           int
	Values        []string     `xml:"ResourceRecords>ResourceRecord>Value"`
	HealthCheckId string       `xml:",omitempty"`
	SetIdentifier string       `xml:",omitempty"`
	Weight        int          `xml:",omitempty"`
	AliasTarget   *AliasTarget `xml:",omitempty"`
	Region        string       `xml:",omitempty"`
	Failover      string       `xml:",omitempty"`
}

type AliasTarget struct {
	HostedZoneId         string
	DNSName              string
	EvaluateTargetHealth bool
}

type ChangeBatch struct {
	zone *Zone

	Comment string    `xml:"Comment"`
	Changes []*Change `xml:"Changes>Change"`
}

type Change struct {
	Action string  `xml:"Action"`
	Record *Record `xml:"ResourceRecordSet"`
}

///////////////////

type encRecord struct {
	Name            string
	Type            string
	TTL             int
	ResourceRecords []encResourceRecord `xml:"ResourceRecords>ResourceRecord"`
	HealthCheckId   string              `xml:",omitempty"`
	SetIdentifier   string              `xml:",omitempty"`
	Weight          int                 `xml:",omitempty"`
	AliasTarget     *AliasTarget        `xml:",omitempty"`
	Region          string              `xml:",omitempty"`
	Failover        string              `xml:",omitempty"`
}

type encResourceRecord struct{ Value string }
