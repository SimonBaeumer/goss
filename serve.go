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

func Serve(ctx app.CliContext) {
	endpoint := ctx.Endpoint
	color.NoColor = true
	cache := cache.New(ctx.Cache, 30*time.Second)

	health := healthHandler{
		c:             ctx,
		gossConfig:    getGossConfig(ctx),
		sys:           system.New(ctx.Package),
		outputer:      getOutputer(ctx),
		cache:         cache,
		gossMu:        &sync.Mutex{},
		maxConcurrent: ctx.MaxConcurrent,
	}
	if ctx.Format == "json" {
		health.contentType = "application/json"
	}
	http.Handle(endpoint, health)
	listenAddr := ctx.ListenAddr
	log.Printf("Starting to listen on: %s", listenAddr)
	log.Fatal(http.ListenAndServe(ctx.ListenAddr, nil))
}

type res struct {
	exitCode int
	b        bytes.Buffer
}
type healthHandler struct {
	c             app.CliContext
	gossConfig    GossConfig
	sys           *system.System
	outputer      outputs.Outputer
	cache         *cache.Cache
	gossMu        *sync.Mutex
	contentType   string
	maxConcurrent int
}

func (h healthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	outputConfig := util.OutputConfig{
		FormatOptions: h.c.FormatOptions,
	}

	log.Printf("%v: requesting health probe", r.RemoteAddr)
	var resp res
	tmp, found := h.cache.Get("res")
	if found {
		resp = tmp.(res)
	} else {
		h.gossMu.Lock()
		defer h.gossMu.Unlock()
		tmp, found := h.cache.Get("res")
		if found {
			resp = tmp.(res)
		} else {
			h.sys = system.New(h.c.Package)
			log.Printf("%v: Stale cache, running tests", r.RemoteAddr)
			iStartTime := time.Now()
			out := validate(h.sys, h.gossConfig, h.maxConcurrent)
			var b bytes.Buffer
			exitCode := h.outputer.Output(&b, out, iStartTime, outputConfig)
			resp = res{exitCode: exitCode, b: b}
			h.cache.Set("res", resp, cache.DefaultExpiration)
		}
	}
	if h.contentType != "" {
		w.Header().Set("Content-Type", h.contentType)
	}
	if resp.exitCode == 0 {
		resp.b.WriteTo(w)
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
		resp.b.WriteTo(w)
	}
}
