package balancer

import (
	"fmt"
	"net"
	"net/url"

	b "github.com/divideandconquer/go-consul-client/src/balancer"
)

func NewK8SBalancer() b.DNS {
	return &k8sBalancer{}
}

type k8sBalancer struct{}

func (d *k8sBalancer) FindService(serviceName string) (*b.ServiceLocation, error) {
	_, addrs, err := net.LookupSRV("", "", serviceName)
	if err != nil {
		return nil, err
	}

	if len(addrs) == 0 {
		return nil, fmt.Errorf("No addresses found")
	}

	return &b.ServiceLocation{
		URL:  addrs[0].Target,
		Port: int(addrs[0].Port), // just take the first one, DNS will return them randomly-ish.  TODO: do actual load balancing
	}, nil

}

func (d *k8sBalancer) GetHttpUrl(serviceName string, useTLS bool) (url.URL, error) {
	u := url.URL{}

	//loc, err := d.FindService(serviceName)
	//if err != nil {
	//	return u, err
	//}

	u.Scheme = "http"
	if useTLS {
		u.Scheme = "https"
	}
	u.Host = fmt.Sprintf("%s", serviceName)

	return u, nil
}
