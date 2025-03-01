package main

import (
	"context"
	"flag"
	"os"

	"github.com/google/subcommands"
	"github.com/vaclovas2020/p2p-tunnel/cmd"
)

func main() {
	subcommands.Register(subcommands.HelpCommand(), "common")
	subcommands.Register(subcommands.FlagsCommand(), "common")
	subcommands.Register(subcommands.CommandsCommand(), "common")
	subcommands.Register(&cmd.RunTunnel{}, "p2p-client")

	flag.Parse()
	ctx := context.Background()
	os.Exit(int(subcommands.Execute(ctx)))
}
