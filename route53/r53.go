package route53

import "github.com/davidmz/aws"

type Conn struct{ keys *aws.Keys }

func New(keys *aws.Keys) *Conn { return &Conn{keys: keys} }
