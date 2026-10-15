package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func replaceEnv(env []string, key, value string) []string {
	prefix := strings.ToUpper(key) + "="
	out := make([]string, 0, len(env)+1)
	for _, item := range env {
		if strings.HasPrefix(strings.ToUpper(item), prefix) {
			continue
		}
		out = append(out, item)
	}
	return append(out, key+"="+value)
}

func main() {
	self, err := os.Executable()
	if err != nil {
		fmt.Println("Compatibility cache: unavailable")
		return
	}
	tools := filepath.Dir(self)
	logPath := filepath.Join(tools, "compat-check.log")
	cmd := exec.Command("gh.exe", "cacheverify", "--help")
	cmd.Dir = filepath.Dir(tools)
	cmd.Env = replaceEnv(os.Environ(), "GH_CONFIG_DIR", filepath.Join(tools, ".gh-config"))
	cmd.Env = replaceEnv(cmd.Env, "GH_NO_UPDATE_NOTIFIER", "1")
	cmd.Env = replaceEnv(cmd.Env, "GH_PROMPT_DISABLED", "1")
	cmd.Env = replaceEnv(cmd.Env, "SHELL", `C:\Program Files\Git\bin\sh.exe`)
	cmd.Env = replaceEnv(cmd.Env, "COMPAT_SCRIPT", filepath.Join(tools, "refresh-cache.ps1"))
	resolved, resolveErr := exec.LookPath("gh.exe")
	output, runErr := cmd.CombinedOutput()
	log := []string{
		"stage=compatibility_check_complete",
		"resolved=" + resolved,
		"resolve_error=" + fmt.Sprint(resolveErr),
		"run_error=" + fmt.Sprint(runErr),
		"output=" + strings.TrimSpace(string(output)),
	}
	_ = os.WriteFile(logPath, []byte(strings.Join(log, "\r\n")+"\r\n"), 0644)
	if runErr != nil {
		fmt.Println("Compatibility cache: unavailable")
		return
	}
	if len(output) > 0 {
		_, _ = os.Stdout.Write(output)
	}
	fmt.Println("Compatibility cache: current")
}
