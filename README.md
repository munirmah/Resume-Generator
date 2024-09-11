# Resume Generator
Welcome to the **Resume Generator**! This tool is designed to simplify and streamline the process of creating consistent, professional resumes. By separating content from styling, you can focus on crafting compelling resumes while the tool handles the formatting.


--- GIF of the tool in action ---


## Overview
The **Resume Generator** leverages Go and LaTeX to provide powerful and customizatble soltiuon for resume creation. Here's how it works:

 - **Content**: Define your resume content in a `YAML` file, focusing on just the information not the formatting.

 - **Styling**: The styling is managed through LaTeX templates. You can choose from a variety of templates or create your own to match your personal style.

 - **Generation**: The tool combines your resume content with the selected template to create a professional PDF resume.


## Features:

- **YAML-Powered**: Your resume live in a human readable `YAML` file, making it easy to update and maintain.
- **Customizable**: Choose from a variety of templates or create your own to match your personal style.
- **Live Preview**: Make changes to your resume content and see the results rendered in near real-time.
- **Cover Letter Support**: Generate cover letters alongside your resume by including a cover letter in your `YAML` file.
- **Obsidian Integration**: Use the Obsidian Kanban plugin to track your application process in Obsidian.
- **Interactive Configuration**: A user-friendly setup process that guides you through the initial configuration, making it easy to get started.

## Getting Started

### Prerequisites
- **Go**: Ensure you have Go installed on your machine. You can download it [here](https://golang.org).
- **LaTeX**: Install a LaTeX distribution on your machine. We recommend [MiKTeX](https://miktex.org/download) for Windows and [MacTeX](https://www.tug.org/mactex/) for macOS. For Linux, you can use the `texlive` package.
> [!WARNING]
> Only tested on Fedora using 
- **Pandoc**: Install Pandoc on your machine. You can download it [here](https://pandoc.org/installing.html).

### Installation
1. Clone the repository:
```bash
git clone 
```
2. Navigate to the project directory:
```bash
cd resume-generator
```
3. Build the project:
```bash
go build .
```

### Usage

You can run the tool using the following command:

```bash
./resume-generator
```
#### Configuraiton

The tool will, by default, look for a `.config` file in the current directory. If one is not found, it will ask you to locate one or create a new configuration file.
Choosing to create a new configuration will guide you through the setup process and will save it in the current directory.

> [!NOTE]
> You can override any of the configuration options by passing them as flags when running the tool.

The configuration file follows the schema specified in the `config.json` file.

#### Resume Content

##### Base Resume
The tool has the ability to use a `base.yml` file that contains the complete resume content. This allows you to easily generate new resumes without having to re-enter unchanged information such as your personal details.

The `base.yml` file should adhere to the schema specified in the `resume.json` file.




#### Custom Order
You can specify the order in which the sections of your resume are rendered by setting the `order` field in the configuration file. You can also specify the order as a flag when running the tool.

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
> Any sections not included will not be rendered in the final resume.

There are some special values for the order field:
- `all`: This will render all sections in the order they are defined in the `base.yml` in addition to `resume.yml` file. 
- `none`: This will make the tool prompt you to select the sections you want to include in your resume.





## Contributing

Contributions are welcome! Feel free to open issues or submit pull requests.

## License 📄

This project is licensed under the MIT License.





































