package template

import (
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// Engine struct
type Engine struct {
	// delimiters
	left  string
	right string
	// views folder
	directory string
	// fs.FS supports embedded files
	fileSystem fs.FS
	// views extension
	extension string
	// layout variable name that incapsulates the template
	layout string
	// determines if the engine parsed all templates
	loaded bool
	// reload on each render
	reload bool
	// debug prints the parsed templates
	debug bool
	// lock for funcmap and templates
	mutex sync.RWMutex
	// template funcmap
	funcmap map[string]any
	// templates
	Templates *template.Template
}

// New returns a HTML render engine
func New(directory, extension string) *Engine {
	fileSystem := os.DirFS(directory)
	return NewFileSystem(fileSystem, extension)
}

// NewFileSystem returns a HTML render engine
func NewFileSystem(fileSystem fs.FS, extension string) *Engine {
	engine := &Engine{
		left:       "{{",
		right:      "}}",
		directory:  ".",
		fileSystem: fileSystem,
		extension:  extension,
		layout:     "embed",
		funcmap:    make(map[string]any),
	}
	engine.AddFunc(engine.layout, func() error {
		return fmt.Errorf("layout called unexpectedly.")
	})
	return engine
}

// Layout defines the variable name that will incapsulate the template
func (e *Engine) Layout(key string) *Engine {
	e.layout = key
	return e
}

// Delims sets the action delimiters to the specified strings, to be used in
// templates. An empty delimiter stands for the
// corresponding default: {{ or }}.
func (e *Engine) Delims(left, right string) *Engine {
	e.left, e.right = left, right
	return e
}

// AddFunc adds the function to the template's function map.
// It is legal to overwrite elements of the default actions
func (e *Engine) AddFunc(name string, fn any) *Engine {
	e.mutex.Lock()
	e.funcmap[name] = fn
	e.mutex.Unlock()
	return e
}

// AddFuncs adds the functions to the template's function map.
// It is legal to overwrite elements of the default actions
func (e *Engine) AddFuncs(fm template.FuncMap) *Engine {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	for k, v := range fm {
		e.funcmap[k] = v
	}
	return e
}

// Reload if set to true the templates are reloading on each render,
// use it when you're in development and you don't want to restart
// the application when you edit a template file.
func (e *Engine) Reload(enabled bool) *Engine {
	e.reload = enabled
	return e
}

// Debug will print the parsed templates when Load is triggered.
func (e *Engine) Debug(enabled bool) *Engine {
	e.debug = enabled
	return e
}

// Load parses the templates to the engine.
func (e *Engine) Load() error {
	if e.loaded {
		return nil
	}
	// race safe
	e.mutex.Lock()
	defer e.mutex.Unlock()
	e.Templates = template.New(e.directory)

	// Set template settings
	e.Templates.Delims(e.left, e.right)
	e.Templates.Funcs(e.funcmap)

	walkFn := func(path string, info fs.DirEntry, err error) error {
		// Return error if exist
		if err != nil {
			return err
		}
		// Skip file if it's a directory or has no file info
		if info == nil || info.IsDir() {
			return nil
		}
		// Skip file if it does not equal the given template extension
		if len(e.extension) >= len(path) || path[len(path)-len(e.extension):] != e.extension {
			return nil
		}
		// Get the relative file path
		// ./views/html/index.tmpl -> index.tmpl
		rel, err := filepath.Rel(e.directory, path)
		if err != nil {
			return err
		}
		// Reverse slashes '\' -> '/' and
		// partials\footer.tmpl -> partials/footer.tmpl
		name := filepath.ToSlash(rel)
		// Remove ext from name 'index.tmpl' -> 'index'
		name = strings.TrimSuffix(name, e.extension)
		// name = strings.Replace(name, e.extension, "", -1)
		// Read the file
		// #gosec G304
		buf, err := readFile(path, e.fileSystem)
		if err != nil {
			return err
		}
		// Create new template associated with the current one
		// This enable use to invoke other templates {{ template .. }}
		_, err = e.Templates.New(name).Parse(string(buf))
		if err != nil {
			return err
		}
		// Debugging
		if e.debug {
			fmt.Printf("views: parsed template: %s\n", name)
		}
		return err
	}
	// notify engine that we parsed all templates
	e.loaded = true
	return fs.WalkDir(e.fileSystem, e.directory, walkFn)
}

// Render will execute the template name along with the given values.
func (e *Engine) Render(out io.Writer, template string, binding any, layout ...string) error {
	if !e.loaded || e.reload {
		if e.reload {
			e.loaded = false
		}
		if err := e.Load(); err != nil {
			return err
		}
	}

	tmpl := e.Templates.Lookup(template)
	if tmpl == nil {
		return fmt.Errorf("render: template %s does not exist", template)
	}
	if len(layout) > 0 && layout[0] != "" {
		lay := e.Templates.Lookup(layout[0])
		if lay == nil {
			return fmt.Errorf("render: layout %s does not exist", layout[0])
		}
		e.mutex.Lock()
		defer e.mutex.Unlock()
		lay.Funcs(map[string]any{
			e.layout: func() error {
				return tmpl.Execute(out, binding)
			},
		})
		return lay.Execute(out, binding)
	}
	return tmpl.Execute(out, binding)
}

func readFile(path string, fileSystem fs.FS) ([]byte, error) {
	f, err := fileSystem.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return io.ReadAll(f)
}
