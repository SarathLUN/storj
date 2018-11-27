// Copyright (C) 2018 Storj Labs, Inc.
// See LICENSE for copying information.

package satelliteweb

import (
	"net/http"
	"path/filepath"

	"go.uber.org/zap"

	"github.com/graphql-go/graphql"
)

// GatewayConfig contains configuration for gateway
type GatewayConfig struct {
	Address    string `help:"server address of the graphql api gateway and frontend app" default:"127.0.0.1:8081"`
	StaticPath string `help:"path to static resources" default:""`
}

type gateway struct {
	schema graphql.Schema
	config GatewayConfig
	log    *zap.Logger
}

func (gw *gateway) run() {
	mux := http.NewServeMux()
	//gw.config.StaticPath = "./web/satellite"
	fs := http.FileServer(http.Dir(gw.config.StaticPath))

	mux.Handle("/api/graphql/v0", http.HandlerFunc(gw.grapqlHandler))

	if gw.config.StaticPath != "" {
		mux.Handle("/", http.HandlerFunc(gw.appHandler))
		mux.Handle("/static/", http.StripPrefix("/static", fs))
	}

	err := http.ListenAndServe(gw.config.Address, mux)
	gw.log.Error("unexpected exit of satellite gateway server: ", zap.Error(err))
}

func (gw *gateway) appHandler(w http.ResponseWriter, req *http.Request) {
	http.ServeFile(w, req, filepath.Join(gw.config.StaticPath, "dist", "public", "index.html"))
}
