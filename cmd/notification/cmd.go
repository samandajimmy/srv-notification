package main

import (
	"flag"
	"fmt"
	"os"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/ncore"
)

type CmdFlags struct {
	CmdShowHelp    *bool
	CmdShowVersion *bool
	OptWorkDir     *string
}

type BootOptions struct {
	Core ncore.BootOptions
}

/// initCmdFlags initiate available command line interface commands and options for parsing
func initCmdFlags() CmdFlags {
	return CmdFlags{
		CmdShowHelp:    flag.Bool("help", false, "Command: Show available commands and options"),
		CmdShowVersion: flag.Bool("version", false, "Command: Show version"),
		OptWorkDir:     flag.String("dir", ".", "Option: Set working directory"),
	}
}

func handleCmdFlags() BootOptions {
	// Parse CLI commands and options
	cmdFlags := initCmdFlags()
	flag.Parse()

	// Intercept help command
	if *cmdFlags.CmdShowHelp {
		printHelp()
		os.Exit(0)
	}

	// Intercept version command
	if *cmdFlags.CmdShowVersion {
		printVersion()
		os.Exit(0)
	}

	return BootOptions{
		Core: ncore.BootOptions{
			Manifest: ncore.Manifest{
				AppName:    AppName,
				AppVersion: AppVersion,
				Metadata: map[string]interface{}{
					"build_hash": BuildHash,
				},
			},
			WorkDir: *cmdFlags.OptWorkDir,
		},
	}
}

func printHelp() {
	fmt.Printf("%s. Available Commands and Options:\n\n", AppName)
	flag.PrintDefaults()
}

// printVersion print app version and integrity
func printVersion() {
	fmt.Printf("%s\n"+
		"  Version    : %s\n"+
		"  Build Hash : %s\n",
		AppName, AppVersion, BuildHash)
}
