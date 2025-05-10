package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/jianongHe/term-rex/game"
	"github.com/nsf/termbox-go"
)

func main() {
	// 设置信号处理，捕获 Ctrl+C
	setupSignalHandler()

	// initialize terminal
	if err := termbox.Init(); err != nil {
		panic(err)
	}
	defer termbox.Close()

	// dynamic width based on terminal size
	//w, _ := termbox.Size()
	//game.SetWidth(w)

	// start game
	g := game.NewGame()
	g.Run()
}

// setupSignalHandler 设置信号处理器，捕获 Ctrl+C 信号
func setupSignalHandler() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		termbox.Close()
		os.Exit(0)
	}()
}
