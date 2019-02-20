package main

import (
    "crypto/tls"
    "crypto/x509"
    "io/ioutil"
    "log"
    "net/http"
)

// Create a minimal https server with
// client cert authentication
func main() {
    // Add root ca
    caCert, _ := ioutil.ReadFile("ca.crt")
    caCertPool := x509.NewCertPool()
    caCertPool.AppendCertsFromPEM(caCert)

    // Add valid client certificates
    clientCert, _ := ioutil.ReadFile("client.crt")
    clientCAs := x509.NewCertPool()
    clientCAs.AppendCertsFromPEM(clientCert)

    server := &http.Server{
        Addr: ":8081",
        TLSConfig: &tls.Config{
            ClientAuth: tls.RequireAndVerifyClientCert,
            RootCAs: caCertPool,
            ClientCAs: clientCAs,
        },
    }
    server.Handler = &handler{}

    err := server.ListenAndServeTLS("server.crt", "server.key")
    if err != nil {
        log.Fatal(err.Error())
    }
}

type handler struct{}

// ServeHTTP serves the handler function
func (h *handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    if _, err := w.Write([]byte("PONG\n")); err != nil {
        log.Fatal(err.Error())
    }
}