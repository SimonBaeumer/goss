package goss

import (
	"bytes"
	"github.com/SimonBaeumer/goss/internal/app"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/SimonBaeumer/goss/outputs"
	"github.com/SimonBaeumer/goss/system"
	"github.com/SimonBaeumer/goss/util"
	"github.com/fatih/color"
	"github.com/patrickmn/go-cache"
)

//TODO: Maybe seperating handler and server?
type HealthHandler struct {
	RunTimeConfig GossRunTime
	C             app.CliContext
	GossConfig    GossConfig
	Sys           *system.System
	Outputer      outputs.Outputer
	Cache         *cache.Cache
	GossMu        *sync.Mutex
	ContentType   string
	MaxConcurrent int
	ListenAddr    string
}

func (h *HealthHandler) Serve(endpoint string) {
	color.NoColor = true

	http.Handle(endpoint, h)
	log.Printf("Starting to listen on: %s", h.ListenAddr)
	log.Fatal(http.ListenAndServe(h.ListenAddr, nil))
}

type res struct {
	exitCode int
	b        bytes.Buffer
}

//ServeHTTP fulfills the handler interface and is called as a handler on the
//health check request.
func (h HealthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	outputConfig := util.OutputConfig{
		FormatOptions: h.C.FormatOptions,
	}

	log.Printf("%v: requesting health probe", r.RemoteAddr)
	var resp res
	tmp, found := h.Cache.Get("res")
	if found {
		resp = tmp.(res)
	} else {
		h.GossMu.Lock()
		defer h.GossMu.Unlock()
		tmp, found := h.Cache.Get("res")
		if found {
			resp = tmp.(res)
		} else {
			h.Sys = system.New()
			log.Printf("%v: Stale Cache, running tests", r.RemoteAddr)
			iStartTime := time.Now()
			out := validate(h.Sys, h.GossConfig, h.MaxConcurrent)
			var b bytes.Buffer
			exitCode := h.Outputer.Output(&b, out, iStartTime, outputConfig)
			resp = res{exitCode: exitCode, b: b}
			h.Cache.Set("res", resp, cache.DefaultExpiration)
		}
	}
	if h.ContentType != "" {
		w.Header().Set("Content-Type", h.ContentType)
	}
	if resp.exitCode == 0 {
		resp.b.WriteTo(w)
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
		resp.b.WriteTo(w)
	}
}
