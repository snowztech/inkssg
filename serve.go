package inkssg

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

type ServeOptions struct {
	Addr string
}

func Serve(dir string, opts ...ServeOptions) error {
	o := ServeOptions{Addr: ":3000"}
	if len(opts) > 0 {
		if opts[0].Addr != "" {
			o.Addr = opts[0].Addr
		}
	}

	if dir == "" {
		dir = "."
	}
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return fmt.Errorf("invalid path: %w", err)
	}

	site, err := NewSite(absDir)
	if err != nil {
		return err
	}
	if err := site.Build(); err != nil {
		return err
	}

	outAbs, err := filepath.Abs(filepath.Join(absDir, site.Output))
	if err != nil {
		return fmt.Errorf("resolve output: %w", err)
	}

	rl := newReloader()

	rebuild := func() error {
		start := time.Now()
		s, err := NewSite(absDir)
		if err != nil {
			return err
		}
		if err := s.Build(); err != nil {
			return err
		}
		fmt.Printf("↻ rebuilt in %s\n", time.Since(start).Round(time.Millisecond))
		rl.broadcast()
		return nil
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	go func() {
		if err := watch(ctx, absDir, outAbs, rebuild); err != nil {
			fmt.Fprintf(os.Stderr, "watcher stopped: %v\n", err)
		}
	}()

	displayAddr := o.Addr
	if host, port, err := net.SplitHostPort(displayAddr); err == nil && host == "" {
		displayAddr = net.JoinHostPort("localhost", port)
	}
	fmt.Printf("inkssg serving %s on http://%s\n", absDir, displayAddr)
	return runServer(ctx, o.Addr, filepath.Join(absDir, site.Output), rl)
}
