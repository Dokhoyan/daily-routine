package config

import (
	"fmt"
	"os"
)

const (
	testModeEnabledEnvName = "ENABLE_TEST_MODE"
	testUserIDEnvName      = "TEST_USER_ID"
)

type testConfig struct {
	enabled bool
	userID  int64
}

func NewTestConfig() TestConfig {
	enabled := os.Getenv(testModeEnabledEnvName) == "true"
	userIDStr := os.Getenv(testUserIDEnvName)
	userID := parseInt64(userIDStr)
	if userID == 0 {
		userID = 123456789
	}

	return &testConfig{
		enabled: enabled,
		userID:  userID,
	}
}

func (cfg *testConfig) IsTestModeEnabled() bool {
	return cfg.enabled
}

func (cfg *testConfig) GetTestUserID() int64 {
	return cfg.userID
}

func parseInt64(s string) int64 {
	var result int64
	fmt.Sscanf(s, "%d", &result)
	return result
}
