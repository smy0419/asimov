// Copyright (c) 2018-2020 The asimov developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.
package util

import (
	"github.com/AsimovNetwork/asimov/rpcs/node"
	"os"
	"os/signal"
	"syscall"
)

func StartNode(stack *node.Node) {
	if err := stack.Start(); err != nil {
		rpcLog.Errorf("Error starting protocol stack: %v", err)
		os.Exit(1)
	}
	go func() {
		sigc := make(chan os.Signal, 1)
		signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
		defer signal.Stop(sigc)
		<-sigc
		rpcLog.Info("Got interrupt, shutting down...")
		go stack.Stop()
		for i := 10; i > 0; i-- {
			<-sigc
			if i > 1 {
				rpcLog.Warn("Already shutting down, interrupt more to panic.", "times", i-1)
			}
		}
	}()
}
