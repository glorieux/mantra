package main

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path"
	"strings"

	"github.com/spf13/afero"
	"pkg.glorieux.io/mantra"
)

var fs = afero.NewOsFs()

var templates = template.New("").Funcs(template.FuncMap{
	"Title": func(s string) string {
		return strings.Title(s)
	},
})

func init() {
	template.Must(templates.New("go.mod").Parse(`module {{ .Name }}

require (
	pkg.glorieux.io/mantra v{{ .MantraVersion }}
)
`))

	template.Must(templates.New("main").Parse(`package main

import (
	"github.com/sirupsen/logrus"

	"pkg.glorieux.io/mantra"
)

func main() {
  log := logrus.New()
  mantra.New(log)
}
`))

	template.Must(templates.New("service").Parse(`package {{ . }}

import (
  "pkg.glorieux.io/mantra"
)

type {{ . | Title }}Service struct{}

// Serve runs the service
func (*{{ . | Title }}Service) Serve(mux mantra.ServeMux) {
	mux.Handle("tmp", func(e mantra.Event) {
		// TODO: Handle event
	})
}

// Stop stops the service
func (*{{ . | Title }}Service) Stop() error {
	return nil
}

func (*{{ . | Title }}Service) String() string {
	return "{{ . }}"
}
`))
}

func createApplication(name string) error {
	fmt.Printf("Creating application %s...\n", name)
	if exists, _ := afero.DirExists(fs, name); exists {
		return fmt.Errorf("Directory named %s already exists", name)
	}

	err := fs.MkdirAll(name, os.ModeDir|os.ModePerm)
	if err != nil {
		return fmt.Errorf("Error creating %s directory: %s", name, err)
	}

	if exists, _ := afero.Exists(fs, path.Join(name, "main.go")); !exists {
		var b bytes.Buffer

		err := templates.ExecuteTemplate(&b, "main", strings.ToLower(name))
		if err != nil {
			return err
		}

		err = afero.WriteFile(fs, path.Join(name, "main.go"), b.Bytes(), 0644)
		if err != nil {
			return err
		}
	}

	if exists, _ := afero.Exists(fs, path.Join(name, "go.mod")); !exists {
		var b bytes.Buffer

		err := templates.ExecuteTemplate(&b, "go.mod", struct {
			Name          string
			MantraVersion string
		}{
			strings.ToLower(name),
			mantra.VERSION,
		})
		if err != nil {
			return err
		}

		err = afero.WriteFile(fs, path.Join(name, "go.mod"), b.Bytes(), 0644)
		if err != nil {
			return err
		}
	}
	fmt.Println("You're good to go!")
	fmt.Printf("You can now run [go build ./%s] to build the application.\n", name)
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
