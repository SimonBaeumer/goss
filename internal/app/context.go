package app

import (
    "github.com/urfave/cli"
    "time"
)

type CliContext struct {
    FormatOptions []string
    Sleep time.Duration
    RetryTimeout time.Duration
    MaxConcurrent int
    Package string
    Vars string
    Gossfile string
    ExcludeAttr []string
    Timeout int
    AllowInsecure bool
    NoFollowRedirects bool
    Server string
    Username string
    Password string
    Header string
    Endpoint string
    ListenAddr string
    Cache time.Duration
    Format string
    NoColor bool
    Color bool
}

func NewCliContext(c *cli.Context) CliContext {
    return CliContext{
        FormatOptions: c.StringSlice("format-options"),
        Sleep: c.Duration("sleep"),
        RetryTimeout: c.Duration("retry-timeout"),
        Package: c.String("package"),
        MaxConcurrent: c.Int("max-concurrent"),
        Vars: c.GlobalString("vars"),
        Gossfile: c.GlobalString("gossfile"),
        ExcludeAttr: c.GlobalStringSlice("exclude-attr"),
        Timeout: int(c.Duration("timeout") / time.Millisecond),
        AllowInsecure: c.Bool("insecure"),
        NoFollowRedirects: c.Bool("no-follow-redirects"),
        Server: c.String("server"),
        Username: c.String("username"),
        Password: c.String("password"),
        Header: c.String("header"),
        Endpoint: c.String("endpoint"),
        ListenAddr: c.String("listen-addr"),
        Cache: c.Duration("cache"),
        Format: c.String("format"),
        NoColor: c.Bool("no-color"),
        Color: c.Bool("color"),
    }
}
