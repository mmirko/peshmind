package peshmind

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
)

func (c *Config) SpawnPrologEngine(kbPool string) error {
	// Check for swipl binary existence, permissions, and path, try the execution with the --version flag
	cmd := exec.Command("swipl", "--version")
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to execute swipl: %w", err)
	}

	// Create the knowledge base
	cmd = exec.Command("./createkb.sh")
	cmd.Dir = kbPool
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to create Prolog engine: %w", err)
	}

	// Create the Prolog engine context and spawn the engine

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go func(ctx context.Context) {
		// Start the engine
		cmd := exec.Command("swipl", "data.pl")
		cmd.Dir = kbPool
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		err := cmd.Start()
		if err != nil {
			fmt.Printf("failed to start Prolog engine: %v\n", err)
			return
		}
		cmd.Wait()

	}(ctx)

	<-ctx.Done()

	return nil
}
