package envHelper

import (
	"crypto/rand"
	"encoding/hex"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/IceWhaleTech/CasaOS-AppManagement/pkg/config"
)

var (
	defaultPasswordOnce sync.Once
	defaultPassword     string
)

const defaultPasswordFile = "default-password"

// ensureDefaultPassword loads or generates the host-unique random secret used
// to substitute $DefaultPassword in app compose templates.
//
// The secret is persisted (0600) under the runtime path so it stays stable
// across restarts. It is never a well-known constant, which removes the
// "admin / casaos" default that every CasaOS install shared.
func ensureDefaultPassword() {
	runtimePath := config.CommonInfo.RuntimePath

	if data, err := os.ReadFile(filepath.Join(runtimePath, defaultPasswordFile)); err == nil {
		if password := strings.TrimSpace(string(data)); len(password) >= 20 {
			defaultPassword = password
			return
		}
	}

	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return
	}
	// hex string: 32 chars, safe for URLs, shell and docker env vars
	password := hex.EncodeToString(buf)

	if runtimePath != "" {
		if err := os.MkdirAll(runtimePath, 0o755); err == nil {
			if err := os.WriteFile(filepath.Join(runtimePath, defaultPasswordFile), []byte(password+"\n"), 0o600); err == nil {
				defaultPassword = password
				return
			}
		}
	}
	// Persistence failed — still use a process-unique random secret rather than
	// falling back to a well-known constant.
	defaultPassword = password
}

// GetDefaultPassword returns the host-unique random secret used for
// $DefaultPassword substitution in app compose templates.
func GetDefaultPassword() string {
	defaultPasswordOnce.Do(ensureDefaultPassword)
	return defaultPassword
}

func ReplaceDefaultENV(key, tz string) string {
	temp := ""
	switch key {
	case "$DefaultPassword":
		if password := GetDefaultPassword(); password != "" {
			temp = password
		}
	case "$DefaultUserName":
		temp = "admin"

	case "$PUID":
		temp = "1000"
	case "$PGID":
		temp = "1000"
	case "$TZ":
		temp = tz
	}
	return temp
}

// replace env default setting
func ReplaceStringDefaultENV(str string) string {
	return strings.ReplaceAll(strings.ReplaceAll(str, "$DefaultPassword", ReplaceDefaultENV("$DefaultPassword", "")), "$DefaultUserName", ReplaceDefaultENV("$DefaultUserName", ""))
}