package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"syscall"
	"time"

	conf "github.com/ad/findmyticker/config"
	"github.com/getlantern/systray"
	"github.com/kardianos/osext"
)

var (
	cancel context.CancelFunc

	version = `0.0.1`

	menuInfo *systray.MenuItem

	config *conf.Config
)

func main() {
	cfg, errInitConfig := conf.InitConfig()
	if errInitConfig != nil {
		log.Fatalf("InitConfig: %s", errInitConfig.Error())

		return
	}

	config = cfg

	_, cancel = context.WithCancel(context.Background())

	defer func() {
		cancel()
		systray.Quit()
	}()

	go func() {
		run()

		for range time.Tick(1 * time.Minute) {
			run()
		}
	}()

	systray.Run(onReady, onExit)

}

func onReady() {
	systray.SetTitle("⚐")
	mTitle := systray.AddMenuItem(fmt.Sprintf("FindMy ticker v%s", version), "App title")
	mTitle.Disable()

	menuInfo = systray.AddMenuItem(fmt.Sprintf("started at: %s", time.Now().Format("15:04:05")), "")
	menuInfo.Disable()

	mRestart := systray.AddMenuItem("Restart", "Restart app")
	mQuit := systray.AddMenuItem("Quit", "Quit app")

	for {
		select {
		case <-mRestart.ClickedCh:
			fmt.Println("Requesting restart")
			cancel()
			_ = Restart()
			return
		case <-mQuit.ClickedCh:
			fmt.Println("Requesting quit")
			cancel()
			systray.Quit()
			return
		}
	}
}

func onExit() {
	// clean up here
}

// Restart app
func Restart() error {
	file, error := osext.Executable()
	if error != nil {
		return error
	}

	error = syscall.Exec(file, os.Args, os.Environ())
	if error != nil {
		return error
	}

	return nil
}
