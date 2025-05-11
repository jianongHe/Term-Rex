package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/jianongHe/term-rex/game"
	"github.com/nsf/termbox-go"
)

// Version information
const (
	Version   = "1.0.0"
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
		panic(err)
	}
	defer termbox.Close()

	// Dynamic width based on terminal size
	//w, _ := termbox.Size()
	//game.SetWidth(w)

	// Start game
	g := game.NewGame()
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
