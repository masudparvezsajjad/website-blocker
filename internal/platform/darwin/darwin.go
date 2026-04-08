package platform

import (
	"fmt"
	"os"
	"runtime"
)

func EnsureMacOS() error {
	if runtime.GOOS != "darwin" {
		return fmt.Errorf("this build is for macOS only")
	}
	return nil
}

func EnsureRoot() error {
	if os.Geteuid() != 0 {
		return fmt.Errorf("this command must be run with sudo")
	}
	return nil
}
