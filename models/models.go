package models

import "fmt"

// Model to design a database
type Model struct {
	ID          string
	Title       string
	Description string
}

// Field define a field in a model to embed data
type Field struct {
	ID          string
	Title       string
	Description string
	Type        string

	DefaultJSON string
	Required    bool
	Widget      string
	Searchable  string
}

var models map[string]interface{}

func initModels() {
	models = make(map[string]interface{})
	models["model"] = &Model{}
	models["field"] = &Field{}
}

// Get return the models Map
func Get() map[string]interface{} {
	if len(models) == 0 {
		initModels()
	}
	return models
}

func (m *Model) String() string {
	return fmt.Sprintf("Model-%s", m.Title)
}
func (f *Field) String() string {
	return fmt.Sprintf("Field-%s", f.Title)
}
