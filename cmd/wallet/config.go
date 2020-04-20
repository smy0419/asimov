// Copyright (c) 2018-2020. The asimov developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package main

import (
	"github.com/jessevdk/go-flags"
	"os"
)

const (
	defaultFormat = "hex"
)

// config defines the configuration options for findcheckpoint.
//
// See loadConfig for details on the configuration load process.
type config struct {
	Help       bool `short:"h" long:"help" description:"Show usage."`
	Cmd        string `short:"c" long:"cmd" description:"Command: genKey, genMultiSigAddress"`
	Format     string `short:"f" long:"format" description:"in/out format, currently, support hex, base64"`
	Net        string `short:"n" long:"net" description:"support main,dev,test,regtest,default is main"`
}


// loadConfig initializes and parses the config using command line options.
func loadConfig() (*config, []string, error) {
	// Default config.
	cfg := config{
		Format: defaultFormat,
	}

	// Parse command line options.
	parser := flags.NewParser(&cfg, flags.Default)
	remainingArgs, err := parser.Parse()
	if err != nil {
		if e, ok := err.(*flags.Error); !ok || e.Type != flags.ErrHelp {
			parser.WriteHelp(os.Stderr)
		}
		return nil, nil, err
	}

	// Validate database type.
	if cfg.Help {
		parser.WriteHelp(os.Stderr)
		return nil, nil, err
	}
	if len(cfg.Cmd) == 0 {
		parser.WriteHelp(os.Stderr)
		return nil, nil, err
	}

	return &cfg, remainingArgs, nil
}