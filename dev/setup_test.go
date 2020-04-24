package dev

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"testing"

	// . "github.com/puma/puma-dev/dev/devtest"

	"github.com/stretchr/testify/assert"
)

func TestInstallIntoSystem(t *testing.T) {
	appLinkDir, _ := ioutil.TempDir("", ".puma-dev")
	libDir, _ := ioutil.TempDir("", "Library")
	logFilePath := filepath.Join(libDir, "Logs", "puma-dev.log")
	launchAgentDir := filepath.Join(libDir, "LaunchAgents")
	expectedPlistPath := filepath.Join(launchAgentDir, "io.puma.dev.plist")
	assert.NoDirExists(t, launchAgentDir)

	defer func() {
		exec.Command("launchctl", "unload", expectedPlistPath).Run()
		os.RemoveAll(appLinkDir)
		os.RemoveAll(logFilePath)
		os.RemoveAll(libDir)
	}()

	err := InstallIntoSystem(&InstallIntoSystemArgs{
		ListenPort:         10080,
		TlsPort:            10443,
		Domains:            "test:localhost",
		Timeout:            "5s",
		ApplinkDirPath:     appLinkDir,
		LaunchAgentDirPath: launchAgentDir,
		LogfilePath:        logFilePath,
	})

	assert.NoError(t, err)

	assert.DirExists(t, launchAgentDir)
	assert.FileExists(t, expectedPlistPath)
	AssertDirUmask(t, "0755", launchAgentDir)
}

func AssertDirUmask(t *testing.T, expectedUmask, path string) {
	info, err := os.Stat(path)
	if !assert.NoError(t, err) {
		assert.True(t, info.IsDir())

		actualUmask := "0" + strconv.FormatInt(int64(info.Mode().Perm()), 8)
		assert.Equal(t, expectedUmask, actualUmask)
	}
}
