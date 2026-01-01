package main

import (
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/alexflint/go-arg"
	"github.com/kaleocheng/goldmark"
	"github.com/kaleocheng/goldmark/ast"
	"github.com/kaleocheng/goldmark/extension"
	"github.com/kaleocheng/goldmark/text"
)

type card struct {
	card_type string
	front     string
	back      string
}

/*
func append_to_anki(md_card card) {

}
*/

func main() {
	// Assign variables & CMDLINE parsing
	var args struct {
		Input_md_file    string `arg:"positional, required, -I"`
		Output_directory string `arg:"-O" default:"."`
	}
	arg.MustParse(&args)

	// Create file
	output_filename := filepath.Join(path.Base(args.Output_directory), strings.Replace(args.Input_md_file, ".md", ".txt", 1))
	output_file, err := os.Create(output_filename)

	if err != nil {
		panic("Failed to create file: " + err.Error())
	}
	defer output_file.Close()

	// Write preamble
	output_file.WriteString("#separator:tab\n#html:false\n#notetype column:1\n#deck column:2")

	// Input Markdown to Goldmark AST
	source, _ := os.ReadFile(args.Input_md_file)
	md := goldmark.New(goldmark.WithExtensions(extension.Table))
	md_parsed := md.Parser().Parse(text.NewReader(source))

	// Walk the md input AST tree and track current section
	var current_section string
	ast.Walk(md_parsed, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		// if it is a heading, table, or list
		switch node := n.(type) {
		case *ast.Heading:
			current_section = string(node.FirstChild().Text(source))

		case *ast.Table:
			if current_section == "0. Terminologies" {
				// TODO
			}

		case *ast.List:
			if current_section == "## 1. Facts" {
				// TODO
			}
		}

		return ast.WalkContinue, nil
	})
}
