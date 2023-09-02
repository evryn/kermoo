package user_config_test

import (
	"kermoo/modules/logger"
	"kermoo/modules/user_config"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

const ROOT = "../../.."

const (
	PATH_VALID_YAML          = ROOT + "/tests/units/stubs/valid.yaml"
	PATH_VALID_JSON          = ROOT + "/tests/units/stubs/valid.json"
	PATH_INVALID_YAML        = ROOT + "/tests/units/stubs/invalid.yaml"
	PATH_INVALID_JSON        = ROOT + "/tests/units/stubs/invalid.json"
	PATH_INVALSTRUCTURE_YAML = ROOT + "/tests/units/stubs/invalid_structure.yaml"
	PATH_INVALSTRUCTURE_JSON = ROOT + "/tests/units/stubs/invalid_structure.json"
)

func TestMakeConfigFromFilename(t *testing.T) {
	logger.MustInitLogger("fatal")

	tt := []struct {
		name         string
		filename     string
		expectsError bool
	}{
		{
			name:     "valid json file",
			filename: PATH_VALID_JSON,
		},
		{
			name:     "valid yaml file",
			filename: PATH_VALID_YAML,
		},
		{
			name:         "invalid json file",
			filename:     PATH_INVALID_JSON,
			expectsError: true,
		},
		{
			name:         "invalid yaml file",
			filename:     PATH_INVALID_YAML,
			expectsError: true,
		},
		{
			name:         "non-existent file",
			filename:     "/path/to/non_existent.json",
			expectsError: true,
		},
		{
			name:         "parsable json but invalid structure",
			filename:     PATH_INVALSTRUCTURE_JSON,
			expectsError: true,
		},
		{
			name:         "parsable yaml but invalid structure",
			filename:     PATH_INVALSTRUCTURE_YAML,
			expectsError: true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			_, err := user_config.MakePreparedConfig(tc.filename)

			if tc.expectsError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestMakeConfigFromStdin(t *testing.T) {
	logger.MustInitLogger("fatal")

	tt := []struct {
		name         string
		stdin        string
		expectsError bool
	}{
		{
			name:  "valid json",
			stdin: getFileContent(t, PATH_VALID_JSON),
		},
		{
			name:  "valid yaml",
			stdin: getFileContent(t, PATH_VALID_YAML),
		},

		{
			name:         "invalid json",
			stdin:        getFileContent(t, PATH_INVALID_JSON),
			expectsError: true,
		},
		{
			name:         "invalid yaml",
			stdin:        getFileContent(t, PATH_INVALID_YAML),
			expectsError: true,
		},
		{
			name:         "parsable json with invalid structure",
			stdin:        getFileContent(t, PATH_INVALSTRUCTURE_JSON),
			expectsError: true,
		},
		{
			name:         "parsable yaml with invalid structure",
			stdin:        getFileContent(t, PATH_INVALSTRUCTURE_YAML),
			expectsError: true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if tc.stdin != "" {
				tmpfile, err := os.CreateTemp("", "stdin")
				if err != nil {
					t.Fatalf("failed to create temporary file: %v", err)
				}
				defer os.Remove(tmpfile.Name())

				tmpfile.WriteString(tc.stdin)
				tmpfile.Seek(0, 0)

				oldStdin := os.Stdin
				defer func() { os.Stdin = oldStdin }()
				os.Stdin = tmpfile
			}

			_, err := user_config.MakePreparedConfig("-")

			if tc.expectsError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestMakeConfigFromString(t *testing.T) {
	logger.MustInitLogger("fatal")

	tt := []struct {
		name         string
		content      string
		expectsError bool
	}{
		{
			name:    "valid json",
			content: getFileContent(t, PATH_VALID_JSON),
		},
		{
			name:    "valid yaml",
			content: getFileContent(t, PATH_VALID_YAML),
		},

		{
			name:         "invalid json",
			content:      getFileContent(t, PATH_INVALID_JSON),
			expectsError: true,
		},
		{
			name:         "invalid yaml",
			content:      getFileContent(t, PATH_INVALID_YAML),
			expectsError: true,
		},
		{
			name:         "parsable json with invalid structure",
			content:      getFileContent(t, PATH_INVALSTRUCTURE_JSON),
			expectsError: true,
		},
		{
			name:         "parsable yaml with invalid structure",
			content:      getFileContent(t, PATH_INVALSTRUCTURE_YAML),
			expectsError: true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			_, err := user_config.MakePreparedConfig(tc.content)

			if tc.expectsError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func getFileContent(t *testing.T, filename string) string {
	bytes, err := os.ReadFile(filename)

	if err != nil {
		t.Fatal("unable to load file from "+filename, err)
	}

	return string(bytes)
}
