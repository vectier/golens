package main

import (
	"bufio"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/vectier/golens/internal/lsp"
)

func main() {
	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt)

	log.SetOutput(os.Stderr)
	r, w := bufio.NewReader(os.Stdin), os.Stdout

	srv := lsp.NewServer()
	srv.Handle(ctx, r, w)
}
