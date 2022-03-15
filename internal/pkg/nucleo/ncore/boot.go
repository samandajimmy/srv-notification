package ncore

import (
	"github.com/google/uuid"
	"path"
)

type BootOptions struct {
	Manifest        Manifest
	NodeId          string
	WorkDir         string
	ResponseMapFile string
}

func Boot(args ...BootOptions) *Core {
	// Load Options
	options := getBootOptions(args)

	// Init Core
	core := Core{
		Manifest: options.Manifest,
		WorkDir:  options.WorkDir,
		NodeId:   options.NodeId,
	}

	return &core
}

func getBootOptions(args []BootOptions) BootOptions {
	// Get options
	var options BootOptions
	if len(args) > 0 {
		options = args[0]
	}

	// If working directory is not set, then set to current directory
	if options.WorkDir == "" {
		options.WorkDir = "."
	}

	// If config file is not set, then set default
	if options.ResponseMapFile == "" {
		options.ResponseMapFile = path.Join(options.WorkDir, "responses.yml")
	}

	// If node number is not set, then generate a random uuid
	if options.NodeId == "" {
		nodeId, err := uuid.NewUUID()
		if err != nil {
			panic(err)
		}
		options.NodeId = nodeId.String()
	}

	return options
}
