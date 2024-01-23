package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"syscall"
	"time"

	conf "github.com/ad/findmyticker/config"
	"github.com/ad/lru"
	"github.com/getlantern/systray"
	"github.com/kardianos/osext"
)

var (
	cancel context.CancelFunc

	version = `0.1.2`

	menuInfo  *systray.MenuItem
	menuError *systray.MenuItem

	config *conf.Config

	lruCache *lru.Cache[string, [2]float64]
)

func main() {
	cfg, errInitConfig := conf.InitConfig()
	if errInitConfig != nil {
		log.Fatalf("InitConfig: %s", errInitConfig.Error())

		return
	}

	config = cfg

	if config.OpenFindMyOnStartup {
		runFindMyApp()
	}

	lruCache = lru.New[string, [2]float64]()

	_, cancel = context.WithCancel(context.Background())

	defer func() {
		cancel()
		systray.Quit()
	}()

	go func() {
		run()

		for range time.Tick(time.Duration(config.Period) * time.Second) {
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

	menuError = systray.AddMenuItem(fmt.Sprintf("error: %s", "none"), "")
	menuError.Disable()
	menuError.Hide()

	mConfig := systray.AddMenuItem("Config", "Open config")
	mRestart := systray.AddMenuItem("Restart", "Restart app")
	mQuit := systray.AddMenuItem("Quit", "Quit app")

	for {
		select {
		case <-mRestart.ClickedCh:
			fmt.Println("Requesting restart")
			cancel()
			_ = Restart()
			return
		case <-mConfig.ClickedCh:
			fmt.Println("Opening config editor")
			_ = OpenConfig()
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

// Open config
func OpenConfig() error {
	_ = conf.OpenConfigEditor()

	return nil
}

func runFindMyApp() error {
	cmd := exec.Command(`open`, "--hide", "--background", "/System/Applications/FindMy.app")
	stderr, err := cmd.StderrPipe()
	// log.SetOutput(os.Stderr)

	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	slurp, _ := io.ReadAll(stderr)
	fmt.Printf("%s\n", slurp)

	if err := cmd.Wait(); err != nil {
		return err
	}

	return nil
}
