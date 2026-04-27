package request_dispatch

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"
)

type Config struct {
	LogLevel   string              `json:"logLevel"`
	MarkHeader string              `json:"markHeader"`
	MarkHosts  map[string][]string `json:"markHosts"`
}

func CreateConfig() *Config {
	return &Config{}
}

type RequestDispatcher struct {
	logger        *Logger
	config        *Config
	next          http.Handler
	randSource    *rand.Rand
	randMutex     sync.Mutex
	reverseProxies map[string]*httputil.ReverseProxy
}

func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	if err := validateConfig(config); err != nil {
		return nil, err
	}

	reverseProxies := make(map[string]*httputil.ReverseProxy)
	for _, hosts := range config.MarkHosts {
		for _, host := range hosts {
			target, err := url.ParseRequestURI(host)
			if err != nil {
				return nil, err
			}
			reverseProxies[host] = httputil.NewSingleHostReverseProxy(target)
		}
	}

	dispatcher := &RequestDispatcher{
		logger:         NewLogger(config.LogLevel),
		config:         config,
		next:           next,
		randSource:     rand.New(rand.NewSource(time.Now().UnixNano())),
		reverseProxies: reverseProxies,
	}

	return dispatcher, nil
}

func (d *RequestDispatcher) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	markValue := req.Header.Get(d.config.MarkHeader)
	if markValue != "" {
		if hosts, ok := d.config.MarkHosts[markValue]; ok {
			d.randMutex.Lock()
			host := hosts[d.randSource.Intn(len(hosts))]
			d.randMutex.Unlock()
			if proxy, ok := d.reverseProxies[host]; ok {
				d.proxyRequest(rw, req, host, proxy)
				return
			}
		}
	}

	// default router
	d.next.ServeHTTP(rw, req)
}

func (d *RequestDispatcher) proxyRequest(rw http.ResponseWriter, req *http.Request, host string, proxy *httputil.ReverseProxy) {
	target, _ := url.ParseRequestURI(host)
	d.logger.Debug("reverse proxy, from: ", req.URL, "to: ", target)

	// Set request host to target host for correct URL parsing in reverse proxy
	req.Host = target.Host
	proxy.ErrorHandler = func(writer http.ResponseWriter, request *http.Request, err error) {
		d.logger.Error("failed to reverse proxy, will use default ServeHTTP, error: ", err)
		d.next.ServeHTTP(rw, req)
	}
	proxy.ServeHTTP(rw, req)
}

func validateConfig(config *Config) error {
	if config.MarkHeader == "" {
		return nil
	}
	if len(config.MarkHosts) == 0 {
		return nil
	}
	for mark, hosts := range config.MarkHosts {
		if len(hosts) == 0 {
			return fmt.Errorf("mark '%s' has no hosts configured", mark)
		}
		for _, host := range hosts {
			if _, err := url.ParseRequestURI(host); err != nil {
				return fmt.Errorf("invalid URL for mark '%s': %s, error: %w", mark, host, err)
			}
		}
	}
	return nil
}

func (d *RequestDispatcher) Close() error {
	// Close releases any resources held by the dispatcher.
	return nil
}
