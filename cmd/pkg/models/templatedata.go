package models

// TemplateData holds data sent to template from handler
type TemplateData struct {
	StringMap map[string]string
	IntMap    map[string]int
	FloatMap  map[string]float32
	Data      map[string]interface{}
	CSFRToken string
	Flash     string
	Warning   string
	Error     string
}
