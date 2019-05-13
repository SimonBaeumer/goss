package main

import (
	"fmt"
    app2 "github.com/SimonBaeumer/goss/internal/app"
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
	//"time"
)

var version string

func main() {
	startTime := time.Now()
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
		cli.StringFlag{
			Name:  "package",
			Usage: "Package type to use [rpm, deb, apk, pacman]",
		},
	}
	app.Commands = []cli.Command{
		{
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
                ctx := app2.NewCliContext(c)

				runtime := getGossRunTime(ctx)

				v := &goss.Validator{
					MaxConcurrent: ctx.MaxConcurrent,
					Package:       ctx.Package,
					Outputer:      outputs.GetOutputer(ctx.Format),
					FormatOptions: ctx.FormatOptions,
					GossConfig:    runtime.GetGossConfig(),
				}

				//TODO: ugly shit to set the color here, tmp fix for the moment!
                if ctx.NoColor {
                    color.NoColor = true
                }
                if ctx.Color {
                    color.NoColor = false
                }

				os.Exit(v.Validate(startTime))
				return nil
			},
		},
		{
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
				ctx := app2.NewCliContext(c)

				gossRunTime := getGossRunTime(ctx)

				h := &goss.HealthHandler{
					Cache:         cache.New(ctx.Cache, 30*time.Second),
					ListenAddr:    ctx.ListenAddr,
					Outputer:      outputs.GetOutputer(ctx.Format),
					Sys:           system.New(ctx.Package),
					GossMu:        &sync.Mutex{},
					MaxConcurrent: ctx.MaxConcurrent,
					GossConfig:    gossRunTime.GetGossConfig(),
				}

				if ctx.Format == "json" {
					h.ContentType = "application/json"
				}

				h.Serve(ctx.Endpoint)
				return nil
			},
		},
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
				fmt.Print(goss.RenderJSON(c))
				return nil
			},
		},
		{
			Name:    "autoadd",
			Aliases: []string{"aa"},
			Usage:   "automatically add all matching resource to the test suite",
			Action: func(c *cli.Context) error {
				ctx := app2.NewCliContext(c)
				goss.AutoAddResources(c.GlobalString("gossfile"), c.Args(), ctx)
				return nil
			},
		},
		{
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
					Name:  "package",
					Usage: "add new package",
					Action: func(c *cli.Context) error {
						ctx := app2.NewCliContext(c)
						goss.AddResources(c.GlobalString("gossfile"), "Package", c.Args(), ctx)
						return nil
					},
				},
				{
					Name:  "file",
					Usage: "add new file",
					Action: func(c *cli.Context) error {
						ctx := app2.NewCliContext(c)
						goss.AddResources(c.GlobalString("gossfile"), "File", c.Args(), ctx)
						return nil
					},
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
					Action: func(c *cli.Context) error {
						ctx := app2.NewCliContext(c)
						goss.AddResources(c.GlobalString("gossfile"), "Addr", c.Args(), ctx)
						return nil
					},
				},
				{
					Name:  "port",
					Usage: "add new listening [protocol]:port - ex: 80 or udp:123",
					Action: func(c *cli.Context) error {
						ctx := app2.NewCliContext(c)
						goss.AddResources(c.GlobalString("gossfile"), "Port", c.Args(), ctx)
						return nil
					},
				},
				{
					Name:  "service",
					Usage: "add new service",
					Action: func(c *cli.Context) error {
						ctx := app2.NewCliContext(c)
						goss.AddResources(c.GlobalString("gossfile"), "Service", c.Args(), ctx)
						return nil
					},
				},
				{
					Name:  "user",
					Usage: "add new user",
					Action: func(c *cli.Context) error {
						ctx := app2.NewCliContext(c)
						goss.AddResources(c.GlobalString("gossfile"), "User", c.Args(), ctx)
						return nil
					},
				},
				{
					Name:  "group",
					Usage: "add new group",
					Action: func(c *cli.Context) error {
						ctx := app2.NewCliContext(c)
						goss.AddResources(c.GlobalString("gossfile"), "Group", c.Args(), ctx)
						return nil
					},
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
					Action: func(c *cli.Context) error {
						ctx := app2.NewCliContext(c)
						goss.AddResources(c.GlobalString("gossfile"), "Command", c.Args(), ctx)
						return nil
					},
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
					Action: func(c *cli.Context) error {
						ctx := app2.NewCliContext(c)
						goss.AddResources(c.GlobalString("gossfile"), "DNS", c.Args(), ctx)
						return nil
					},
				},
				{
					Name:  "process",
					Usage: "add new process name",
					Action: func(c *cli.Context) error {
						ctx := app2.NewCliContext(c)
						goss.AddResources(c.GlobalString("gossfile"), "Process", c.Args(), ctx)
						return nil
					},
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
					Action: func(c *cli.Context) error {
						ctx := app2.NewCliContext(c)
						goss.AddResources(c.GlobalString("gossfile"), "HTTP", c.Args(), ctx)
						return nil
					},
				},
				{
					Name:  "goss",
					Usage: "add new goss file, it will be imported from this one",
					Action: func(c *cli.Context) error {
						ctx := app2.NewCliContext(c)
						goss.AddResources(c.GlobalString("gossfile"), "Gossfile", c.Args(), ctx)
						return nil
					},
				},
				{
					Name:  "kernel-param",
					Usage: "add new goss kernel param",
					Action: func(c *cli.Context) error {
						ctx := app2.NewCliContext(c)
						goss.AddResources(c.GlobalString("gossfile"), "KernelParam", c.Args(), ctx)
						return nil
					},
				},
				{
					Name:  "mount",
					Usage: "add new mount",
					Action: func(c *cli.Context) error {
						ctx := app2.NewCliContext(c)
						goss.AddResources(c.GlobalString("gossfile"), "Mount", c.Args(), ctx)
						return nil
					},
				},
				{
					Name:  "interface",
					Usage: "add new interface",
					Action: func(c *cli.Context) error {
						ctx := app2.NewCliContext(c)
						goss.AddResources(c.GlobalString("gossfile"), "Interface", c.Args(), ctx)
						return nil
					},
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func getGossRunTime(ctx app2.CliContext) goss.GossRunTime {
	runtime := goss.GossRunTime{
		Gossfile: ctx.Gossfile,
		Vars:     ctx.Vars,
		Package:  ctx.Package,
	}
	return runtime
}
