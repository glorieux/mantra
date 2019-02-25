package main

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/spf13/afero"
)

func TestMain(m *testing.M) {
	fs = afero.NewMemMapFs()
	os.Exit(m.Run())
}

func TestCreateApplication(t *testing.T) {
	t.Run("No error", func(t *testing.T) {
		if err := createApplication("test"); err != nil {
			t.Error(err)
		}
	})

	t.Run("Directory exists", func(t *testing.T) {
		if err := createApplication("test"); err == nil {
			t.Error("Folder already exists")
		}
	})

	t.Run("Directories and files", func(t *testing.T) {
		appName := "test_file_structure"
		if err := createApplication(appName); err != nil {
			t.Error(err)
		}
		files := []string{
			appName,
			path.Join(appName, "cmd"),
			path.Join(appName, "cmd", "main.go"),
			path.Join(appName, fmt.Sprintf("%s.go", appName)),
		}

		for _, path := range files {
			if exists, _ := afero.Exists(fs, path); !exists {
				t.Errorf("%s doesn't exist", path)
			}
		}
	})
}

func TestCreateService(t *testing.T) {
	if err := createService("test"); err != nil {
		t.Error(err)
	}
	if exists, _ := afero.Exists(fs, "test.go"); !exists {
		t.Error("test.go doesn't exist")
	}
}
