package goss

import (
	"fmt"
	"github.com/SimonBaeumer/goss/internal/app"
	"github.com/SimonBaeumer/goss/system"
	"github.com/SimonBaeumer/goss/util"
	"os"
	"strconv"
	"strings"
)

// AddResources is a sSimple wrapper to add multiple resources
func AddResources(fileName, resourceName string, keys []string, ctx app.CliContext) error {
	OutStoreFormat = getStoreFormatFromFileName(fileName)
	header := extractHeaderArgument(ctx.Header)

	config := util.Config{
		IgnoreList:        ctx.ExcludeAttr,
		Timeout:           ctx.Timeout,
		AllowInsecure:     ctx.AllowInsecure,
		NoFollowRedirects: ctx.NoFollowRedirects,
		Server:            ctx.Server,
		Username:          ctx.Username,
		Password:          ctx.Password,
		Header:            header,
	}

	var gossConfig GossConfig
	if _, err := os.Stat(fileName); err == nil {
		gossConfig = ReadJSON(fileName)
	} else {
		gossConfig = *NewGossConfig()
	}

	sys := system.New()

	for _, key := range keys {
		if err := AddResource(fileName, gossConfig, resourceName, key, config, sys); err != nil {
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
func AddResource(fileName string, gossConfig GossConfig, resourceName, key string, config util.Config, sys *system.System) error {
	// Need to figure out a good way to refactor this
	switch resourceName {
	case "Addr":
		res, err := gossConfig.Addrs.AppendSysResource(key, sys, config)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		resourcePrint(fileName, res)
	case "Command":
		res, err := gossConfig.Commands.AppendSysResource(key, sys, config)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		resourcePrint(fileName, res)
	case "DNS":
		res, err := gossConfig.DNS.AppendSysResource(key, sys, config)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		resourcePrint(fileName, res)
	case "File":
		res, err := gossConfig.Files.AppendSysResource(key, sys, config)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		resourcePrint(fileName, res)
	case "Group":
		res, err := gossConfig.Groups.AppendSysResource(key, sys, config)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		resourcePrint(fileName, res)
	case "Package":
		res, err := gossConfig.Packages.AppendSysResource(key, sys, config)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		resourcePrint(fileName, res)
	case "Port":
		res, err := gossConfig.Ports.AppendSysResource(key, sys, config)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		resourcePrint(fileName, res)
	case "Process":
		res, err := gossConfig.Processes.AppendSysResource(key, sys, config)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		resourcePrint(fileName, res)
	case "Service":
		res, err := gossConfig.Services.AppendSysResource(key, sys, config)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		resourcePrint(fileName, res)
	case "User":
		res, err := gossConfig.Users.AppendSysResource(key, sys, config)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		resourcePrint(fileName, res)
	case "Gossfile":
		res, err := gossConfig.Gossfiles.AppendSysResource(key, sys, config)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		resourcePrint(fileName, res)
	case "KernelParam":
		res, err := gossConfig.KernelParams.AppendSysResource(key, sys, config)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		resourcePrint(fileName, res)
	case "Mount":
		res, err := gossConfig.Mounts.AppendSysResource(key, sys, config)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		resourcePrint(fileName, res)
	case "Interface":
		res, err := gossConfig.Interfaces.AppendSysResource(key, sys, config)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		resourcePrint(fileName, res)
	case "HTTP":
		res, err := gossConfig.HTTPs.AppendSysResource(key, sys, config)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		resourcePrint(fileName, res)
	default:
		panic("Undefined resource name: " + resourceName)
	}

	return nil
}

// Simple wrapper to add multiple resources
func AutoAddResources(fileName string, keys []string, ctx app.CliContext) error {
	OutStoreFormat = getStoreFormatFromFileName(fileName)

	var gossConfig GossConfig
	if _, err := os.Stat(fileName); err == nil {
		gossConfig = ReadJSON(fileName)
	} else {
		gossConfig = *NewGossConfig()
	}

	sys := system.New()

	for _, key := range keys {
		if err := AutoAddResource(fileName, gossConfig, key, sys); err != nil {
			return err
		}
	}
	if err := WriteJSON(fileName, gossConfig); err != nil {
		return err
	}

	return nil
}

// Autoadds all resources to the config file
func AutoAddResource(fileName string, gossConfig GossConfig, key string, sys *system.System) error {
	// file
	if strings.Contains(key, "/") {
		if res, _, ok := gossConfig.Files.AppendSysResourceIfExists(key, sys); ok == true {
			resourcePrint(fileName, res)
		}
	}

	// group
	if res, _, ok := gossConfig.Groups.AppendSysResourceIfExists(key, sys); ok == true {
		resourcePrint(fileName, res)
	}

	// package
	if res, _, ok := gossConfig.Packages.AppendSysResourceIfExists(key, sys); ok == true {
		resourcePrint(fileName, res)
	}

	// port
	if res, _, ok := gossConfig.Ports.AppendSysResourceIfExists(key, sys); ok == true {
		resourcePrint(fileName, res)
	}

	// process
	if res, sysres, ok := gossConfig.Processes.AppendSysResourceIfExists(key, sys); ok == true {
		resourcePrint(fileName, res)
		ports := system.GetPorts(true)
		pids, _ := sysres.Pids()
		for _, pid := range pids {
			pidS := strconv.Itoa(pid)
			for port, entries := range ports {
				for _, entry := range entries {
					if entry.Pid == pidS {
						// port
						if res, _, ok := gossConfig.Ports.AppendSysResourceIfExists(port, sys); ok == true {
							resourcePrint(fileName, res)
						}
					}
				}
			}
		}
	}

	// Service
	if res, _, ok := gossConfig.Services.AppendSysResourceIfExists(key, sys); ok == true {
		resourcePrint(fileName, res)
	}

	// user
	if res, _, ok := gossConfig.Users.AppendSysResourceIfExists(key, sys); ok == true {
		resourcePrint(fileName, res)
	}

	return nil
}
