package testutil

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// Note: This file is copied verbatim from SPD golang lib

// LoadTestSQLFile load test sql data from a file
func LoadTestSQLFile(t *testing.T, tx *gorm.DB, filename string) {
	body, err := ReadFile(filename)
	require.NoError(t, err)

	err = tx.Exec(string(body)).Error
	require.NoError(t, err)
}

// ReadFile reads a file completely. But if the file does not exist, try to find it in the parent directory, [repeat...]
func ReadFile(filename string) ([]byte, error) {
	return os.ReadFile(abs(filename))
}

// abs returns absolute path of a file in project directory
// on errors, we return the requested `filename` for simplicity
// since eventually using that file (read/write) would error in a more convenient place to handle
// i.e. those read/write methods has `error` included in return value
func abs(filename string) string {
	cwd, err := os.Getwd()
	if err != nil {
		return filename
	}

	return _abs(cwd, filename)
}

// the reason we split into `_abs` is because the `filename`
// could be a relative path, e.g. `config/secret.yml` and we want to keep the `config/` hierarchy while traversing
// the parent directory of `dirname` recursively
func _abs(dirname, filename string) string {
	fullpath, err := filepath.Abs(filepath.Join(dirname, filename))
	if err != nil {
		return filename
	}

	if _, err = os.Stat(fullpath); err != nil {
		parentdir := filepath.Dir(dirname)
		if parentdir == "/" {
			return filename
		}
		return _abs(parentdir, filename)
	}

	return fullpath
}
