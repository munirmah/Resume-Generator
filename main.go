package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
	"github.com/charmbracelet/log"
	"github.com/fsnotify/fsnotify"
)

const (
	githubURL   = "https://www.github.com"
	linkedinURL = "https://www.linkedin.com"
)

// TODO: Read from config file instead of hardcoding
var mapping = map[rune]string{
	'e': "Education",
	'x': "Experience",
	'p': "Projects",
	's': "Skills",
	'c': "Certifications",
	't': "Custom",
	'm': "Summary",
}

func main() {
	configFile := ".config"
	var (
		c            config
		updateConfig bool
		reload       bool
		logLevel     string
		resFile      string
	)

	flag.StringVar(&logLevel, "l", "error", "Set the log level: debug, info, warn, error")
	flag.StringVar(&c.BaseFile, "b", c.BaseFile, "The resume that will be used as a basis for missing information")
	flag.StringVar(&c.TemplateDir, "temp", "./templates", "The directory containing resume templates")
	flag.StringVar(&c.TexDir, "tex", "tex", "The directory where TeX files will be generated. Leave empty to auto create ./tex directory")
	flag.StringVar(&c.PdfDir, "dir", "pdf", "The directory where PDF files will be saved. Leave empty to auto create ./pdf directory")
	flag.StringVar(&c.CoverFile, "cvr", "default", "The name of the generated cover letter file. Default option will autogenerate the name")
	flag.StringVar(&c.PdfFile, "pdf", "default", "The name of the generated PDF file. Default option will autogenerate the name")
	flag.StringVar(&c.KanbanFile, "k", c.KanbanFile, "The Markdown file for your Kanban board")
	flag.StringVar(&c.Order, "o", c.Order, "Enter the order of sections. Missing section will be omitted\nEnter none to be prompted everytime")
	flag.BoolVar(&reload, "r", reload, "Enable live reloading of the resume file")
	flag.BoolVar(&c.Cover, "c", c.Cover, "Generate a Cover Letter?")
	flag.BoolVar(&updateConfig, "config", false, "Update the current configuration file")
	flag.BoolVar(&c.Track, "t", c.Track, "Whether to track changes in Obsidian?")
	flag.BoolVar(&c.Show, "s", c.Show, "Show PDF after creation?")
	flag.Parse()

	switch logLevel {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	default:
		log.Errorf("Invalid log level: %s", logLevel)
		log.SetLevel(log.ErrorLevel)
	}

	err := c.readConfig(configFile)
	if err != nil || updateConfig {
		var genConfig bool
		if updateConfig {
			log.Warnf("Updating configuration file")
			genConfig = true
		} else {
			log.Warnf("Error reading configuration file: %v", err)

			form := huh.NewConfirm().
				Title("New Configuration").
				Description("Would you like to generate the config file?").
				Affirmative("Yes").
				Negative("No").
				Value(&genConfig)
			if err := form.Run(); err != nil {
				log.Fatalf("Error generating configuration file: %v", err)
			}
		}
		if genConfig {
			err := c.generateConfiguration(configFile)
			if err != nil {
				log.Fatalf("Error generating configuration file: %v", err)
			}
		} else {
			log.Printf("You must configure this application to run. Either create interactively or pass in the following flags:")
			flag.PrintDefaults()
			os.Exit(1)
		}
	}

	if _, err := os.Stat(c.TemplateDir); os.IsNotExist(err) || c.TemplateDir == "" {
		log.Fatalf("Template directory does not exist: %s", c.TemplateDir)
	}
	if _, err := os.Stat(c.TexDir); os.IsNotExist(err) {
		log.Warnf("Tex directory does not exist: %s", c.TexDir)
		if err := os.MkdirAll(c.TexDir, 0755); err != nil {
			log.Fatalf("Error creating Tex directory: %v", err)
		}
		log.Infof("Created Tex directory: %s", c.TexDir)
	}
	if _, err := os.Stat(c.PdfDir); os.IsNotExist(err) {
		log.Warnf("PDF directory does not exist: %s", c.PdfDir)
		if err := os.MkdirAll(c.PdfDir, 0755); err != nil {
			log.Fatalf("Error creating PDF directory: %v", err)
		}
		log.Infof("Created PDF directory: %s", c.PdfDir)
	}

	var obsidianDir string
	if c.Track || c.KanbanFile != "" {
		obsidianDir = path.Dir(c.KanbanFile)
		if _, err := os.Stat(obsidianDir); os.IsNotExist(err) {
			log.Warnf("Obsidian directory %s does not exist. Skipping tracking", obsidianDir)
			c.Track = false
		}
	}

	if resFile == "" {
		log.Warnf("No resume file provided. Prompting for file")
		var find bool
		form := huh.NewConfirm().
			Title("Resume file not Provided").
			Description("Would you like to specify the file path?").
			Affirmative("Yes").
			Negative("No").
			Value(&find)
		if err := form.Run(); err != nil {
			log.Errorf("Error running configuration file prompt: %v", err)
		}
		if find {
			cwd, err := os.Getwd()
			if err != nil {
				log.Fatalf("Error getting current working directory: %v", err)
			}
			form := huh.NewFilePicker().
				Title("Select resume file").
				DirAllowed(false).
				Picking(true).
				FileAllowed(true).
				ShowPermissions(false).
				ShowHidden(false).
				CurrentDirectory(cwd).
				AllowedTypes([]string{"yaml", "yml"}).
				Value(&resFile).
				WithHeight(20)
			if err := form.Run(); err != nil {
				log.Errorf("Error running file picker: %v", err)
			}
		} else {
			log.Fatalf("No resume file provided. Pass in the resume file with the -f flag")
		}
	}

	var res resume
	if c.BaseFile == "" {
		log.Warnf("No base resume file provided. Skipping base resume")
	} else {
		err := res.parseResume(c.BaseFile)
		if err != nil {
			log.Fatalf("Error parsing resume file: %s - %v", c.BaseFile, err)
		}
	}

	err := res.parseResume(resFile)
	if err != nil {
		log.Fatalf("Error parsing resume file: %s - %v", resFile, err)
	}

	err = res.sanitizeResume()
	if err != nil {
		log.Fatal("Error sanitizing resume file:", err)
	}

	c.Order = strings.ToLower(c.Order)

	if c.Order == "all" {
		c.Order = ""
		for k := range mapping {
			c.Order += string(k)
		}
	}

	if c.Order == "none" || c.Order == "" {
		c.Order = ""
		huh.NewInput().
			Title("Order of Sections").
			DescriptionFunc(func() string {
				var desc string
				for k, v := range mapping {
					desc += fmt.Sprintf("[%c]%s\t", k, v)
				}
				return desc
			}, "").
			Placeholder("expstc").
			//TODO: Add validation for order
			// Validate(func(str string) error {
			// 	reg := regexp.MustCompile()
			// 	if !reg.MatchString(str) {
			// 		return fmt.Errorf("invalid order")
			// 	}
			// 	return nil
			// }).
			Value(&c.Order).
			Run()
	}

	resumeName := getFilename(resFile)
	name := strings.ReplaceAll(res.Info.Name, " ", "_")

	if c.PdfFile == "default" {
		c.PdfFile = fmt.Sprintf("%s_%s", name, resumeName)
	}
	if c.CoverFile == "default" {
		c.CoverFile = fmt.Sprintf("%s_%s_cvr", name, resumeName)
	}

	err = res.execTmpl(c.TemplateDir, c.TexDir, c.PdfFile, c.Order, "resume", true)
	if err != nil {
		log.Fatalf("Error executing templates: %v", err)
	}

	err = generatePDF(c.TexDir, c.PdfDir, c.PdfFile)
	if err != nil {
		log.Fatalf("Error generating PDF: %v", err)
	}
	log.Infof("Generated PDF: %s", c.PdfFile)

	if c.Cover {
		err = res.execTmpl(c.TemplateDir, c.TexDir, c.CoverFile, "", "cover", true)
		if err != nil {
			log.Fatalf("Error executing cover letter template: %v", err)
		}
		log.Infof("Generated cover letter TeX file: %s", c.CoverFile)

		err = generatePDF(c.TexDir, c.PdfDir, c.CoverFile)
		if err != nil {
			log.Fatalf("Error generating cover letter: %v", err)
		}
		log.Infof("Generated cover letter: %s", c.CoverFile)
	}
	if c.Show {
		pdf := path.Join(c.PdfDir, c.PdfFile+".pdf")
		if err := openFile(pdf); err != nil {
			log.Fatalf("Error opening file %s: %v", c.PdfFile, err)
		}
		if c.Cover {
			cvr := path.Join(c.PdfDir, c.CoverFile+".pdf")
			if err := openFile(cvr); err != nil {
				log.Fatalf("Error opening file %s: %v", c.CoverFile, err)
			}
		}
	}

	if c.Track {
		//FIXME: Test Obsidian tracking
		log.Warnf("Tracking PDF in Obsidian is not fully tested yet. Expect bugs")
		var mdBuff bytes.Buffer
		tmplFuncs := template.FuncMap{
			"today": func() string { return time.Now().Format("2006-01-02") },
		}
		md, err := template.New("obsidian").Funcs(tmplFuncs).ParseGlob(path.Join(c.TemplateDir, "*.tmpl"))
		err = md.ExecuteTemplate(&mdBuff, "obsidian", res)
		if err != nil {
			log.Fatalf("Error executing Obsidian template: %v", err)
		}
		fname := fmt.Sprintf("%s.md", c.PdfFile)
		_, exists := os.Stat(path.Join(obsidianDir, fname))
		if exists == nil {
			log.Warnf("Obsidian file already exists: %s", fname)
			form := huh.NewInput().
				Title("File already exists").
				DescriptionFunc(func() string {
					return fmt.Sprintf("%s already exists. Provide a new file name", fname)
				}, "").
				Placeholder(fname).
				Value(&fname).
				Validate(func(str string) error {
					s := filepath.Base(str)
					if s == "" || s == "." {
						return fmt.Errorf("file name cannot be empty")
					}
					_, exists := os.Stat(path.Join(obsidianDir, s))
					if exists == nil {
						return fmt.Errorf("file already exists")
					}
					return nil
				})
			if err := form.Run(); err != nil {
				log.Fatalf("Error prompting for overwrite: %v", err)
			}
			fname = filepath.Base(fname)
		}

		err = os.WriteFile(path.Join(obsidianDir, fname), mdBuff.Bytes(), 0644)
		if err != nil {
			log.Fatalf("Error writing Obsidian file: %v", err)
		}
		log.Infof("Generated Obsidian file: %s", fname)
		kanban, err := os.Open(path.Join(obsidianDir, c.KanbanFile))
		if err != nil {
			log.Fatalf("Error opening Status file: %v", err)
		}
		defer kanban.Close()
		scanner := bufio.NewScanner(kanban)
		var updated []string
		for scanner.Scan() {
			line := scanner.Text()
			updated = append(updated, line)
			list := "## " + c.KanbanListName
			if strings.Contains(line, list) {
				updated = append(updated, fmt.Sprintf(" - [ ] [[%s|%s]]\n", fname, res.Job.Title))
			}
		}
		if err := scanner.Err(); err != nil {
			log.Fatalf("Error reading Status file: %v", err)
		}
		err = os.WriteFile(path.Join(obsidianDir, c.KanbanFile), []byte(strings.Join(updated, "")), 0644)
		if err != nil {
			log.Fatalf("Error writing Status file: %v", err)
		}
		log.Infof("Updated Kanban file: %s", c.KanbanFile)
	}

	if reload {
		log.Infof("Live realoding enabled")

		w, err := fsnotify.NewWatcher()
		if err != nil {
			log.Fatalf("Error creating watcher: %v", err)
		}
		defer w.Close()

		title := fmt.Sprintf("Watching for changes to %s\n  Press CTRL+C to exit", filepath.Base(resFile))
		err = spinner.New().
			Title(title).
			Action(func() {
				st, err := os.Lstat(resFile)
				if err != nil {
					log.Fatalf("Error getting file info: %v", err)
				}
				if st.IsDir() {
					log.Fatalf("File is a directory: %s", resFile)
				}
				err = w.Add(filepath.Dir(resFile))
				if err != nil {
					log.Fatalf("Error adding file to watcher: %v", err)
				}
				for {
					select {
					case e, ok := <-w.Events:
						if !ok {
							return
						}
						if e.Name == resFile {
							if e.Op.Has(fsnotify.Write) {
								time.Sleep(1 * time.Second)
								res.parseResume(resFile)
								err = res.sanitizeResume()
								if err != nil {
									log.Fatalf("Error sanitizing resume file: %v", err)
								}
								err = res.execTmpl(c.TemplateDir, c.TexDir, c.PdfFile, c.Order, "resume", false)
								if err != nil {
									log.Fatalf("Error executing templates: %v", err)
								}
								err = generatePDF(c.TexDir, c.PdfDir, c.PdfFile)
								if err != nil {
									log.Fatalf("Error generating PDF: %v", err)
								}
								log.Infof("Generated PDF: %s", c.PdfFile)
								openFile(path.Join(c.PdfDir, c.PdfFile+".pdf"))
							}
						}
					case err := <-w.Errors:
						log.Errorf("Error watching file: %v", err)
						return
					}
				}
			}).
			Type(spinner.Dots).
			Run()

		log.Printf("Stopped watching for changes. Exiting...")
		os.Exit(0)
	}

}
