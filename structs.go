package main

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	yaml "gopkg.in/yaml.v3"
)

type config struct {
	BaseFile       string `yaml:"base,omitempty" form:"file; title=Base Resume File; desc=The resume that will be used as a basis for missing information\nLeave empty to ignore; ext=yml"`
	TemplateDir    string `yaml:"template,omitempty" form:"file; title=Template Directory; desc=The directory containing resume templates;ext=tmpl"`
	TexDir         string `yaml:"tex,omitempty" form:"dir; title=TeX Output Directory; desc=The directory where TeX files will be generated\nLeave empty to auto create ./tex directory"`
	PdfDir         string `yaml:"pdf_dir,omitempty" form:"dir; title=PDF Output Directory; desc=The directory where PDF files will be saved\nLeave empty to auto create ./pdf directory"`
	CoverFile      string `yaml:"cover_file,omitempty" form:"input; title=Cover Letter File Name; desc=The name of the generated cover letter file\nDefault option with autogenerate the name; placeholder=default"`
	PdfFile        string `yaml:"pdf,omitempty" form:"input; title=PDF File Name; desc=The name of the generated PDF file\nDefault option will autogenerate the name; placeholder=default"`
	Track          bool   `yaml:"track,omitempty" form:"confirm; title=Track changes in Obsidian"`
	KanbanFile     string `yaml:"kanban,omitempty" form:"file; title=Kanban Board; desc=The Markdown file for your Kanban board; ext=md"`
	KanbanListName string `yaml:"kanban_list_name,omitempty" form:"input; title=Kanban List Name; desc=The name of the list in the Kanban board that new jobs will be added under; placeholder=To Apply"`
	Order          string `yaml:"order,omitempty" form:"input; title=Default Resume Section Order; desc=Enter the order of sections. Missing section will be omitted:\n\t[e]ducation, e[x]perience, [p]rojects, [s]kills, [c]ertifications, cus[t]om, su[m]mary\nEnter none to be prompted everytime; placeholder=none"`
	Cover          bool   `yaml:"cover,omitempty" form:"confirm; title=Generate a Cover Letter"`
	Show           bool   `yaml:"show,omitempty" form:"confirm; title=Show PDF after creation"`
}

type resume struct {
	Job            job             `yaml:"job"`            // Job application details (Optional)
	Info           info            `yaml:"information"`    // Information of the Person in the Resume (Required)
	Education      []school        `yaml:"education"`      // Education of the Person in the Resume (Required)
	Experiences    []experience    `yaml:"experience"`     // Experiences of the Person in the Resume (Required)
	Projects       []project       `yaml:"projects"`       // Projects of the Person in the Resume (Optional)
	Skills         []skill         `yaml:"skills"`         // Skills of the Person in the Resume (Optional)
	Certifications []certification `yaml:"certifications"` // Certifications of the Person in the Resume (Optional)
	Custom         custom          `yaml:"custom"`         // Custom Section of the Person in the Resume (Optional)
	Summary        summary         `yaml:"summary"`        // Summary Section of the Person in the Resume (Optional)
	CoverLetter    coverLetter     `yaml:"cover_letter"`   // Cover Letter of the Person in the Resume (Optional)
}

type job struct {
	Title    string `yaml:"title"`    // Title of the Job (Required) Example: Software Engineer
	Company  string `yaml:"company"`  // Company of the Job (Required) Example: Google
	Location string `yaml:"location"` // Location of the Job (Required) Example: Mountain View, CA
	URL      string `yaml:"url"`      // URL of the Job (Optional) Example: https://www.google.com
	UUID     string
}

type info struct {
	Name    string   `yaml:"name"`    // Name of the Person in the Resume (Required) Example: John Decode
	Address address  `yaml:"address"` // Address of the Person in the Resume (Optional)
	Email   string   `yaml:"email"`   // Email of the Person in the Resume (Required) Example: name@example.com
	Phone   phone    `yaml:"phone"`   // Phone of the Person in the Resume (Required) Example: 1234567890
	Socials []social `yaml:"socials"` // Social Media of the Person in the Resume (Optional)
}

type address struct {
	Street string `yaml:"street"` // Street of the Address (Required) Example: 123 Main St
	City   string `yaml:"city"`   // City of the Address (Required) Example: New York
	State  string `yaml:"state"`  // State of the Address (Required) Example: NY
	Zip    string `yaml:"zip"`    // Zip of the Address (Required) Example: 10001
}

type phone struct {
	Number string `yaml:"phone"` // Number of the Phone (Required) Example: 1234567890
}

type social struct {
	Platform string `yaml:"platform"` // Platform of the Social Media (Optional) Example: GitHub
	Username string `yaml:"username"` // Username of the Person in the Social Media (Optional) Example: johndecode
	url      string
	icon     string
}

type school struct {
	Name      string `yaml:"name"`       // Name of the School (Required) Example: University of Science
	StartDate date   `yaml:"start_date"` // Start Date of the School (Required) Example: 2018-08-01
	EndDate   date   `yaml:"end_date"`   // End Date of the School (Required) Example: 2022-05-01
	Major     string `yaml:"major"`      // Major of the School (Required) Example: Computer Science
	Minor     string `yaml:"minor"`      // Minor of the School (Optional) Example: Mathematics
	Location  string `yaml:"location"`   // Location of the School (Required) Example: New York, NY
}

type date struct {
	time time.Time `yaml:"start_date,issue_date,expiration_date,end_date,omitempty"` // Time of the Date (Optional) Example: 2022-05-01
	text string    `yaml:"text,omitempty"`                                           // Text of the Date (Optional) Example: Present
}

