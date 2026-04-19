package generator

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/alifengineer/pii/internal/analyzer"
	"golang.org/x/tools/imports"
)

type generator struct {
	buf    *bytes.Buffer
	indent string
}

func newGenerator(buf *bytes.Buffer) *generator {
	return &generator{buf: buf}
}

func (g *generator) p(format string, args ...any) {
	fmt.Fprintf(g.buf, g.indent+format+"\n", args...)
}

func (g *generator) in() {
	g.indent += "\t"
}

func (g *generator) out() {
	g.indent = g.indent[1:]
}

func Generate(packageName string, targets []analyzer.SanitizeTarget, buf *bytes.Buffer) ([]byte, error) {
	g := newGenerator(buf)
	g.emitHeader(packageName)

	for _, target := range targets {
		alloc := newIdentifierAllocator()
		g.emitSanitizeMethod(target, alloc)
	}

	return imports.Process("", buf.Bytes(), nil) //nolint:wrapcheck //best effort
}

func ShouldWrite(output []byte, destination string) bool {
	existing, err := os.ReadFile(destination) //nolint:gosec // path from go:generate, not user input
	if err != nil {
		return true
	}
	return !bytes.Equal(existing, output)
}

func WriteFile(output []byte, destination string) error {
	dir := filepath.Dir(destination)
	tmp, err := os.CreateTemp(dir, "pi-*.go")
	if err != nil {
		return fmt.Errorf("creating temp file: %w", err)
	}
	tmpName := tmp.Name()

	if _, err := tmp.Write(output); err != nil {
		_ = tmp.Close()
		_ = os.Remove(tmpName)
		return fmt.Errorf("writing temp file: %w", err)
	}
	if err := tmp.Close(); err != nil {
		_ = os.Remove(tmpName)
		return fmt.Errorf("closing temp file: %w", err)
	}
	if err := os.Rename(tmpName, destination); err != nil {
		_ = os.Remove(tmpName)
		return fmt.Errorf("renaming to destination: %w", err)
	}
	return nil
}
