package agentclient

import (
	"context"
	"net/http"
	"time"

	"github.com/ironzhang/superlib/httputils/httpclient"
	"github.com/ironzhang/superlib/timeutil"
)

// Options agent client options.
type Options struct {
	Addr    string
	Timeout time.Duration
}

// Client is a client to call agent api.
type Client struct {
	hc httpclient.Client
}

// New returns an instance of agent client.
func New(opts Options) *Client {
	return &Client{
		hc: httpclient.Client{
			Addr: opts.Addr,
			Client: http.Client{
				Timeout: opts.Timeout,
			},
		},
	}
}

// SubscribeDomains subscribes the given domains.
func (p *Client) SubscribeDomains(ctx context.Context, domains []string, ttl time.Duration, waitForReady time.Duration) error {
	req := _SubscribeDomainsReq{
		Domains:      domains,
		TTL:          timeutil.Duration(ttl),
		WaitForReady: timeutil.Duration(waitForReady),
	}
	return p.hc.Post(ctx, "/superdns/agent/v1/api/subscribe/domains", nil, req, nil)
}
