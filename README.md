# Resume Generator
Welcome to the **Resume Generator**! This tool is designed to simplify and streamline the process of creating consistent, professional resumes. By separating content from styling, you can focus on crafting compelling resumes while the tool handles the formatting.

## Motivation
Crafting tailored resumes for each job application can be a time-consuming and tedious process. I wanted to easily change out the content of my LaTeX resume without the hassle of manually updating the file each time. Thus **Resume Generator** was born.

![Relevent XKCD](https://imgs.xkcd.com/comics/automation.png)

## Overview
The **Resume Generator** leverages Go and LaTeX to provide powerful and customizatble solution for resume creation. The princile behind this tool is to create a separation of concerns:

 - **Content**: Define your resume content in a `YAML` file.
 - **Styling**: The typesetting is managed through LaTeX templates.
 - **Generation**: The tool combines your resume content with the selected template to create a PDF output.

> [!CAUTION]
> Developed and tested on Fedora Linux. The tool may not work as expected on other operating systems. 

## Features:

- **YAML-Powered**: Your resume live in a human readable `YAML` file, making it easy to update and maintain.
- **Customizable**: Choose from a variety of templates or create your own to match your personal style.
- **Live Preview**: Make changes to your resume content and see the results rendered in *near* real-time.
- **Cover Letter Support**: Generate cover letters alongside your resume by including a cover letter in your `YAML` file.
- **Obsidian Integration**: Use the Obsidian Kanban plugin to track your application process in Obsidian.
- **Interactive Configuration**: A user-friendly setup process that guides you through the initial configuration, making it easy to get started.

## Getting Started

### Prerequisites
- **Go**: Ensure you have Go installed on your machine. You can download it [here](https://golang.org).
- **LaTeX**: Install a LaTeX distribution on your machine such as [TeX Live](https://www.tug.org/texlive/).
- **Obsidian**: If you plan to use the Obsidian integration, ensure you have Obsidian installed on your machine. You can download it [here](https://obsidian.md).

### Installation
1. Clone the repository:

```bash
git clone https://github.com/munirmah/Resume-Generator.git 
```

2. Navigate to the project directory:

```bash
cd Resume-Generator
```

3. Install the dependencies:

```bash
go mod tidy
```

4. Build the project:

```bash
go build .
```

## Usage

You can run the tool using the following command:

```bash
./Resume-Generator
```

### Resume Flags


#### Base Resume Selection
You can specify a `base.yml` file to use as the base resume content by passing the `-b` flag.
This allows you to easily generate new resumes without having to include unchanging information such as your contact details.

The `base.yml` file should adhere to the schema specified in the `resume.json` file.

#### Cover Letter Generation
You can generate a cover letter alongside your resume by passing the `-c` flag.

#### Config Regeneration
You can regenerate the configuration file by passing the `-config` flag and overwriting the existing configuration.

#### Cover Letter File Name
You can specify the name of the cover letter file by passing the `-cvr` flag.

> [!NOTE]
> The default value of `"default"` will generate a file with Name, resume file name and "cvr".
> For example, if the resume file name is `resume.yml` and the name is `John Doe`, the cover letter file will be named `John_Doe_resume_cvr.pdf`.

#### Directory for Output
You can specify the directory where the PDF output files will be saved by passing the `-dir` flag.

#### File with Resume Content
You can specify the file containing the resume content by passing the `-f` flag.

#### Kanban File
You can specify the file `md` containing the Kanban board for the Obsidian integration by passing the `-k` flag.

#### Logging Level
You can specify the logging level by passing the `-l` flag. The available options are:
 - `info`
 - `warn`
 - `error`
 - `debug`

#### Order of Sections
You can specify the order in which the sections of your resume are rendered by passing the `-o` flag.

> [!NOTE]
> The letter to section mapping is as follows:
> - `e`: Education
> - `x`: Experience
> - `s`: Skills
> - `p`: Projects
> - `c`: Certifications
> - `t`: Custom
> - `m`: Summary
>
> **Any sections not included will not be rendered in the final resume.**
>
>>  For example:
>> ```bash
>> ./Resume-Generator -o xsm
>> ```
>> This will render the Experience, Skills, and Summary sections in that order.

There are some special values for the order field:
- `all`: This will render all sections in the order they are defined in the `base.yml` in addition to `resume.yml` file. 
- `none`: This will make the tool prompt you to select the sections you want to include in your resume.


#### PDF Output File Name
You can specify the name of the PDF output file by passing the `-pdf` flag.

> [!NOTE]
> The default value of `"default"` will generate a file with Name, resume file name.
> For example, if the resume file name is `resume.yml` and the name is `John Doe`, the cover letter file will be named `John_Doe_resume.pdf`.

#### Real-Time Preview
You can enable real-time preview of the resume by passing the `-r` flag.

#### Show the Generated PDF
You can enable the tool to open the generated PDF file by passing the `-s` flag.

#### Track Application in Obsidian
You can enable the Obsidian integration by passing the `-t` flag.

#### Template Directory
You can specify the directory where the LaTeX templates are stored by passing the `-temp` flag.

#### Latex Directory
You can specify the directory where the LaTeX files are stored by passing the `-tex` flag.

### Configuraiton

The tool will, by default, look for a `.config` file in the current directory. If not found, it will ask you to locate one or create a new configuration file.
Choosing to create a new configuration will guide you through the setup process and will save it in the current directory.

This helps you to avoid passing the same flags every time you run the tool.

> [!NOTE]
> You can still override any of the configuration options by passing them as flags when running the tool.

The configuration file follows the schema specified in the `config.json` file.

## Templates
The templates are stored in the `templates` directory. You can create your own templates or modify the existing ones to match your personal style. They are written in LaTeX and follow the standard `html/template` syntax.

## Contributing

Contributions are welcome! Feel free to open issues or submit pull requests.

## License

This project is licensed under the MIT License.





































