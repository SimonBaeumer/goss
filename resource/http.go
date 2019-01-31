package resource

import (
    "github.com/SimonBaeumer/goss/system"
    "github.com/SimonBaeumer/goss/util"
    "reflect"
    "strings"
    "time"
)

const TimeoutMS = 5000

type HTTP struct {
	Title             string   			  	`json:"title,omitempty" yaml:"title,omitempty"`
	Meta              meta     			  	`json:"meta,omitempty" yaml:"meta,omitempty"`
	HTTP              string   			  	`json:"-" yaml:"-"`
	Status            matcher  			  	`json:"status" yaml:"status"`
	AllowInsecure     bool     			  	`json:"allow-insecure" yaml:"allow-insecure"`
	NoFollowRedirects bool     			  	`json:"no-follow-redirects" yaml:"no-follow-redirects"`
	Timeout           int      			  	`json:"timeout" yaml:"timeout"`
	Body              []string 			  	`json:"body" yaml:"body"`
	Username          string   			  	`json:"username,omitempty" yaml:"username,omitempty"`
	Password          string   			  	`json:"password,omitempty" yaml:"password,omitempty"`
	Headers			  map[string][]string   `json:"headers,omitempty" yaml:"headers,omitempty"`
	RequestHeaders 	  map[string][]string   `json:"request-headers,omitempty" yaml:"request-headers,omitempty"`
}

func (u *HTTP) ID() string      { return u.HTTP }
func (u *HTTP) SetID(id string) { u.HTTP = id }

func (r *HTTP) GetTitle() string { return r.Title }
func (r *HTTP) GetMeta() meta    { return r.Meta }

func (u *HTTP) Validate(sys *system.System) []TestResult {
	skip := false

	conf := util.Config{
		AllowInsecure: u.AllowInsecure,
		NoFollowRedirects: u.NoFollowRedirects,
		Timeout: u.Timeout,
		Username: u.Username,
		Password: u.Password,
		RequestHeaders: u.RequestHeaders,
	}

	sysHTTP := sys.NewHTTP(
		u.HTTP,
		sys,
		conf,
	)

	sysHTTP.SetAllowInsecure(u.AllowInsecure)
	sysHTTP.SetNoFollowRedirects(u.NoFollowRedirects)

	var results []TestResult
	results = append(results, ValidateValue(u, "status", u.Status, sysHTTP.Status, skip))
	if shouldSkip(results) {
		skip = true
	}
	if len(u.Body) > 0 {
		results = append(results, ValidateContains(u, "Body", u.Body, sysHTTP.Body, skip))
	}

	if len(u.Headers) > 0 {
		headers, _ := sysHTTP.Headers()
		results = append(results, validateHeader(u, "Headers", u.Headers, headers, skip))
	}

	return results
}

func NewHTTP(sysHTTP system.HTTP, config util.Config) (*HTTP, error) {
	if config.Timeout == 0 {
		config.Timeout = TimeoutMS
	}

	http := sysHTTP.HTTP()
	status, err := sysHTTP.Status()
	headers, _ := sysHTTP.Headers()
	u := &HTTP{
		HTTP:              http,
		Status:            status,
		Body:              []string{},
		AllowInsecure:     config.AllowInsecure,
		NoFollowRedirects: config.NoFollowRedirects,
		Timeout:           config.Timeout,
		Username:		   config.Username,
		Password:          config.Password,
		Headers:		   headers,
	}
	return u, err
}


func validateHeader(res ResourceRead, property string, expectedHeaders map[string][]string, actualHeaders system.Header, skip bool) TestResult {
	id := res.ID()
	title := res.GetTitle()
	meta := res.GetMeta()
	typ := reflect.TypeOf(res)
	typeS := strings.Split(typ.String(), ".")[1]
	startTime := time.Now()
	if skip {
		return skipResult(
			typeS,
			Values,
			id,
			title,
			meta,
			property,
			startTime,
		)
	}

	actualString := convertHeaderMapToString(actualHeaders)

	for expectedKey, expectedValues := range expectedHeaders {
		if _, ok := actualHeaders[expectedKey]; !ok {
			return TestResult{
				Successful:   false,
				Result:       FAIL,
				Title:        title,
				ResourceType: typeS,
				ResourceId:   id,
				TestType:     Header,
				Property:     property,
				Err:          nil,
				Human:        "Did not find header " + expectedKey + " got \n" + actualString,
				Expected:     []string{expectedKey},
				Found:        []string{actualString},
			}
		}

		actualValues := actualHeaders[expectedKey]

		for _, expectedValue := range expectedValues {
			if !isInStringSlice(actualValues, expectedValue) {
				return TestResult{
					Successful:   false,
					Result:       FAIL,
					ResourceType: typeS,
					ResourceId:   id,
					TestType:     Header,
					Title:        title,
					Property:     property,
					Err:          nil,
					Human:        "Did not find header " + expectedKey + ": " + expectedValue,
					Expected:     []string{expectedValue},
					Found:        []string{actualString},
				}
			}
		}
	}

	return TestResult{
		Successful:   true,
		Title:        title,
		ResourceId:   id,
		Result:       SUCCESS,
		ResourceType: typeS,
		TestType:     Header,
	}
}

func convertHeaderMapToString(actualHeaders system.Header) string {
	var actualString string
	for k, values := range actualHeaders {
		for _, v := range values {
			actualString += k + ": " + v + "\n"
		}
	}
	return actualString
}

func isInStringSlice(haystack []string, needle string) bool {
	for _, value := range haystack {
		if needle == value {
			return true
		}
	}
	return false
}
