package Utils

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

func FileExists(filePath string) bool {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	}
	return true
}

func GetDurationFlag(cmd *cobra.Command, flag string) (*time.Duration, error) {
	value, _ := cmd.Flags().GetString(flag)
	duration, err := time.ParseDuration(value)

	if err != nil {
		return nil, fmt.Errorf("'%s' is not a valid duration. Valid duration examples: 200ms, 5s, 10m, 2h, 1d", value)
	}

	return &duration, nil
}
