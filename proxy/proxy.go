package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/streadway/amqp"

    "github.com/mishas/prometheus_amqp_proxy/proxy/config"
    "github.com/mishas/prometheus_amqp_proxy/proxy/rpc"
)

type externalAuth struct{}

func (a *externalAuth) Mechanism() string {
	return "EXTERNAL"
}
func (a *externalAuth) Response() string {
	return fmt.Sprintf("\000")
}

// getAMQPConfig returns a reference to the amqp.Config object.
// certsDir should be the path to the directory holding {cacert,cert,key}.pem files for the TLS
// connection, or nil for no TLS.
// serverName should be set to the CN defined in the certificates expected to be received from the
// server, or nil, if CN is the same as the server's DNS name.
func getAMQPConfig(cfg config.TLSConfig) (*amqp.Config, error) {
	tlscfg := new(tls.Config)
	tlscfg.RootCAs = x509.NewCertPool()
	if ca, err := ioutil.ReadFile(cfg.CAFile); err == nil {
		tlscfg.RootCAs.AppendCertsFromPEM(ca)
	} else {
		return nil, fmt.Errorf("Failed reading CA certificate: %v", err)
	}

	if cert, err := tls.LoadX509KeyPair(cfg.CertFile, cfg.KeyFile); err == nil {
		tlscfg.Certificates = append(tlscfg.Certificates, cert)
	} else {
		return nil, fmt.Errorf("Failed reading client certificate: %v", err)
	}

	tlscfg.ServerName = cfg.ServerName

	return &amqp.Config{
		SASL:            []amqp.Authentication{&externalAuth{}},
		TLSClientConfig: tlscfg,
	}, nil
}

type handler struct {
	c      *rpc.RpcClient
	prefix string
}

func newHandler(scrapeCfg config.ScrapeConfig) handler {
	cfg, err := getAMQPConfig(scrapeCfg.TLSConfig)
	if err != nil {
		log.Fatalf("Failed creating AMQP config: %v", err)
	}

	c := rpc.NewRPCClient(scrapeCfg.AMQPConfig.URL, scrapeCfg.AMQPConfig.Exchange, cfg)
	if err := c.Init(); err != nil {
		log.Fatalf("Failed to create RPCClient: %v", err)
	}

	return handler{c, scrapeCfg.JobName}
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		log.Printf("Got unexpected method: %s. Full request: %q", r.Method, r)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	s := strings.Split(r.Host, ":")
	port := s[len(s)-1]

	ch, err := h.c.Call(h.prefix + ":" + port)
	if err != nil {
		log.Printf("Failed to publish: %v", err)
		return
	}
	if res := <-ch; res == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if _, err := w.Write(res); err != nil {
		log.Printf("Failed to get response: %v", err)
	}
}

func runListeners(addr string) {
	// TODO: Should this be http.ListenAndServeTLS?
	log.Fatal(http.ListenAndServe(addr, nil))
}

func main() {
	cfg, err := config.LoadFile("./config.yaml")
	if err != nil {
		log.Fatalf("Failed loading config: %v", err)
	}

	if cfgCount := len(cfg.ScrapeConfigs); cfgCount != 1 {
		log.Fatalf("Only supporting 1 scrape_config, %d supplied.", cfgCount)
	}

	scrapeCfg := cfg.ScrapeConfigs[0]

	println("Starting main")
	http.Handle("/metrics", newHandler(*scrapeCfg))

	for _, targetGroup := range scrapeCfg.TargetGroups {
		for _, target := range targetGroup.Targets {
			s := strings.Split(target, ":")
			go runListeners(":" + s[len(s)-1])
		}
	}

	select {}
}
