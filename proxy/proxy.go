package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/streadway/amqp"

	"github.com/mishas/prometheus_amqp_proxy/proxy/rpc"
)

var (
	amqpURL      = flag.String("amqp_url", "", "URL of the AMQP server to dial")
	amqpExchange = flag.String("amqp_exchange", "", "Name of the AMQP exchange to use")

	certsDir = flag.String("certs_dir", "", "Directory of certs for TLS connection to AMQP, or empty for non-TLS connection. "+
		"Expected files are: cacert.pem, cert.pem and key.pem.")
	serverName = flag.String("server_name", "", "Name of the server for TLS verification, or nil for default")
)

type externalAuth struct{}

func (a *externalAuth) Mechanism() string {
	return "EXTERNAL"
}
func (a *externalAuth) Response() string {
	return fmt.Sprintf("\000")
}

// getAMQPConfig returns a reference to the amqp.Config object.
func getAMQPConfig() (*amqp.Config, error) {
	tlscfg := new(tls.Config)
	if *certsDir != "" {
		tlscfg.RootCAs = x509.NewCertPool()
		if ca, err := ioutil.ReadFile(*certsDir + "/cacert.pem"); err == nil {
			tlscfg.RootCAs.AppendCertsFromPEM(ca)
		} else {
			return nil, fmt.Errorf("Failed reading CA certificate: %v", err)
		}

		if cert, err := tls.LoadX509KeyPair(*certsDir+"/cert.pem", *certsDir+"/key.pem"); err == nil {
			tlscfg.Certificates = append(tlscfg.Certificates, cert)
		} else {
			return nil, fmt.Errorf("Failed reading client certificate: %v", err)
		}

		if *serverName != "" {
			tlscfg.ServerName = *serverName
		}
	}

	return &amqp.Config{
		SASL:            []amqp.Authentication{&externalAuth{}},
		TLSClientConfig: tlscfg,
	}, nil
}

type handler struct {
	c *rpc.RpcClient
}

func newHandler() handler {
	cfg, err := getAMQPConfig()
	if err != nil {
		log.Fatalf("Failed creating AMQP config: %v", err)
	}

	c := rpc.NewRPCClient(*amqpURL, *amqpExchange, cfg)
	if err := c.Init(); err != nil {
		log.Fatalf("Failed to create RPCClient: %v", err)
	}

	return handler{c}
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		log.Printf("Got unexpected method: %s. Full request: %q", r.Method, r)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	params := r.URL.Query()
	target := params.Get("target")

	ch, err := h.c.Call(target)
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

func main() {
	flag.Parse()

	if !strings.HasPrefix(*amqpURL, "amqp://") && !strings.HasPrefix(*amqpURL, "amqps://") {
		fmt.Println("Please provide a valid URL for an AMQP server with the -amqp_url flag.")
		os.Exit(1)
	}

	if *amqpExchange == "" {
		fmt.Println("Please provide an AMQP exchange to use with the -amqp_exchange flag.")
		os.Exit(1)
	}

	println("Starting main")
	http.Handle("/proxy", newHandler())

	log.Fatal(http.ListenAndServe(":8200", nil))
}
