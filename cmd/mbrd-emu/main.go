package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/murych/si-medlab-mbrd-emu/internal/srvint"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	flags := flag.NewFlagSet("mbrd-emu", flag.ContinueOnError)
	flags.SetOutput(os.Stderr)
	device := flags.String("device", "", "serial device path")
	logPath := flags.String("log-file", "", "optional JSONL state log path")
	if err := flags.Parse(os.Args[1:]); err != nil {
		return err
	}
	if *device == "" {
		return errors.New("--device is required")
	}

	var logFile *os.File
	if *logPath != "" {
		file, err := os.OpenFile(*logPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o644)
		if err != nil {
			return fmt.Errorf("open log file %q: %w", *logPath, err)
		}
		logFile = file
		defer logFile.Close()
	}

	server, err := srvint.New(*device)
	if err != nil {
		return err
	}
	defer server.Close()
	if err := server.Connect(); err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "mbrd-emu ready on %s\n", *device)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	if err := server.Serve(ctx); err != nil {
		log.Printf("SrvInt runtime failure: %v", err)
		return err
	}
	return nil
}
