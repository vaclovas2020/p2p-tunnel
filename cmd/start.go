package cmd

import (
	"context"
	"flag"
	"fmt"

	"github.com/google/subcommands"
	"github.com/vaclovas2020/p2p-tunnel/p2p"
)

type StartServerCmd struct {
	port int
}

func (*StartServerCmd) Name() string     { return "start" }
func (*StartServerCmd) Synopsis() string { return "start tunnel" }
func (*StartServerCmd) Usage() string {
	return `start [-port]:
	Start tunnel on specific port.
  `
}

func (p *StartServerCmd) SetFlags(f *flag.FlagSet) {
	f.IntVar(&p.port, "port", 7777, "tunnel port (default: 7777)")
}

func (p *StartServerCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...any) subcommands.ExitStatus {
	p2p.StartServer(p.port)
	fmt.Println()
	return subcommands.ExitSuccess
}
