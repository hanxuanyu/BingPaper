package config

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	err := Init("")
	if err != nil {
		t.Fatalf("Failed to init config: %v", err)
	}

	cfg := GetConfig()
	if cfg.Server.Port != 8080 {
		t.Errorf("Expected port 8080, got %d", cfg.Server.Port)
	}

	if cfg.DB.Type != "sqlite" {
		t.Errorf("Expected DB type sqlite, got %s", cfg.DB.Type)
	}
}

func TestDebugFunctions(t *testing.T) {
	// 设置一个环境变量
	os.Setenv("BINGPAPER_SERVER_PORT", "9999")
	defer os.Unsetenv("BINGPAPER_SERVER_PORT")

	err := Init("")
	if err != nil {
		t.Fatalf("Failed to init config: %v", err)
	}

	settings := GetAllSettings()
	serverCfg, ok := settings["server"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected server config map, got %v", settings["server"])
	}

	// Viper numbers in AllSettings are often int
	portValue := serverCfg["port"]
	// 允许不同的数字类型，因为 viper 内部实现可能变化
	portStr := fmt.Sprintf("%v", portValue)
	if portStr != "9999" {
		t.Errorf("Expected port 9999 in settings, got %v (%T)", portValue, portValue)
	}

	overrides := GetEnvOverrides()
	found := false
	for _, o := range overrides {
		if strings.Contains(o, "server.port") && strings.Contains(o, "9999") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected server.port override in %v", overrides)
	}

	// 验证格式化输出
	formatted := GetFormattedSettings()
	if !strings.Contains(formatted, "server.port: 9999") {
		t.Errorf("Expected formatted settings to contain server.port: 9999, got %s", formatted)
	}
}
