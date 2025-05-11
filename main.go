package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/jianongHe/term-rex/game"
	"github.com/nsf/termbox-go"
)

// Version information - this should match version.js
const (
	Version   = "0.1.5"
	BuildDate = "2025-05-11"
)

func main() {
	// Check for version flag
	if len(os.Args) > 1 && (os.Args[1] == "--version" || os.Args[1] == "-v") {
		fmt.Printf("Term-Rex v%s (built on %s)\n", Version, BuildDate)
		os.Exit(0)
	}

	// Set up signal handler to catch Ctrl+C
	setupSignalHandler()

	// Initialize terminal
	if err := termbox.Init(); err != nil {
		fmt.Printf("Failed to initialize terminal: %v\n", err)
		os.Exit(1)
	}
	defer termbox.Close()

	// Set game width based on terminal size
	//w, _ := termbox.Size()
	//game.SetWidth(w)

	// Create a new game instance
	g := game.NewGame()

	// Run the game
	g.Run()
}

// setupSignalHandler sets up a signal handler to catch Ctrl+C
func setupSignalHandler() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		termbox.Close()
		os.Exit(0)
	}()
}
