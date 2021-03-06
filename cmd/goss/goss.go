package main

import (
	"fmt"
	"github.com/SimonBaeumer/goss/system"
	"github.com/fatih/color"
	"github.com/patrickmn/go-cache"
	"log"
	"os"
	"sync"
	"time"

	"github.com/SimonBaeumer/goss"
	"github.com/SimonBaeumer/goss/outputs"
	"github.com/urfave/cli"
)

var version string

func main() {
	app := createApp()

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

// CliConfig represents all configurations passed by the cli app
type CliConfig struct {
	FormatOptions     []string
	Sleep             time.Duration
	RetryTimeout      time.Duration
	MaxConcurrent     int
	Vars              string
	Gossfile          string
	ExcludeAttr       []string
	Timeout           int
	AllowInsecure     bool
	NoFollowRedirects bool
	Server            string
	Username          string
	Password          string
	Header            string
	Endpoint          string
	ListenAddr        string
	Cache             time.Duration
	Format            string
	NoColor           bool
	Color             bool
}

// NewCliConfig creates an object from the cli.Context. It is used as a constructor and will do simple type conversions
// and validations
func NewCliConfig(c *cli.Context) CliConfig {
	return CliConfig{
		FormatOptions:     c.StringSlice("format-options"),
		Sleep:             c.Duration("sleep"),
		RetryTimeout:      c.Duration("retry-timeout"),
		MaxConcurrent:     c.Int("max-concurrent"),
		Vars:              c.GlobalString("vars"),
		Gossfile:          c.GlobalString("gossfile"),
		ExcludeAttr:       c.GlobalStringSlice("exclude-attr"),
		Timeout:           int(c.Duration("timeout") / time.Millisecond),
		AllowInsecure:     c.Bool("insecure"),
		NoFollowRedirects: c.Bool("no-follow-redirects"),
		Server:            c.String("server"),
		Username:          c.String("username"),
		Password:          c.String("password"),
		Header:            c.String("header"),
		Endpoint:          c.String("endpoint"),
		ListenAddr:        c.String("listen-addr"),
		Cache:             c.Duration("cache"),
		Format:            c.String("format"),
		NoColor:           c.Bool("no-color"),
		Color:             c.Bool("color"),
	}
}

// Create the cli.App object
func createApp() *cli.App {
	app := cli.NewApp()
	app.EnableBashCompletion = true
	app.Version = version
	app.Name = "goss"
	app.Usage = "Quick and Easy server validation"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "gossfile, g",
			Value:  "./goss.yaml",
			Usage:  "GossRunTime file to read from / write to",
			EnvVar: "GOSS_FILE",
		},
		cli.StringFlag{
			Name:   "vars",
			Usage:  "json/yaml file containing variables for template",
			EnvVar: "GOSS_VARS",
		},
	}
	app.Commands = []cli.Command{
		createValidateCommand(app),
		createServeCommand(),
		createAddCommand(app),
		{
			Name:    "render",
			Aliases: []string{"r"},
			Usage:   "render gossfile after imports",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "debug, d",
					Usage: fmt.Sprintf("Print debugging info when rendering"),
				},
			},
			Action: func(c *cli.Context) error {
				conf := NewCliConfig(c)
				runtime := goss.GossRunTime{
					Vars:     conf.Vars,
					Gossfile: conf.Gossfile,
				}

				fmt.Print(runtime.Render())
				return nil
			},
		},
		{
			Name:    "autoadd",
			Aliases: []string{"aa"},
			Usage:   "automatically add all matching resource to the test suite",
			Action: func(c *cli.Context) error {
				conf := NewCliConfig(c)

				a := goss.Add{
					Writer:            app.Writer,
					Sys:               system.New(),
					Username:          conf.Username,
					Password:          conf.Password,
					Server:            conf.Server,
					Timeout:           conf.Timeout,
					NoFollowRedirects: conf.NoFollowRedirects,
					AllowInsecure:     conf.AllowInsecure,
					ExcludeAttr:       conf.ExcludeAttr,
					Header:            conf.Header,
				}
				return a.AutoAddResources(c.GlobalString("gossfile"), c.Args())
			},
		},
	}
	return app
}

