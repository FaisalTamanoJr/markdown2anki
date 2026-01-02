package main

import (
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/alexflint/go-arg"
	"github.com/kaleocheng/goldmark"
	"github.com/kaleocheng/goldmark/ast"
	"github.com/kaleocheng/goldmark/extension"
	extensionast "github.com/kaleocheng/goldmark/extension/ast"
	"github.com/kaleocheng/goldmark/text"
)

type card struct {
	card_type string
	front     string
	back      string
}

func extract_text(node ast.Node, source []byte) string {
	var sb strings.Builder

	ast.Walk(node, func(child ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		if text, ok := child.(*ast.Text); ok {
			sb.Write(text.Text(source))
		}

		return ast.WalkContinue, nil
	})

	return strings.TrimSpace(sb.String())
}

func extract_table_basic(table_node *extensionast.Table, title string, source []byte, output_file os.File) {
	ast.Walk(table_node, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		if row, ok := n.(*extensionast.TableRow); ok { // if node is a table row
			append_to_anki(create_basic_card(row, source), title, output_file)
			return ast.WalkSkipChildren, nil
		}

		return ast.WalkContinue, nil
	})
}

func create_basic_card(row *extensionast.TableRow, source []byte) card {
	return card{
		front:     extract_text(row.FirstChild(), source),
		back:      extract_text(row.FirstChild().NextSibling(), source),
		card_type: "Basic",
	}
}

func extract_list_cloze(list_node *ast.List, title string, source []byte, output_file os.File) {
	ast.Walk(list_node, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		if point, ok := n.(*ast.ListItem); ok { // if node is a bullet point
			append_to_anki(create_cloze_card(point, source), title, output_file)
			return ast.WalkSkipChildren, nil
		}

		return ast.WalkContinue, nil
	})
}

func create_cloze_card(list_node *ast.ListItem, source []byte) card {
	var sb strings.Builder

	cloze_no := 1
	ast.Walk(list_node, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		switch node := n.(type) {
		case *ast.Text:
			sb.Write(node.Text(source))

		case *ast.Emphasis:
			emphasized_text := extract_text(node, source)
			sb.WriteString("{{c" + strconv.Itoa(cloze_no) + "::" + emphasized_text + "}}")
			if node.Level == 2 {
				cloze_no++
			}
			return ast.WalkSkipChildren, nil
		}

		return ast.WalkContinue, nil
	})

	return card{
		front:     "",
		back:      sb.String(),
		card_type: "Cloze",
	}
}

func append_to_anki(md_card card, title string, output_file os.File) {
	output_file.WriteString("\n" + md_card.card_type + "\t" + title + "\t" + md_card.back + "\t" + md_card.front)
}

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
	var deck_title string
	ast.Walk(md_parsed, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		// if it is a heading, table, or list
		switch node := n.(type) {
		case *ast.Heading:
			if node.Level == 1 {
				deck_title = extract_text(node, source)
			} else if node.Level == 2 {
				current_section = extract_text(node, source)
			}

		case *extensionast.Table:
			if current_section == "0. Terminologies" {
				extract_table_basic(node, deck_title, source, *output_file)
			}

		case *ast.List:
			if current_section == "1. Facts" {
				extract_list_cloze(node, deck_title, source, *output_file)
			}
		}

		return ast.WalkContinue, nil
	})

	// Add a newline at the end (haven't tested if this is necessary but my reference export anki card did this)
	output_file.WriteString("\n")
}
