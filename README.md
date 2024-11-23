# 📄 Resume Generator

[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat-square&logo=go)](https://golang.org)
[![LaTeX](https://img.shields.io/badge/LaTeX-Required-008080?style=flat-square&logo=latex)](https://www.latex-project.org/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg?style=flat-square)](https://opensource.org/licenses/MIT)

Generate professional resumes and cover letters with the power of Go and LaTeX. Separate your content from styling, and focus on what matters most—showcasing your experience.

<div align="center">
  <img src="https://imgs.xkcd.com/comics/automation.png" alt="XKCD Automation Comic" width="400"/>
  <p><i>Yes, I know I'm spending more time on this than I'll save. That's half the fun!</i></p>
</div>

## ✨ Features

- **Content-Style Separation**: Write your resume content in YAML, let LaTeX handle the formatting
- **Live Preview**: Watch your changes render in near real-time as you edit
- **Multiple Output Support**: Generate both resumes and cover letters from a single source
- **Base Resume Templates**: Use a base resume for common information, customize for specific applications
- **Interactive Setup**: User-friendly configuration process to get you started quickly

## 📋 Example Output

<div align="center">
  <a href="John_Doe_example.pdf">
  <p><i>Click to view the full example PDF</i></p>
  </a>
</div>

> This is a sample resume generated using the default template. Your resume can look different based on your chosen template and content.


## 🚀 Quick Start

### Prerequisites

- Go 1.21 or higher
- LaTeX distribution with `pdflatex`

### Installation

```bash
# Clone the repository
git clone https://github.com/munirmah/Resume-Generator.git

# Navigate to project directory
cd Resume-Generator

# Install dependencies
go get .

# Build the project
go build .
```

### Basic Usage

```bash
# Generate a resume with default settings
./Resume-Generator -f your-resume.yml

# Enable live preview while editing
./Resume-Generator -f your-resume.yml -r

# Generate both resume and cover letter
./Resume-Generator -f your-resume.yml -c
```

## 🎯 Core Concepts

### YAML-Based Content

Your resume content lives in a clean, human-readable YAML file:

```yaml
information:
  name: "Jane Doe"
  email: "jane@example.com"
  phone: "123-456-7890"

experience:
  - company: "Tech Corp"
    title: "Senior Developer"
    start_date: "2020-01-01"
    end_date: "Present"
    description:
      - "Led development of core platform features"
      - "Mentored junior developers"
```

### Section Ordering

Customize your resume's section order using simple letter codes:

- `e` - Education
- `x` - Experience
- `s` - Skills
- `p` - Projects
- `c` - Certifications
- `t` - Custom
- `m` - Summary

Example: `-o xsep` renders `Experience → Skills → Education → Projects`

## 🛠️ Configuration

### Command Line Flags

| Flag | Description | Default |
|------|-------------|---------|
| `-f` | Resume YAML file | Required |
| `-b` | Base resume template | Optional |
| `-c` | Generate cover letter | false |
| `-r` | Enable live preview | false |
| `-o` | Section order | Required |
| `-s` | Show PDF after generation | false |
| `-l` | Log level (debug,info,warn,error) | error |

### Configuration File

The tool uses a `.config` file for persistent settings. First run creates this interactively, or override with flags.

## 🔍 Advanced Features

### Base Resume System

Use a base resume for unchanging information:

```bash
./Resume-Generator -b base.yml -f job-specific.yml
```

### Live Preview Mode

Enable real-time PDF updates while editing:

```bash
./Resume-Generator -f resume.yml -r
```

## 🤝 Contributing

Contributions welcome! Feel free to:

- Open issues for bugs or suggestions
- Submit pull requests
- Improve documentation
- Share your templates

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

Made with ❤️ by [munirmah](https://github.com/munirmah)
