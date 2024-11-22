package main

import (
	"bytes"
	"fmt"
	"html"
	"html/template"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/log"
	yaml "gopkg.in/yaml.v3"
)

func (r *resume) execTmpl(templateDir, outDir, filename, order, tmplType string, check bool) error {
	tmplFuncs := template.FuncMap{
		"phone":   phone.String,
		"date":    date.String,
		"url":     social.String,
		"getURL":  social.getURL,
		"icon":    social.getIcon,
		"listify": listify,
		"have":    inFuture,
		"trim":    strings.TrimSpace,
		"today":   func() string { return time.Now().Format("2006-01-02") },
	}

	tex, err := template.New("All").Funcs(tmplFuncs).ParseGlob(path.Join(templateDir, "*.tmpl"))
	if err != nil {
		log.Fatalf("Error parsing templates: %v", err)
	}
	log.Infof("Parsed templates from: %s", templateDir)

	var buffer bytes.Buffer

	switch tmplType {
	case "cover":
		err = tex.ExecuteTemplate(&buffer, "cover", r)
		if err != nil {
			return fmt.Errorf("Error executing cover template: %w", err)
		}
		log.Infof("Successfully executed cover template")
	case "resume":
		err = tex.ExecuteTemplate(&buffer, "header", r)
		if err != nil {
			return fmt.Errorf("Error executing header template: %w", err)
		}
		for _, section := range order {
			err = tex.ExecuteTemplate(&buffer, mapping[section], r)
			if err != nil {
				return fmt.Errorf("Error executing %s template: %w", mapping[section], err)
			}
		}
		err = tex.ExecuteTemplate(&buffer, "footer", r)
		if err != nil {
			return fmt.Errorf("Error executing footer template: %w", err)
		}
		log.Infof("Successfully executed all the templates")
	}

	filepath := path.Join(outDir, filename+".tex")
	if check {
		if _, err := os.Stat(filepath); err == nil {
			var overwrite bool
			log.Warnf("File already exists: %s", filepath)
			form := huh.NewConfirm().
				Title("Overwrite existing file?").
				DescriptionFunc(func() string { return fmt.Sprintf("File already exists: %s", filepath) }, "").
				Affirmative("Overwrite").
				Negative("Cancel").
				Value(&overwrite)
			if err := form.Run(); err != nil {
				return fmt.Errorf("Error running form: %w", err)
			}

			if !overwrite {
				return fmt.Errorf("User chose not to overwrite %s", filepath)
			}
		}
	}
	decoded := html.UnescapeString(buffer.String())
	texFile, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("Error creating tex file: %w", err)
	}
	if _, err := texFile.WriteString(decoded); err != nil {
		return fmt.Errorf("Error writing tex file: %w", err)
	}
	if err := texFile.Close(); err != nil {
		return fmt.Errorf("Error closing tex file: %w", err)
	}
	log.Infof("Successfully wrote tex file")
	return nil
}

func generatePDF(inDir, outDir, filename string) error {
	var cmd *exec.Cmd
	tex := path.Join(inDir, filename+".tex")
	if _, err := os.Stat(tex); err != nil {
		return fmt.Errorf("Error finding tex file: %w", err)
	}
	if _, err := exec.LookPath("pdflatex"); err != nil {
		log.Fatalf("Please make sure pdflatex is installed and in your PATH: %v", err)
	}
	cmd = exec.Command("pdflatex", "-output-directory", outDir, tex)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Error generating PDF: %w\nOutput: %s", err, out)
	}
	if runtime.GOOS == "linux" {
		cleanup := exec.Command("find", outDir, "-type", "f", "!", "-name", "*.pdf", "-delete")
		if _, err := cleanup.CombinedOutput(); err != nil {
			return fmt.Errorf("Error cleaning up files: %w", err)
		}
	}
	log.Infof("Successfully generated PDF")

	return nil
}

func scrapeLinkedin(dir, py string) error {
	if _, err := os.Stat(py); err != nil {
		log.Errorf("Error finding python script: %v", err)
		return err
	}
	cmd := exec.Command("python3", py, dir)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Errorf("Error scraping LinkedIn: %v", err)
		log.Errorf("Output: %s", out)
		return err
	}
	return nil
}

func (r *resume) parseResume(resumeFile string) error {
	f, err := os.Open(resumeFile)
	if err != nil {
		return err
	}
	defer f.Close()
	decoder := yaml.NewDecoder(f)
	if err := decoder.Decode(&r); err != nil {
		return err
	}
	log.Infof("Parsed resume file: %s", resumeFile)
	return nil
}
