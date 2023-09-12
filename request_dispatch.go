package request_dispatch

import (
	"context"
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

type Dispatch struct {
	logger *Logger
	config *Config
	next   http.Handler
	rand   *rand.Rand
	mu     sync.Mutex
}

func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	plugin := &Dispatch{
		logger: NewLogger(config.LogLevel),
		config: config,
		next:   next,
		rand:   rand.New(rand.NewSource(time.Now().UnixNano())),
	}

	return plugin, nil
}

func (d *Dispatch) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	markValue := req.Header.Get(d.config.MarkHeader)
	if markValue != "" {
		if hosts, ok := d.config.MarkHosts[markValue]; ok {
			d.mu.Lock()
			host := hosts[d.rand.Intn(len(hosts))]
			d.mu.Unlock()
			target, err := url.ParseRequestURI(host)
			if err == nil {
				d.reverseProxy(rw, req, target)
				return
			} else {
				d.logger.Error("failed to parse request uri, uri: ", host, "error: ", err)
			}
		}
	}

	// default router
	d.next.ServeHTTP(rw, req)
	return
}

func (d *Dispatch) reverseProxy(rw http.ResponseWriter, req *http.Request, target *url.URL) {
	d.logger.Debug("reverse proxy, from: ", req.URL, "to: ", target)

	// replace the host of request, otherwise it will cause incorrect parsing
	req.Host = target.Host
	// FIXME cache proxy?
	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.ErrorHandler = func(writer http.ResponseWriter, request *http.Request, err error) {
		d.logger.Error("failed to reverse proxy, will be use default ServeHTTP, error: ", err)
		d.next.ServeHTTP(rw, req)
	}
	proxy.ServeHTTP(rw, req)
}

func (d *Dispatch) Close() error {
	// Close your logger and any other resources.
	return nil
}
