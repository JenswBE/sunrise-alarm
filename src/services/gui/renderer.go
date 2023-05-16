package gui

import (
	"embed"
	"errors"
	"fmt"
	"html/template"
	"path"
	"path/filepath"
	"strings"

	"github.com/JenswBE/sunrise-alarm/src/entities"
	"github.com/gin-contrib/multitemplate"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
)

//go:embed html
var htmlContent embed.FS

func (h *Handler) NewRenderer() multitemplate.Renderer {
	pages := map[string][]string{
		"clock":      {"pages/clock"},
		"debug":      {"pages/debug"},
		"alarmsList": {"pages/alarms_list"},
		"alarmsForm": {"pages/alarms_form"},
	}

	r := multitemplate.NewRenderer()
	for pageName, templates := range pages {
		// Create new template with functions
		templates = append([]string{"layouts/empty", "layouts/base"}, templates...)
		templatePaths := lo.Map(templates, func(i string, _ int) string { return fmt.Sprintf("html/%s.gohtml", i) })
		templateName := filepath.Base(templatePaths[len(templatePaths)-1]) // Last template is main template
		tmpl := template.New(templateName).Funcs(template.FuncMap{
			"add":          add,
			"contains":     lo.Contains[string],
			"formatDays":   formatDays,
			"getURL":       getURL,
			"getStaticURL": getStaticURL,
			"product":      product,
			"rawJS":        rawJS,
			"rawWeekday":   rawWeekday,
			"repeat":       lo.Range,
			"substract":    substract,
			"toLower":      strings.ToLower,
		})

		// Parse and add templates
		_, err := tmpl.ParseFS(htmlContent, templatePaths...)
		if err != nil {
			log.Fatal().Err(err).Strs("template_paths", templatePaths).Msg("Failed to parse template files")
		}

		// Add template to renderer
		r.Add(pageName, tmpl)
	}
	return r
}

func add(a, b int) int {
	return a + b
}

func substract(a, b int) int {
	return a - b
}

func product(a, b int) int {
	return a * b
}

func formatDays(days []entities.ISOWeekday) string {
	var builder strings.Builder
	builder.Grow(13) // 7 letters of days and 6 spaces
	for _, isoWeekday := range entities.ISOWeekdays() {
		if lo.Contains(days, isoWeekday) {
			builder.WriteByte(isoWeekday.String()[0]) // Add first letter of day
		} else {
			builder.WriteRune('_')
		}
		if isoWeekday != entities.ISOSunday {
			// Append space after each day, except for Sunday
			builder.WriteRune(' ')
		}
	}
	return builder.String()
}

func getURL(parts ...string) string {
	return path.Join(parts...) + "/"
}

func getStaticURL(parts ...string) string {
	if len(parts) == 0 || parts[0] == "" {
		log.Error().Stack().Err(errors.New("missing URL for static asset")).Msg("Missing URL for static asset")
	}
	parts = append([]string{"/static"}, parts...)
	return path.Join(parts...)
}

func rawJS(input string) template.JS {
	return template.JS(input)
}

func rawWeekday(input entities.ISOWeekday) int {
	return int(input)
}
