package agentclient

import "github.com/ironzhang/superlib/timeutil"

type _SubscribeDomainsReq struct {
	Domains      []string          // the domain list that require to subscribe
	TTL          timeutil.Duration // time to live, <= 0 means forever
	Asynchronous bool              // asynchronous call
}
