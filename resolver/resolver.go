package resolver

import (
	"fmt"
	"log"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/miekg/dns"
)

// cache stores cached from service discovery addresses
type cache struct {
	cache map[string]recordInfo
	mtx   *sync.RWMutex
}

const recordTTL = 86400 * time.Second

// recordInfo contains information about a single instance of a service
type recordInfo struct {
	hostname string
	port     uint16
	ip       string
	elapsed  time.Time
}

// Resolver implements a DNS resolver for SRV records
type Resolver struct {
	config    dns.ClientConfig
	client    *dns.Client
	addrCache cache
}

// NewResolver instantiates a new Resolver instance
func NewResolver(port int, servers ...string) *Resolver {
	r := &Resolver{
		config: dns.ClientConfig{Servers: servers, Port: fmt.Sprintf("%d", port)},
		client: new(dns.Client),
		addrCache: cache{
			cache: make(map[string]recordInfo),
			mtx:   &sync.RWMutex{},
		},
	}
	return r
}

// lookup performs an SRV lookup against the given name
func (resolver *Resolver) lookup(name string) ([]recordInfo, error) {
	m := new(dns.Msg)

	m.SetQuestion(dns.Fqdn(name), dns.TypeSRV)
	m.SetEdns0(4096, true)
	r, _, err := resolver.client.Exchange(m, resolver.config.Servers[0]+":"+resolver.config.Port)
	if err != nil {
		return nil, err
	}

	if r.Rcode != dns.RcodeSuccess {
		return nil, fmt.Errorf("lookup failed: %d", r.Rcode)
	}

	info := map[string]recordInfo{}

	for _, k := range r.Extra {
		switch k.(type) {
		case *dns.A:
			a := k.(*dns.A)
			info[a.Hdr.Name] = recordInfo{ip: a.A.String(), hostname: a.Hdr.Name}
		}
	}

	for _, k := range r.Answer {
		if key, ok := k.(*dns.SRV); ok {
			if v, ok := info[key.Target]; ok {
				v.port = key.Port
				info[key.Target] = v
			}
		}
	}

	out := make([]recordInfo, 0, len(info))
	for _, v := range info {
		out = append(out, v)
	}

	return out, nil
}

func (c *cache) deleteRecord(uri string) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	delete(c.cache, uri)
}

func (c *cache) getRecord(uri string) (string, bool) {
	var (
		dnsRecord = recordInfo{}
		hasRecord bool
	)

	c.mtx.RLock()
	if _, hasRecord = c.cache[uri]; hasRecord {
		dnsRecord = c.cache[uri]
	}
	c.mtx.RUnlock()

	if !hasRecord {
		return "", false
	}

	now := time.Now()
	recordAge := now.Sub(dnsRecord.elapsed)
	if recordAge >= recordTTL {
		c.deleteRecord(uri)
		return "", false
	}
	return dnsRecord.hostname, true
}

func (c *cache) addRecord(hostname, ip string, port uint16) {
	addr := hostname + ":" + strconv.Itoa(int(port))
	c.mtx.Lock()
	defer c.mtx.Unlock()

	c.cache[addr] = recordInfo{
		hostname: addr,
		ip:       ip,
		port:     port,
		elapsed:  time.Now(),
	}
}

// Discover provides server discovery for floop
func (resolver *Resolver) Discover(uri string) (string, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return "", err
	}

	addr := u.Hostname() + ":" + u.Port()
	if record, ok := resolver.addrCache.getRecord(addr); ok {
		return record, nil
	}

	// Perform SRV lookup
	rslvResp, err := resolver.lookup(u.Hostname())
	if err != nil {
		return "", err
	}

	if len(rslvResp) == 0 {
		return "", fmt.Errorf("service not found")
	}

	u.Host = rslvResp[0].hostname + ":" + strconv.Itoa(int(rslvResp[0].port))

	for _, r := range rslvResp {
		log.Printf("[DEBUG] host=%s ip=%s port=%d\n", r.hostname, r.ip, r.port)
		resolver.addrCache.addRecord(r.hostname, r.ip, r.port)
	}
	return u.String(), nil
}
