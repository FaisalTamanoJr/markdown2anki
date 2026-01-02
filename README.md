## Objectives

I made this to assist me in my studies, so that I can easily make flash cards while taking notes and reading.

**Main design goal: easy to install, learn, and use**

1. To be able to easily convert markdown notes to a file that can easily be imported in Anki via the command line.
2. Markdown notes structure **should strictly follow** the structure I define.
3. The template for the notes should be intuitive and easy to use. It being Anki-convertable shouldn’t compromise it’s readability and conduciveness for learning.
4. It shouldn’t be difficult to install/setup.
5. It has an isolated design, that is, a converted deck corresponds to a single lesson.
6. It **only utilizes Basic and Cloze**, thereby image-occlusion will not be taken into account.
7. It is intentionally not flexible: **minimal or no configuration** options will be provided. Output should be very predictable.

## Installation

Download it in the [releases](https://github.com/FaisalTamanoJr/markdown2anki/releases) page

If you want to compile it yourself, just download [Go](https://go.dev/), clone the repository and open the directory. Then, run in the terminal (while in the repo directory) `go build markdown2anki.go`.

## Usage

- In the markdown file, make sure that a bullet point text (in the **1. Facts** section) only has either italics or bold but not both.
- Check out the [template.md](https://github.com/FaisalTamanoJr/markdown2anki/blob/main/template.md?plain=1) as a basis for writing md files for anki conversion

```
Usage: markdown2anki.exe [--output_directory OUTPUT_DIRECTORY] INPUT_MD_FILE

Positional arguments:
  INPUT_MD_FILE

Options:
  --output_directory OUTPUT_DIRECTORY, -O OUTPUT_DIRECTORY [default: .]
  --help, -h             display this help and exit

```