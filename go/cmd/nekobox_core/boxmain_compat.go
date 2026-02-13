package main

import (
	"context"
	"encoding/json"
	"os"
	"os/signal"
	"syscall"

	box "github.com/sagernet/sing-box"
	"github.com/sagernet/sing-box/log"
	"github.com/sagernet/sing-box/option"
	E "github.com/sagernet/sing/common/exceptions"
	"github.com/sagernet/sing/service"
	"github.com/spf13/cobra"
)

// Compatibility layer for standard sing-box to work with NekoBox
// This replaces the MatsuriDayo-specific cmd/sing-box implementations

var (
	disableColor bool
	nekoCtx      context.Context
)

func init() {
	// Initialize context with default registry
	nekoCtx = service.ContextWithDefaultRegistry(context.Background())
}

// SetDisableColor sets whether to disable colored log output
func SetDisableColor(dc bool) {
	disableColor = dc
}

// Create creates a new sing-box instance from JSON configuration
func Create(nekoConfigContent []byte) (*box.Box, context.CancelFunc, error) {
	var options option.Options
	err := json.Unmarshal(nekoConfigContent, &options)
	if err != nil {
		return nil, nil, err
	}

	// Apply disable color setting
	if disableColor {
		if options.Log == nil {
			options.Log = &option.LogOptions{}
		}
		options.Log.DisableColor = true
	}

	ctx, cancel := context.WithCancel(nekoCtx)
	instance, err := box.New(box.Options{
		Context: ctx,
		Options: options,
	})
	if err != nil {
		cancel()
		return nil, nil, E.Cause(err, "create service")
	}

	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)
	defer func() {
		signal.Stop(osSignals)
		close(osSignals)
	}()

	startCtx, finishStart := context.WithCancel(context.Background())
	go func() {
		_, loaded := <-osSignals
		if loaded {
			cancel()
			closeMonitor(startCtx)
		}
	}()

	err = instance.Start()
	finishStart()
	if err != nil {
		cancel()
		return nil, nil, E.Cause(err, "start service")
	}

	return instance, cancel, nil
}

// closeMonitor is called when the instance is being shut down
func closeMonitor(ctx context.Context) {
	select {
	case <-ctx.Done():
		return
	default:
	}
}

// Main provides the standard sing-box CLI entry point
func Main() {
	if err := mainCommand.Execute(); err != nil {
		log.Fatal(err)
	}
}

var mainCommand = &cobra.Command{
	Use: "sing-box",
}
