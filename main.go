package main

import (
	dockerfile "github.com/moby/buildkit/frontend/dockerfile/builder"
	"github.com/moby/buildkit/frontend/gateway/grpcclient"
	"github.com/moby/buildkit/util/appcontext"
	"github.com/moby/buildkit/util/bklog"
)

func main() {
	//dockerfile.Build
	//reads a Dockerfile provided to the frontend
	if err := grpcclient.RunFromEnvironment(appcontext.Context(), dockerfile.Build); err != nil {
		bklog.L.Errorf("fatal error: %+v", err)
		panic(err)
	}
}
