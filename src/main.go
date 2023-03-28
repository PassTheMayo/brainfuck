package main

import (
	"context"
	"errors"
	"io"
	"log"
	"os"
	"time"

	"github.com/jessevdk/go-flags"
)

var (
	opts *options = &options{
		MemorySize: 30000,
		Timeout:    -1,
	}
)

type options struct {
	MemorySize  int  `short:"M" long:"memory-size" description:"The size of the memory strip in bytes"`
	Timeout     int  `short:"T" long:"timeout" description:"The timeout in seconds, or -1 for no timeout"`
	PrintCycles bool `short:"C" long:"print-cycles" description:"Print the amount of cycles when the program exits"`
	Debug       bool `short:"D" long:"debug" description:"Print the list of commands to the console before running"`
	NoRun       bool `short:"N" long:"no-run" description:"Do not run the program"`
	Silent      bool `short:"s" long:"silent" description:"Does not print any program output"`
}

func main() {
	args, err := flags.Parse(opts)

	if err != nil {
		if flags.WroteHelp(err) {
			os.Exit(0)
		}

		log.Fatal(err)
	}

	if len(args) < 1 {
		log.Fatalf("missing input file argument")
	}

	p := NewProgram(opts.MemorySize, OutOfBoundsError)

	if opts.Silent {
		p.Output = io.Discard
	}

	if err := p.ReadFile(args[0]); err != nil {
		log.Fatal(err)
	}

	if opts.Debug {
		p.DebugCommands()
	}

	if !opts.NoRun {
		var ctx context.Context
		var cancel context.CancelFunc

		if opts.Timeout != -1 {
			ctx, cancel = context.WithTimeout(context.Background(), time.Second*time.Duration(opts.Timeout))
		} else {
			ctx, cancel = context.WithCancel(context.Background())
		}

		defer cancel()

		if err := p.Run(ctx); err != nil && !errors.Is(err, context.DeadlineExceeded) {
			log.Println(err)
		}
	}

	if opts.PrintCycles {
		log.Printf("Cycles: %d\n", p.CycleCount())
	}
}
