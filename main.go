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

	version = `0.1.3`

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
		errRunFindMyApp := runFindMyApp()
		if errRunFindMyApp != nil {
			fmt.Printf("runFindMyApp: %s\n", errRunFindMyApp.Error())
		}

	}

	if config.BringFindMyToFrontOnIdle {
		go bringFindMyToFrontOnIdle()
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

func bringFindMyToFrontOnIdle() error {
	appleScript := `repeat
	set num to (do shell script "ioreg -c IOHIDSystem | awk '/HIDIdleTime/ {print $NF/1000000000; exit}'")
	
	set o to (offset of "." in num)
	if ((o > 0) and (0.0 as text is "0,0")) then set num to (text 1 thru (o - 1) of num & "," & text (o + 1) thru -1 of num)
	set idleTime to num as integer

	if idleTime is greater than or equal to (1 * 10) then
		log "is idle"
		tell application "System Events"
			tell process "Find My"
				set frontmost to true
			end tell
			delay 16
			tell process "Finder"
				set frontmost to true
			end tell
			delay 16
		end tell
	end if
	delay 1
end repeat`

	cmd := exec.Command("/usr/bin/osascript", "-e", appleScript)
	stderr, err := cmd.StderrPipe()

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
