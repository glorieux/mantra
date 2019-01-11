package main

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path"
	"strings"

	"github.com/spf13/afero"
)

var fs = afero.NewOsFs()

var templates = template.New("").Funcs(template.FuncMap{
	"Title": func(s string) string {
		return strings.Title(s)
	},
})

func init() {
	template.Must(templates.New("go.mod").Parse(`module {{ . }}`))

	template.Must(templates.New("main").Parse(`package main

import (
	"github.com/sirupsen/logrus"

	"glorieux.io/mantra"

  "{{ . }}"
)

func main() {
  log := logrus.New()
  err := mantra.New({{ . }}.Application{}, log)
  if err != nil {
    log.Fatal(err)
  }
}
`))

	template.Must(templates.New("application").Parse(`package {{ . }}

import "glorieux.io/mantra"

// Application is a mantra application
type Application struct {
	send mantra.SendFunc
}

// Init initializes the application
func (app Application) Init(send mantra.SendFunc) error {
	app.send = send
	return nil
}

// String returns the application name
func (Application) String() string {
	return "{{ . }}"
}
`))

	template.Must(templates.New("service").Parse(`package {{ . }}

import (
  "fmt"

  "glorieux.io/mantra"
)

type {{ . | Title }}Service struct {}

func (s *Service) HandleMessage(msg mantra.Message) error {
	switch msg.(type) {
		default:
			return fmt.Errof("Unknown message: %+v", msg)
	}
}

func (*Service) String() string {
	return "{{ . }}"
}
`))
}

// TODO: Create go.mod
func createApplication(name string) error {
	if exists, _ := afero.DirExists(fs, name); exists {
		return fmt.Errorf("Directory named %s already exists", name)
	}

	cmdDirPath := path.Join(name, "cmd", name)
	err := fs.MkdirAll(cmdDirPath, os.ModeDir|os.ModePerm)
	if err != nil {
		return fmt.Errorf("Creating cmd directory: %s", err)
	}

	if exists, _ := afero.Exists(fs, path.Join(cmdDirPath, "main.go")); !exists {
		var b bytes.Buffer

		err := templates.ExecuteTemplate(&b, "main", strings.ToLower(name))
		if err != nil {
			return err
		}

		err = afero.WriteFile(fs, path.Join(cmdDirPath, "main.go"), b.Bytes(), 0644)
		if err != nil {
			return err
		}
	}

	if exists, _ := afero.Exists(fs, path.Join(name, fmt.Sprintf("%s.go", name))); !exists {
		var b bytes.Buffer

		err := templates.ExecuteTemplate(&b, "application", strings.ToLower(name))
		if err != nil {
			return err
		}

		err = afero.WriteFile(fs, path.Join(name, fmt.Sprintf("%s.go", name)), b.Bytes(), 0644)
		if err != nil {
			return err
		}
	}

	if exists, _ := afero.Exists(fs, path.Join(name, "go.mod")); !exists {
		var b bytes.Buffer

		err := templates.ExecuteTemplate(&b, "go.mod", strings.ToLower(name))
		if err != nil {
			return err
		}

		err = afero.WriteFile(fs, path.Join(name, "go.mod"), b.Bytes(), 0644)
		if err != nil {
			return err
		}
	}
	return nil
}

func createService(name string) error {
	if exists, _ := afero.Exists(fs, fmt.Sprintf("%s.go", name)); !exists {
		var b bytes.Buffer

		err := templates.ExecuteTemplate(&b, "service", strings.ToLower(name))
		if err != nil {
			return err
		}

		err = afero.WriteFile(fs, fmt.Sprintf("%s.go", name), b.Bytes(), 0644)
		if err != nil {
			return err
		}
	}
	return nil
}
