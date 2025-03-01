package cmd

import (
	"context"
	"flag"
	"fmt"

	"github.com/google/subcommands"
)

type RunTunnel struct {
	port int
}

func (*RunTunnel) Name() string     { return "start" }
func (*RunTunnel) Synopsis() string { return "start tunnel" }
func (*RunTunnel) Usage() string {
	return `start [-port]:
	Start tunnel on specific port.
  `
}

func (p *RunTunnel) SetFlags(f *flag.FlagSet) {
	f.IntVar(&p.port, "port", 7777, "tunnel port (default 7777)")
}

func (p *RunTunnel) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	fmt.Printf("Starting P2P tunnel on port %d...", p.port)
	return subcommands.ExitSuccess
}