func createServeCommand() cli.Command {
	return cli.Command{
		Name:    "serve",
		Aliases: []string{"s"},
		Usage:   "Serve a health endpoint",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:   "format, f",
				Value:  "rspecish",
				Usage:  fmt.Sprintf("Format to output in, valid options: %s", outputs.Outputers()),
				EnvVar: "GOSS_FMT",
			},
			cli.StringSliceFlag{
				Name:   "format-options, o",
				Usage:  fmt.Sprintf("Extra options passed to the formatter, valid options: %s", outputs.FormatOptions()),
				EnvVar: "GOSS_FMT_OPTIONS",
			},
			cli.DurationFlag{
				Name:   "cache,c",
				Usage:  "Time to cache the results",
				Value:  5 * time.Second,
				EnvVar: "GOSS_CACHE",
			},
			cli.StringFlag{
				Name:   "listen-addr,l",
				Value:  ":8080",
				Usage:  "Address to listen on [ip]:port",
				EnvVar: "GOSS_LISTEN",
			},
			cli.StringFlag{
				Name:   "endpoint,e",
				Value:  "/healthz",
				Usage:  "Endpoint to expose",
				EnvVar: "GOSS_ENDPOINT",
			},
			cli.IntFlag{
				Name:   "max-concurrent",
				Usage:  "Max number of tests to run concurrently",
				Value:  50,
				EnvVar: "GOSS_MAX_CONCURRENT",
			},
		},
		Action: func(c *cli.Context) error {
			conf := NewCliConfig(c)

			gossRunTime := goss.GossRunTime{
				Gossfile: c.GlobalString("gossfile"),
				Vars:     c.GlobalString("vars"),
			}

			h := &goss.HealthHandler{
				Cache:         cache.New(conf.Cache, 30*time.Second),
				Outputer:      outputs.GetOutputer(conf.Format),
				Sys:           system.New(),
				GossMu:        &sync.Mutex{},
				MaxConcurrent: conf.MaxConcurrent,
				GossConfig:    gossRunTime.GetGossConfig(),
				FormatOptions: conf.FormatOptions,
			}

			if conf.Format == "json" {
				h.ContentType = "application/json"
			}

			gossRunTime.Serve(conf.Endpoint, h)
			return nil
		},
	}
}

func createAddHandler(app *cli.App, resourceName string) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		conf := NewCliConfig(c)

		a := goss.Add{
			Sys:               system.New(),
			Writer:            app.Writer,
			Username:          conf.Username,
			Password:          conf.Password,
			Server:            conf.Server,
			Timeout:           conf.Timeout,
			NoFollowRedirects: conf.NoFollowRedirects,
			AllowInsecure:     conf.AllowInsecure,
			ExcludeAttr:       conf.ExcludeAttr,
			Header:            conf.Header,
		}

		a.AddResources(c.GlobalString("gossfile"), resourceName, c.Args())
		return nil
	}
}

