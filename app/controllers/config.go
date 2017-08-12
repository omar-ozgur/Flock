package controllers

import (
	"html/template"
)

// Initialize templates
var templates = template.Must(template.ParseGlob("app/views/pages/*"))
