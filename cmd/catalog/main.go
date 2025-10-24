package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/rselbach/utils/internal/catalog"
)

func main() {
	root := flag.String("root", ".", "root directory to scan for utilities")
	baseURL := flag.String("base-url", "https://utils.rselbach.com", "public base URL for utilities")
	out := flag.String("out", "site/index.html", "output HTML path")
	flag.Parse()

	if err := run(*root, *baseURL, *out); err != nil {
		fmt.Fprintf(os.Stderr, "catalog: %v\n", err)
		os.Exit(1)
	}
}

func run(root, baseURL, out string) error {
	utils, err := catalog.Discover(root)
	if err != nil {
		return err
	}

	index, err := catalog.RenderIndex(baseURL, utils)
	if err != nil {
		return err
	}

	if len(index) == 0 {
		return errors.New("empty index output")
	}

	if err := ensureDir(filepath.Dir(out)); err != nil {
		return err
	}

	if err := os.WriteFile(out, index, 0o644); err != nil {
		return fmt.Errorf("write index: %w", err)
	}

	return nil
}

func ensureDir(dir string) error {
	if dir == "" || dir == "." {
		return nil
	}

	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("create directory %q: %w", dir, err)
	}

	return nil
}
