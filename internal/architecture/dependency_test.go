package architecture

import (
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Enforces that public SDK/instrumentation packages do not depend on internal packages.
func TestPkgLayerDoesNotImportInternal(t *testing.T) {
	root := filepath.Join("..", "..")
	pkgRoot := filepath.Join(root, "pkg")
	fset := token.NewFileSet()

	err := filepath.WalkDir(pkgRoot, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".go") {
			return nil
		}
		file, parseErr := parser.ParseFile(fset, path, nil, parser.ImportsOnly)
		if parseErr != nil {
			t.Fatalf("parse %s: %v", path, parseErr)
		}
		for _, imp := range file.Imports {
			v := strings.Trim(imp.Path.Value, "\"")
			if strings.Contains(v, "/internal/") {
				t.Fatalf("public package file %s imports internal package %s", path, v)
			}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}
