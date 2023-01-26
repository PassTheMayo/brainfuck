package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func main() {
	p := NewProgram(30000, OutOfBoundsError)

	if err := p.ReadFile(os.Args[1]); err != nil {
		log.Fatal(err)
	}

	p.DebugCommands()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)

	defer cancel()

	start := time.Now()

	if err := p.Run(ctx); err != nil {
		log.Println(err)
	}

	fmt.Printf("----------\nCYCLES: %s (%s)\n", message.NewPrinter(language.English).Sprintf("%d", p.CycleCount()), time.Since(start))
}
