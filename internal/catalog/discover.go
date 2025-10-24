package catalog

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"

	"gopkg.in/yaml.v3"
)

var errSkip = errors.New("skip")

// Discover loads all utilities that contain a metadata file beneath root.
func Discover(root string) ([]Utility, error) {
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, fmt.Errorf("read root: %w", err)
	}

	var utils []Utility

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		dir := filepath.Join(root, entry.Name())
		util, err := loadUtility(dir, entry.Name())
		if err != nil {
			if errors.Is(err, errSkip) {
				continue
			}

			return nil, fmt.Errorf("load utility %q: %w", entry.Name(), err)
		}

		utils = append(utils, util)
	}

	sort.Slice(utils, func(i, j int) bool {
		return utils[i].Name < utils[j].Name
	})

	return utils, nil
}

func loadUtility(dir, defaultSlug string) (Utility, error) {
	metaFile := filepath.Join(dir, MetadataFileName)
	f, err := os.Open(metaFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Utility{}, errSkip
		}

		return Utility{}, fmt.Errorf("open metadata: %w", err)
	}
	defer func() {
		_ = f.Close()
	}()

	util, err := decodeUtility(f)
	if err != nil {
		return Utility{}, fmt.Errorf("decode metadata: %w", err)
	}

	if err := util.Validate(defaultSlug); err != nil {
		return Utility{}, fmt.Errorf("validate metadata: %w", err)
	}

	return util, nil
}

func decodeUtility(r io.Reader) (Utility, error) {
	var util Utility

	dec := yaml.NewDecoder(r)
	if err := dec.Decode(&util); err != nil {
		return Utility{}, err
	}

	return util, nil
}
