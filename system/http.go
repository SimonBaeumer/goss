package system

import (
	"crypto/tls"
	"io"
	"net/http"
	"time"
	"github.com/SimonBaeumer/goss/util"
)

type Header map[string][]string

type HTTP interface {
	HTTP() string
	Status() (int, error)
	Body() (io.Reader, error)
	Exists() (bool, error)
	SetAllowInsecure(bool)
	SetNoFollowRedirects(bool)
	Headers() (Header, error)
}

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
}

func NewDefHTTP(http string, system *System, config util.Config) HTTP {
	return &DefHTTP{
		http:              http,
		allowInsecure:     config.AllowInsecure,
		noFollowRedirects: config.NoFollowRedirects,
		Timeout:           config.Timeout,
		Username:		   config.Username,
		Password:          config.Password,
		RequestHeaders:    config.RequestHeaders,
	}
}

//The setup method configures the http client and sends the request.
func (u *DefHTTP) setup() error {
	if u.loaded {
		return u.err
	}
	u.loaded = true

	tr := &http.Transport{
		TLSClientConfig:   &tls.Config{InsecureSkipVerify: u.allowInsecure},
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

func (u *DefHTTP) Exists() (bool, error) {
	if _, err := u.Status(); err != nil {
		return false, err
	}
	return true, nil
}

func (u *DefHTTP) SetNoFollowRedirects(t bool) {
	u.noFollowRedirects = t
}

func (u *DefHTTP) SetAllowInsecure(t bool) {
	u.allowInsecure = t
}

func (u *DefHTTP) ID() string {
	return u.http
}
func (u *DefHTTP) HTTP() string {
	return u.http
}

func (u *DefHTTP) Status() (int, error) {
	if err := u.setup(); err != nil {
		return 0, err
	}

	return u.resp.StatusCode, nil
}

func (u *DefHTTP) Body() (io.Reader, error) {
	if err := u.setup(); err != nil {
		return nil, err
	}

	return u.resp.Body, nil
}

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