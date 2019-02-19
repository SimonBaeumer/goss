package system

import (
	"crypto/tls"
	"github.com/SimonBaeumer/goss/util"
	"io"
	"net/http"
	"time"
)

// Header is an alias for the header type
type Header map[string][]string

// HTTP defines the interface to access the request data
type HTTP interface {
	HTTP() string
	Status() (int, error)
	Body() (io.Reader, error)
	Exists() (bool, error)
	SetAllowInsecure(bool)
	SetNoFollowRedirects(bool)
	Headers() (Header, error)
}

// DefHTTP is the system package representation
type DefHTTP struct {
	http              string
	allowInsecure     bool
	noFollowRedirects bool
	resp              *http.Response
	Timeout           int
	loaded            bool
	err               error
	Username          string
	Password          string
	RequestHeaders    Header
	ClientCertificate tls.Certificate
}

// NewDefHTTP is the constructor of the DefHTTP struct
func NewDefHTTP(http string, system *System, config util.Config) HTTP {
	return &DefHTTP{
		http:              http,
		allowInsecure:     config.AllowInsecure,
		noFollowRedirects: config.NoFollowRedirects,
		Timeout:           config.Timeout,
		Username:		   config.Username,
		Password:          config.Password,
		RequestHeaders:    config.RequestHeaders,
		ClientCertificate: config.Certificate,
	}
}

//The setup method configures the http client and sends the request.
func (u *DefHTTP) setup() error {
	if u.loaded {
		return u.err
	}
	u.loaded = true


    tr := &http.Transport{
		TLSClientConfig:   &tls.Config{
			InsecureSkipVerify: u.allowInsecure,
            Certificates: []tls.Certificate{u.ClientCertificate},
	    },
		DisableKeepAlives: true,
	}
	client := &http.Client{
		Transport: tr,
		Timeout:   time.Duration(u.Timeout) * time.Millisecond,
	}

	if u.noFollowRedirects {
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}

	req, err := http.NewRequest("GET", u.http, nil)

	for key, values := range u.RequestHeaders {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	if err != nil {
		return u.err
	}
	if u.Username != "" || u.Password != "" {
		req.SetBasicAuth(u.Username, u.Password)
	}
	u.resp, u.err = client.Do(req)

	return u.err
}

// Exists checks if the given uri is reachable
func (u *DefHTTP) Exists() (bool, error) {
	if _, err := u.Status(); err != nil {
		return false, err
	}
	return true, nil
}

// SetNoFollowRedirects disables the go default to follow redirect links
func (u *DefHTTP) SetNoFollowRedirects(t bool) {
	u.noFollowRedirects = t
}

// SetAllowInsecure allows bad ssl certificates
func (u *DefHTTP) SetAllowInsecure(t bool) {
	u.allowInsecure = t
}

// ID returns the id of the http resource
func (u *DefHTTP) ID() string {
	return u.http
}

// HTTP returns the url
func (u *DefHTTP) HTTP() string {
	return u.http
}

// Status returns the http status code
func (u *DefHTTP) Status() (int, error) {
	if err := u.setup(); err != nil {
		return 0, err
	}

	return u.resp.StatusCode, nil
}

// Body returns the body of the http response
func (u *DefHTTP) Body() (io.Reader, error) {
	if err := u.setup(); err != nil {
		return nil, err
	}

	return u.resp.Body, nil
}

// Headers returns the headers of the response
func (u *DefHTTP) Headers() (Header, error) {
	if err := u.setup(); err != nil {
		return nil, err
	}

	headers := make(Header)
	for k, v := range u.resp.Header {
		headers[k] = v
	}

	return headers, nil
}