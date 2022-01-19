package request_dispatch

import (
	"context"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
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
}

func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	plugin := &Dispatch{
		logger: NewLogger(config.LogLevel),
		config: config,
		next:   next,
	}

	return plugin, nil
}

func (d *Dispatch) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	markValue := req.Header.Get(d.config.MarkHeader)
	if markValue != "" {
		if hosts, ok := d.config.MarkHosts[markValue]; ok {
			host := hosts[rand.Intn(len(hosts))]
			target, err := url.ParseRequestURI(host)
			if err == nil {
				d.reverseProxy(rw, req, target)
				return
			} else {
				d.logger.Debug("ParseRequestURI failed", "error", err, "url", host)
			}
		}
	}

	// default
	d.next.ServeHTTP(rw, req)
	return
}

func (d *Dispatch) reverseProxy(rw http.ResponseWriter, req *http.Request, target *url.URL) {
	d.logger.Debug("ReverseProxy", "from", req.URL, "to", target)

	// replace the host of request, otherwise it will cause incorrect parsing
	req.Host = target.Host
	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.ServeHTTP(rw, req)
}
