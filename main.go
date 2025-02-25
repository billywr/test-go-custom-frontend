package main

import (
	"context"
	"log"

	"github.com/moby/buildkit/client/llb"
	"github.com/moby/buildkit/frontend/gateway/client"
	"github.com/moby/buildkit/frontend/gateway/grpcclient"
)

func build(ctx context.Context, c client.Client) (*client.Result, error) {
	// Start the custom BuildKit frontend
	log.Println("Starting custom BuildKit frontend...")

	// Define the Windows base image
	baseImage := "mcr.microsoft.com/windows/nanoserver:ltsc2022"
	st := llb.Image(baseImage)

	// Run a simple Windows command in the container
	cmd := "cmd /c echo Hello from WCOW Custom Frontend"
	exec := st.Run(llb.Shlex(cmd))

	// Get the final LLB state
	finalState := exec.Root()

	// Marshal the LLB state to a definition
	def, err := finalState.Marshal(ctx)
	if err != nil {
		log.Fatalf("Failed to marshal LLB state: %v", err)
		return nil, err
	}

	// Convert LLB definition to protocol buffer
	pbDef := def.ToPB()

	// Solve the LLB definition using BuildKit client
	log.Println("Solving the LLB definition...")
	ref, err := c.Solve(ctx, client.SolveRequest{Definition: pbDef})
	if err != nil {
		log.Fatalf("Failed to solve LLB definition: %v", err)
		return nil, err
	}

	// Create and return the build result
	log.Println("Returning successful build result...")
	res := &client.Result{}
	res.SetRef(ref.Ref)

	return res, nil
}

func main() {
	ctx := context.Background()

	// No need to explicitly set the address; use the environment variable
	log.Println("Connecting to BuildKit using environment settings...")

	// Run the custom frontend using BuildKit environment connection
	if err := grpcclient.RunFromEnvironment(ctx, build); err != nil {
		log.Fatalf("Failed to run frontend: %v", err)
	}
}