type experience struct {
	Company     string   `yaml:"company"`     // Company of the Job (Required) Example: Google
	Title       string   `yaml:"title"`       // Title of the Job (Required) Example: Software Engineer
	StartDate   date     `yaml:"start_date"`  // Start Date of the Job (Required) Example: 2022-05-01
	EndDate     date     `yaml:"end_date"`    // End Date of the Job or "Present" (Optional) Example: 2022-05-01
	Location    string   `yaml:"location"`    // Location of the Job (Required) Example: Mountain View, CA
	Description []string `yaml:"description"` // Description of the Job (Required)
}

type project struct {
	Name         string   `yaml:"name"`         // Name of the Project (Required) Example: Resume Builder
	Description  string   `yaml:"description"`  // Description of the Project (Required) Example: A tool to generate resumes
	Technologies []string `yaml:"technologies"` // Technologies used in the Project (Required) Example: [Go, LaTeX]
}

type skill struct {
	Name     string   `yaml:"name"`     // Name of the Skill (Required) Example: Programming
	Keywords []string `yaml:"keywords"` // Keywords of the Skill (Required) Example: [Go, Python]
}

type certification struct {
	Name           string `yaml:"name"`            // Name of the Certification (Required) Example: AWS Certified Solutions Architect
	IssuingOrg     string `yaml:"issuing_org"`     // Issuing Organization of the Certification (Required) Example: Amazon Web Services
	URL            string `yaml:"url"`             // URL of the Certification (Optional) Example: https://www.aws.com
	IssueDate      date   `yaml:"issue_date"`      // Issue Date of the Certification (Required) Example: 2022-05-01
	ExpirationDate date   `yaml:"expiration_date"` // Expiration Date of the Certification (Optional) Example: 2022-05-01
}

type custom struct {
	Title       string   `yaml:"title"`       // Title of the Custom Section (Optional) Example: Hobbies
	Description string   `yaml:"description"` // Description of the Custom Section (Optional) Example: Playing Chess
	Body        []string `yaml:"body"`        // Body of the Custom Section (Optional)
}

type summary struct {
	Title string `yaml:"title"` // Title of the Summary Section (Optional)
	Body  string `yaml:"body"`  // Body of the Summary Section (Optional)
}

type coverLetter struct {
	Name     string  `yaml:"name"`     // Name of the Person in the Cover Letter (Optional) Example: John Deo
	Company  string  `yaml:"company"`  // Company of the Cover Letter (Optional) Example: Google
	Title    string  `yaml:"title"`    // Title of the Cover letter recipient (Optional) Example: Hiring Manager
	Address  address `yaml:"address"`  // Address of the Cover Letter (Optional)
	Greeting string  `yaml:"greeting"` // Greeting of the Cover Letter (Required) Example: Dear Hiring Manager,
	Body     string  `yaml:"body"`     // Body of the Cover Letter (Required)
}

func (s social) getURL() string {
	p, u := strings.ToLower(s.Platform), strings.ToLower(s.Username)
	switch p {
	case "github":
		return fmt.Sprintf("%s/%s", githubURL, u)
	case "linkedin":
		return fmt.Sprintf("%s/in/%s", linkedinURL, u)
	}
	return ""
}

func (s social) getIcon() string {
	p := strings.ToLower(s.Platform)
	switch p {
	case "github":
		return "\\faGithub"
	case "linkedin":
		return "\\faLinkedin"
	}
	return ""
}

func inFuture(d date) bool {
	if d.text != "" {
		return true
	}
	return d.time.After(time.Now())
}

func (t date) String() string {

	if t.text != "" {
		caser := cases.Title(language.English)
		t.text = caser.String(t.text)
		return t.text
	} else {
		return t.time.Format("Jan 2006")
	}
}

func (p phone) String() string {
	n := p.Number
	if len(n) != 10 {
		return n
	}
	return fmt.Sprintf("(%s) %s-%s", n[0:3], n[3:6], n[6:])
}

func (r resume) String() string {
	return fmt.Sprintf("Name: %s\nEmail: %s\nPhone: %s\n", r.Info.Name, r.Info.Email, r.Info.Phone)
}

func (s social) String() string {
	p := strings.ToLower(s.Platform)
	u := strings.ToLower(s.Username)
	switch p {
	case "github":
		return fmt.Sprintf("%s/%s", githubURL[12:], u)
	case "linkedin":
		return fmt.Sprintf("%s/in/%s", linkedinURL[12:], u)
	}
	return fmt.Sprintf("%s@%s", u, p)
}

func (c config) String() string {
	var sb strings.Builder

	sb.WriteString("config:\n")

	// Iterate over each field using reflection
	v := reflect.ValueOf(c)
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		// Get the YAML tag, if any
		yamlTag := field.Name
		if yamlTag == "" {
			continue // Skip fields without YAML tags
		}

		// Format the field and its value
		sb.WriteString("  ")
		sb.WriteString(yamlTag)
		sb.WriteString(": ")

		// Handle booleans specially
		if value.Kind() == reflect.Bool {
			sb.WriteString(strconv.FormatBool(value.Bool()))
		} else {
			sb.WriteString(fmt.Sprintf("%v", value.Interface()))
		}

		sb.WriteString("\n")
	}

	return sb.String()
}

func (p *phone) UnmarshalYAML(value *yaml.Node) error {
	var s string
	err := value.Decode(&s)
	if err != nil {
		return err
	}
	p.Number = s
	return nil
}

func (t *date) UnmarshalYAML(value *yaml.Node) error {
	var s string
	err := value.Decode(&s)
	if err != nil {
		return err
	}
	layout := "2006-01-02"
	parsedTime, err := time.Parse(layout, s)
	if err != nil {
		t.text = s
		return nil
	}
	t.text = ""
	t.time = parsedTime
	return nil
}
