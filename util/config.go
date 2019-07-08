package util

import (
	"crypto/tls"
	"fmt"
	"reflect"
	"strings"

	"github.com/oleiade/reflections"
)

// Config type is a helper type which holds the configuration of the given resource
type Config struct {
	IgnoreList        []string
	Timeout           int
	AllowInsecure     bool
	NoFollowRedirects bool
	Server            string
	Username          string
	Password          string
	Header            map[string][]string
	RequestHeaders    map[string][]string
	Certificate       tls.Certificate
}

type Request struct {
}

type OutputConfig struct {
	FormatOptions []string
}

type format string

const (
	JSON format = "json"
	YAML format = "yaml"
)

// ValidateSections validates the sections of the config file
func ValidateSections(unmarshal func(interface{}) error, i interface{}, whitelist map[string]bool) error {
	// Get generic input
	var toValidate map[string]map[string]interface{}
	if err := unmarshal(&toValidate); err != nil {
		return err
	}

	// Run input through whitelist
	typ := reflect.TypeOf(i)
	typs := strings.Split(typ.String(), ".")[1]
	for id, v := range toValidate {
		for k, _ := range v {
			if !whitelist[k] {
				return fmt.Errorf("Invalid Attribute for %s:%s: %s", typs, id, k)
			}
		}
	}

	return nil
}

func WhitelistAttrs(i interface{}, format format) (map[string]bool, error) {
	validAttrs := make(map[string]bool)
	tags, err := reflections.Tags(i, string(format))
	if err != nil {
		return nil, err
	}
	for _, v := range tags {
		validAttrs[strings.Split(v, ",")[0]] = true
	}
	return validAttrs, nil
}

// IsValueInList checks if a value is in the string slice
func IsValueInList(value string, list []string) bool {
	for _, v := range list {
		if strings.ToLower(v) == strings.ToLower(value) {
			return true
		}
	}
	return false
}
