package cmd

import (
	"context"
	"flag"
	"fmt"

	"github.com/google/subcommands"
	"github.com/vaclovas2020/p2p-tunnel/p2p"
)

type SendMessageCmd struct {
	port int
	host string
}

func (*SendMessageCmd) Name() string     { return "send" }
func (*SendMessageCmd) Synopsis() string { return "send message to peer" }
func (*SendMessageCmd) Usage() string {
	return `send [-port] [-host] <message>:
	Send message to peer.
  `
}

func (p *SendMessageCmd) SetFlags(f *flag.FlagSet) {
	f.IntVar(&p.port, "port", 7777, "tunnel port (default: 7777)")
	f.StringVar(&p.host, "host", "localhost", "peer host (default: localhost)")
}

func (p *SendMessageCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	f.Parse(f.Args())
	p2p.SendMessageToServer(p.host, p.port, f.Arg(0))
	fmt.Println()
	return subcommands.ExitSuccess
}
