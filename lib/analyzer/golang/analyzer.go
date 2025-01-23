package golang

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path/filepath"

	"github.com/snowmerak/indexer/lib/analyzer"
)

var _ analyzer.Analyzer = &Analyzer{}

type Analyzer struct {
}

func (a *Analyzer) Walk(path string, recursive bool, callback func(codeBlock, filePath string, line int) error) error {
	type Entry struct {
		Path  string
		Entry os.DirEntry
	}
	entries := make([]Entry, 0)
	dir, err := os.ReadDir(path)
	if err != nil {
		return fmt.Errorf("failed to read dir: %w", err)
	}

	for _, entry := range dir {
		switch entry.IsDir() {
		case true:
			if recursive {
				entries = append(entries, Entry{Path: path, Entry: entry})
			}
		case false:
			if filepath.Ext(entry.Name()) == ".go" {
				entries = append(entries, Entry{Path: path, Entry: entry})
			}
		}
	}

	for _, entry := range entries {
		if entry.Entry.IsDir() {
			if err := a.Walk(filepath.Join(entry.Path, entry.Entry.Name()), recursive, callback); err != nil {
				return err
			}
		} else {
			fileSet := token.NewFileSet()
			sourceFile := filepath.Join(path, entry.Entry.Name())
			f, err := parser.ParseFile(fileSet, sourceFile, nil, parser.ParseComments)
			if err != nil {
				return fmt.Errorf("failed to parse file: %w", err)
			}

			printCfg := printer.Config{
				Mode:     printer.UseSpaces | printer.TabIndent,
				Tabwidth: 4,
			}
			buf := bytes.NewBuffer(nil)

			ast.Inspect(f, func(node ast.Node) bool {
				use := false
				switch node := node.(type) {
				case *ast.FuncDecl:
					use = true
				case *ast.GenDecl:
					use = node.Tok == token.TYPE || node.Tok == token.CONST || node.Tok == token.VAR
				case *ast.IfStmt:
					use = true
				case *ast.ForStmt:
					use = true
				case *ast.RangeStmt:
					use = true
				case *ast.SwitchStmt:
					use = true
				case *ast.TypeSwitchStmt:
					use = true
				}

				if use {
					buf.Reset()
					if err := printCfg.Fprint(buf, fileSet, node); err != nil {
						return false
					}
					if err := callback(buf.String(), sourceFile, fileSet.Position(node.Pos()).Line); err != nil {
						return false
					}
				}

				return true
			})
		}
	}

	return nil
}
