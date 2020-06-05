package balancer

import (
	"fmt"
	"net"
	"net/url"

	b "github.com/divideandconquer/go-consul-client/src/balancer"
)

// portName should match the "service" portion of the SRV record, e.g. _http._tcp.profile... or _postgres._tcp.profile-db...
func NewSRVBalancer(portName string) b.DNS {
	return &dnsBalancer{
		portName: portName,
	}
}

type dnsBalancer struct {
	portName string
}

func (d *dnsBalancer) FindService(serviceName string) (*b.ServiceLocation, error) {
	cname, addrs, err := net.LookupSRV(d.portName, "tcp", serviceName)
	if err != nil {
		return nil, err
	}

	if len(addrs) == 0 {
		return nil, fmt.Errorf("No addresses found")
	}

	return &b.ServiceLocation{
		URL:  cname,
		Port: int(addrs[0].Port), // just take the first one, DNS will return them randomly-ish.  TODO: do actual load balancing
	}, nil

}

func (d *dnsBalancer) GetHttpUrl(serviceName string, useTLS bool) (url.URL, error) {
	u := url.URL{}

	loc, err := d.FindService(serviceName)
	if err != nil {
		return u, err
	}

	u.Scheme = "http"
	if useTLS {
		u.Scheme = "https"
	}
	u.Host = fmt.Sprintf("%s:%d", loc.URL, loc.Port)

	return u, nil
}
