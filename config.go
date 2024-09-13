package main

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/log"
	yaml "gopkg.in/yaml.v3"
)

func (c *config) readConfig(file string) error {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		log.Errorf("Error finding configuration file in current directory: %v", err)
		var find bool
		form := huh.NewConfirm().
			Title("Configuration file not found").
			Description("Would you like to specify the file path?").
			Affirmative("Yes").
			Negative("No / Generate one").
			Value(&find)
		if err := form.Run(); err != nil {
			return fmt.Errorf("Error running configuration file prompt: %w", err)
		}
		if find {
			form := huh.NewFilePicker().
				Title("Configuration file").
				Description("Select the config file").
				DirAllowed(true).
				FileAllowed(true).
				ShowPermissions(true).
				ShowHidden(true).
				Value(&file)
			if err := form.Run(); err != nil {
				return fmt.Errorf("Error running file picker: %w", err)
			}
		} else {
			return fmt.Errorf("No configuration file found")
		}
	}

	cfg, err := os.Open(file)
	defer cfg.Close()
	if err != nil {
		return fmt.Errorf("Error opening configuration file: %w", err)
	}
	decoder := yaml.NewDecoder(cfg)
	if err := decoder.Decode(&c); err != nil {
		return fmt.Errorf("Error decoding configuration file: %w", err)
	}
	f, err := os.Stat(c.TemplateDir)
	if err != nil {
		return fmt.Errorf("Error finding template directory: %w", err)
	}
	if !f.IsDir() {
		c.TemplateDir = filepath.Dir(c.TemplateDir)
		return nil
	}
	return nil
}

func (c *config) generateConfiguration(file string) error {
	if err := c.buildForm().Run(); err != nil {
		return fmt.Errorf("Error generating interactive configuration: %w", err)
	}
	cfg, err := os.Create(file)
	if err != nil {
		return fmt.Errorf("Error creating configuration file: %w", err)
	}
	defer cfg.Close()

	conf, err := yaml.Marshal(&c)
	if err != nil {
		return fmt.Errorf("Error marshalling configuration: %w", err)
	}
	err = os.WriteFile(file, conf, 0644)
	if err != nil {
		return fmt.Errorf("Error writing configuration file: %w", err)
	}
	log.Infof("Configuration written to %s", file)
	return nil
}
func (cfg *config) buildForm() *huh.Form {
	var formBits []huh.Field

	formBits = append(formBits,
		huh.NewNote().
			Title("Interactive Configuration Tool").
			Description("A config file will be generated based on your input. Note all options can be overridden at runtime with flags."))

	v := reflect.ValueOf(cfg)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	t := v.Type()
	cwd, err := os.Getwd()
	if err != nil {
		log.Errorf("Error getting home directory: %v", err)
	}
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("form")
		if tag == "" {
			continue
		}
		tagParts := strings.Split(tag, ";")
		for i, p := range tagParts {
			tagParts[i] = strings.TrimSpace(p)
		}
		formType := tagParts[0]
		switch formType {
		case "file":

			title, desc, ext, _ := parseTagOptions(tagParts[1:], field.Name)
			formBits = append(formBits,
				huh.NewFilePicker().
					Title(title).
					Description(desc).
					Picking(false).
					DirAllowed(false).
					FileAllowed(true).
					ShowPermissions(false).
					ShowHidden(false).
					CurrentDirectory(cwd).
					AllowedTypes([]string{ext}).
					Value(v.Field(i).Addr().Interface().(*string)))

		case "dir":
			title, desc, _, _ := parseTagOptions(tagParts[1:], field.Name)
			formBits = append(formBits,
				huh.NewFilePicker().
					Title(title).
					Description(desc).
					Picking(false).
					DirAllowed(true).
					FileAllowed(false).
					ShowPermissions(false).
					ShowHidden(false).
					CurrentDirectory(cwd).
					Value(v.Field(i).Addr().Interface().(*string)))

		case "input":
			title, desc, _, pholder := parseTagOptions(tagParts[1:], field.Name)
			formBits = append(formBits,
				huh.NewInput().
					Title(title).
					Description(desc).
					Placeholder(pholder).
					Value(v.Field(i).Addr().Interface().(*string)))

		case "confirm":
			title, _, _, _ := parseTagOptions(tagParts[1:], field.Name)
			formBits = append(formBits,
				huh.NewConfirm().
					Title(title).
					Value(v.Field(i).Addr().Interface().(*bool)))
		}
	}
	form := huh.NewForm(huh.NewGroup(formBits...)).WithShowHelp(true).WithTheme(huh.ThemeCharm()).WithProgramOptions(tea.WithAltScreen())
	return form
}

func parseTagOptions(options []string, fieldName string) (title, desc, ext, placeholder string) {
	var t, d, e, ph string
	for _, p := range options {
		parts := strings.SplitN(p, "=", 2)
		if len(parts) != 2 {
			log.Errorf("invalid tag %q for field %q", p, fieldName)
			continue
		}
		k, v := parts[0], parts[1]
		switch k {
		case "title":
			t = v
		case "desc":
			d = v
		case "ext":
			e = v
		case "placeholder":
			ph = v
		}
	}
	return t, d, e, ph
}