func createAddCommand(app *cli.App) cli.Command {
	return cli.Command{
		Name:    "add",
		Aliases: []string{"a"},
		Usage:   "add a resource to the test suite",
		Flags: []cli.Flag{
			cli.StringSliceFlag{
				Name:  "exclude-attr",
				Usage: "Exclude the following attributes when adding a new resource",
			},
		},
		Subcommands: []cli.Command{
			{
				Name:   "Package",
				Usage:  "add new package",
				Action: createAddHandler(app, "package"),
			},
			{
				Name:   "file",
				Usage:  "add new file",
				Action: createAddHandler(app, "file"),
			},
			{
				Name:  "addr",
				Usage: "add new remote address:port - ex: google.com:80",
				Flags: []cli.Flag{
					cli.DurationFlag{
						Name:  "timeout",
						Value: 500 * time.Millisecond,
					},
				},
				Action: createAddHandler(app, "addr"),
			},
			{
				Name:   "port",
				Usage:  "add new listening [protocol]:port - ex: 80 or udp:123",
				Action: createAddHandler(app, "port"),
			},
			{
				Name:   "service",
				Usage:  "add new service",
				Action: createAddHandler(app, "service"),
			},
			{
				Name:   "user",
				Usage:  "add new user",
				Action: createAddHandler(app, "user"),
			},
			{
				Name:   "group",
				Usage:  "add new group",
				Action: createAddHandler(app, "group"),
			},
			{
				Name:  "command",
				Usage: "add new command",
				Flags: []cli.Flag{
					cli.DurationFlag{
						Name:  "timeout",
						Value: 10 * time.Second,
					},
				},
				Action: createAddHandler(app, "command"),
			},
			{
				Name:  "dns",
				Usage: "add new dns",
				Flags: []cli.Flag{
					cli.DurationFlag{
						Name:  "timeout",
						Value: 500 * time.Millisecond,
					},
					cli.StringFlag{
						Name:  "server",
						Usage: "The IP address of a DNS server to query",
					},
				},
				Action: createAddHandler(app, "dns"),
			},
			{
				Name:   "process",
				Usage:  "add new process name",
				Action: createAddHandler(app, "process"),
			},
			{
				Name:  "http",
				Usage: "add new http",
				Flags: []cli.Flag{
					cli.BoolFlag{
						Name: "insecure, k",
					},
					cli.BoolFlag{
						Name: "no-follow-redirects, r",
					},
					cli.DurationFlag{
						Name:  "timeout",
						Value: 5 * time.Second,
					},
					cli.StringFlag{
						Name:  "username, u",
						Usage: "Username for basic auth",
					},
					cli.StringFlag{
						Name:  "password, p",
						Usage: "Password for basic auth",
					},
					cli.StringFlag{
						Name:  "header",
						Usage: "Set-Cookie: Value",
					},
				},
				Action: createAddHandler(app, "http"),
			},
			{
				Name:   "goss",
				Usage:  "add new goss file, it will be imported from this one",
				Action: createAddHandler(app, "goss"),
			},
			{
				Name:   "kernel-param",
				Usage:  "add new goss kernel param",
				Action: createAddHandler(app, "kernel-param"),
			},
			{
				Name:   "mount",
				Usage:  "add new mount",
				Action: createAddHandler(app, "mount"),
			},
			{
				Name:   "interface",
				Usage:  "add new interface",
				Action: createAddHandler(app, "interface"),
			},
		},
	}
}

func createValidateCommand(app *cli.App) cli.Command {
	startTime := time.Now()

	return cli.Command{
		Name:    "validate",
		Aliases: []string{"v"},
		Usage:   "Validate system",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:   "format, f",
				Value:  "rspecish",
				Usage:  fmt.Sprintf("Format to output in, valid options: %s", outputs.Outputers()),
				EnvVar: "GOSS_FMT",
			},
			cli.StringSliceFlag{
				Name:   "format-options, o",
				Usage:  fmt.Sprintf("Extra options passed to the formatter, valid options: %s", outputs.FormatOptions()),
				EnvVar: "GOSS_FMT_OPTIONS",
			},
			cli.BoolFlag{
				Name:   "color",
				Usage:  "Force color on",
				EnvVar: "GOSS_COLOR",
			},
			cli.BoolFlag{
				Name:   "no-color",
				Usage:  "Force color off",
				EnvVar: "GOSS_NOCOLOR",
			},
			cli.DurationFlag{
				Name:   "sleep,s",
				Usage:  "Time to sleep between retries, only active when -r is set",
				Value:  1 * time.Second,
				EnvVar: "GOSS_SLEEP",
			},
			cli.DurationFlag{
				Name:   "retry-timeout,r",
				Usage:  "Retry on failure so long as elapsed + sleep time is less than this",
				Value:  0,
				EnvVar: "GOSS_RETRY_TIMEOUT",
			},
			cli.IntFlag{
				Name:   "max-concurrent",
				Usage:  "Max number of tests to run concurrently",
				Value:  50,
				EnvVar: "GOSS_MAX_CONCURRENT",
			},
		},
		Action: func(c *cli.Context) error {
			conf := NewCliConfig(c)

			runtime := getGossRunTime(conf.Gossfile, conf.Vars)

			v := &goss.Validator{
				OutputWriter:  app.Writer,
				MaxConcurrent: conf.MaxConcurrent,
				Outputer:      outputs.GetOutputer(conf.Format),
				FormatOptions: conf.FormatOptions,
				GossConfig:    runtime.GetGossConfig(),
			}

			//TODO: ugly shit to set the color here, tmp fix for the moment!
			if conf.NoColor {
				color.NoColor = true
			}
			if conf.Color {
				color.NoColor = false
			}

			if conf.Gossfile == "testing" {
				v.Validate(startTime)
				return nil
			}

			os.Exit(runtime.Validate(v))
			return nil
		},
	}
}

func getGossRunTime(gossfile string, vars string) goss.GossRunTime {
	runtime := goss.GossRunTime{
		Gossfile: gossfile,
		Vars:     vars,
	}
	return runtime
}
