package goss

import (
	"fmt"
	"github.com/SimonBaeumer/goss/resource"
	"github.com/SimonBaeumer/goss/system"
	"github.com/SimonBaeumer/goss/util"
	"io"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type Add struct {
	Writer            io.Writer
	ExcludeAttr       []string
	Timeout           int
	AllowInsecure     bool
	NoFollowRedirects bool
	Server            string
	Username          string
	Password          string
	Header            string
	Sys               *system.System
}

// AddResources is a sSimple wrapper to add multiple resources
func (a *Add) AddResources(fileName, resourceName string, keys []string) error {
	OutStoreFormat = getStoreFormatFromFileName(fileName)
	header := extractHeaderArgument(a.Header)

	config := util.Config{
		IgnoreList:        a.ExcludeAttr,
		Timeout:           a.Timeout,
		AllowInsecure:     a.AllowInsecure,
		NoFollowRedirects: a.NoFollowRedirects,
		Server:            a.Server,
		Username:          a.Username,
		Password:          a.Password,
		Header:            header,
	}

	var gossConfig GossConfig
	if _, err := os.Stat(fileName); err == nil {
		gossConfig = ReadJSON(fileName)
	} else {
		gossConfig = *NewGossConfig()
	}

	for _, key := range keys {
		if err := a.AddResource(fileName, gossConfig, resourceName, key, config); err != nil {
			return err
		}
	}
	WriteJSON(fileName, gossConfig)

	return nil
}

func extractHeaderArgument(headerArg string) map[string][]string {
	if headerArg == "" {
		return make(map[string][]string)
	}
	rawHeaders := strings.Split(headerArg, ":")
	headers := make(map[string][]string)
	headers[rawHeaders[0]] = []string{strings.TrimSpace(rawHeaders[1])}
	return headers
}

// AddResource adds a resource to the configuration file
func (a *Add) AddResource(fileName string, gossConfig GossConfig, resourceName, key string, config util.Config) error {
	// Need to figure out a good way to refactor this
	switch resourceName {
	case "addr":
		res, err := gossConfig.Addrs.AppendSysResource(key, a.Sys, config)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		a.resourcePrint(fileName, res)
	case "command":
		res, err := gossConfig.Commands.AppendSysResource(key, a.Sys, config)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
        a.resourcePrint(fileName, res)
	case "dns":
		res, err := gossConfig.DNS.AppendSysResource(key, a.Sys, config)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
        a.resourcePrint(fileName, res)
	case "file":
		res, err := gossConfig.Files.AppendSysResource(key, a.Sys, config)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
        a.resourcePrint(fileName, res)
	case "Group":
		res, err := gossConfig.Groups.AppendSysResource(key, a.Sys, config)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		a.resourcePrint(fileName, res)
	case "package":
		res, err := gossConfig.Packages.AppendSysResource(key, a.Sys, config)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
        a.resourcePrint(fileName, res)
	case "port":
		res, err := gossConfig.Ports.AppendSysResource(key, a.Sys, config)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
        a.resourcePrint(fileName, res)
	case "process":
		res, err := gossConfig.Processes.AppendSysResource(key, a.Sys, config)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
        a.resourcePrint(fileName, res)
	case "service":
		res, err := gossConfig.Services.AppendSysResource(key, a.Sys, config)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
        a.resourcePrint(fileName, res)
	case "user":
		res, err := gossConfig.Users.AppendSysResource(key, a.Sys, config)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
        a.resourcePrint(fileName, res)
	case "gossfile":
		res, err := gossConfig.Gossfiles.AppendSysResource(key, a.Sys, config)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
        a.resourcePrint(fileName, res)
	case "kernel-param":
		res, err := gossConfig.KernelParams.AppendSysResource(key, a.Sys, config)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
        a.resourcePrint(fileName, res)
	case "Mount":
		res, err := gossConfig.Mounts.AppendSysResource(key, a.Sys, config)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
        a.resourcePrint(fileName, res)
	case "interface":
		res, err := gossConfig.Interfaces.AppendSysResource(key, a.Sys, config)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		a.resourcePrint(fileName, res)
	case "http":
		res, err := gossConfig.HTTPs.AppendSysResource(key, a.Sys, config)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
        a.resourcePrint(fileName, res)
	default:
		panic("Undefined resource name: " + resourceName)
	}

	return nil
}

// Simple wrapper to add multiple resources
func (a *Add) AutoAddResources(fileName string, keys []string) error {
	OutStoreFormat = getStoreFormatFromFileName(fileName)

	var gossConfig GossConfig
	if _, err := os.Stat(fileName); err == nil {
		gossConfig = ReadJSON(fileName)
	} else {
		gossConfig = *NewGossConfig()
	}

	for _, key := range keys {
		if err := a.AutoAddResource(fileName, gossConfig, key); err != nil {
			return err
		}
	}
	if err := WriteJSON(fileName, gossConfig); err != nil {
		return err
	}

	return nil
}

// Autoadds all resources to the config file
func (a *Add) AutoAddResource(fileName string, gossConfig GossConfig, key string) error {
	// file
	if strings.Contains(key, "/") {
		if res, _, ok := gossConfig.Files.AppendSysResourceIfExists(key, a.Sys); ok == true {
            a.resourcePrint(fileName, res)
		}
	}

	// group
	if res, _, ok := gossConfig.Groups.AppendSysResourceIfExists(key, a.Sys); ok == true {
        a.resourcePrint(fileName, res)
	}

	// package
	if res, _, ok := gossConfig.Packages.AppendSysResourceIfExists(key, a.Sys); ok == true {
        a.resourcePrint(fileName, res)
	}

	// port
	if res, _, ok := gossConfig.Ports.AppendSysResourceIfExists(key, a.Sys); ok == true {
        a.resourcePrint(fileName, res)
	}

	// process
	if res, sysres, ok := gossConfig.Processes.AppendSysResourceIfExists(key, a.Sys); ok == true {
        a.resourcePrint(fileName, res)
		ports := system.GetPorts(true)
		pids, _ := sysres.Pids()
		for _, pid := range pids {
			pidS := strconv.Itoa(pid)
			for port, entries := range ports {
				for _, entry := range entries {
					if entry.Pid == pidS {
						// port
						if res, _, ok := gossConfig.Ports.AppendSysResourceIfExists(port, a.Sys); ok == true {
                            a.resourcePrint(fileName, res)
						}
					}
				}
			}
		}
	}

	// Service
	if res, _, ok := gossConfig.Services.AppendSysResourceIfExists(key, a.Sys); ok == true {
        a.resourcePrint(fileName, res)
	}

	// user
	if res, _, ok := gossConfig.Users.AppendSysResourceIfExists(key, a.Sys); ok == true {
        a.resourcePrint(fileName, res)
	}

	return nil
}

func (a *Add) resourcePrint(fileName string, res resource.ResourceRead) {
    resMap := map[string]resource.ResourceRead{res.ID(): res}

    oj, _ := marshal(resMap)
    typ := reflect.TypeOf(res)
    typs := strings.Split(typ.String(), ".")[1]

	fmt.Fprintf(a.Writer, "Adding %s to '%s':\n\n%s\n\n", typs, fileName, string(oj))
}