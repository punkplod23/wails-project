package main

import (
	"context"
	"fmt"

	"github.com/punkplod23/wails-project/internal/parsecsv"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	parser := parsecsv.NewCSVParser(a.ctx)
	results := parser.Query(name)
	return fmt.Sprintf("%s", results)
}

func (a *App) SelectFile() string {
	file, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{})
	if err != nil {
		return err.Error()
	}
	return file
}

func (a *App) RunCSV(filename string) string {
	parser := parsecsv.NewCSVParser(a.ctx)
	results := parser.RunFile(filename)
	fmt.Println(results)
	return fmt.Sprintf("Hello %s, It's show time!", results)
}
