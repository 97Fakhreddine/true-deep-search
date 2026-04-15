package browser

import (
	"fmt"
	"os/exec"
	"runtime"
)

func Open(target string) error {
	if target == "" {
		return fmt.Errorf("empty target")
	}

	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("xdg-open", target)
	case "darwin":
		cmd = exec.Command("open", target)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", target)
	default:
		return fmt.Errorf("unsupported platform")
	}

	return cmd.Start()
}
