package catalog

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLinkFor(t *testing.T) {
	tests := map[string]struct {
		base string
		slug string
		want string
	}{
		"relative": {
			base: "",
			slug: "foo",
			want: "./foo/",
		},
		"absolute": {
			base: "https://example.com",
			slug: "foo",
			want: "https://example.com/foo/",
		},
		"base handles trailing slash": {
			base: "https://example.com/",
			slug: "foo",
			want: "https://example.com/foo/",
		},
		"slug trims separators": {
			base: "https://example.com",
			slug: "/foo/bar/",
			want: "https://example.com/foo/bar/",
		},
		"root relative": {
			base: "",
			slug: "",
			want: "./",
		},
		"root absolute": {
			base: "https://example.com/utilities/",
			slug: "",
			want: "https://example.com/utilities/",
		},
	}

	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			r := require.New(t)
			got := linkFor(tc.base, tc.slug)
			r.Equal(tc.want, got)
		})
	}
}

func TestDiscover(t *testing.T) {
	r := require.New(t)

	root := t.TempDir()
	r.NoError(os.MkdirAll(filepath.Join(root, "util-one"), 0o755))
	r.NoError(os.MkdirAll(filepath.Join(root, "util-two"), 0o755))

	writeMeta := func(dir, name, description string) {
		meta := []byte("name: " + name + "\ndescription: " + description + "\n")
		r.NoError(os.WriteFile(filepath.Join(dir, MetadataFileName), meta, 0o644))
	}

	writeMeta(filepath.Join(root, "util-one"), "One", "first util")
	writeMeta(filepath.Join(root, "util-two"), "Two", "second util")

	utils, err := Discover(root)
	r.NoError(err)
	r.Len(utils, 2)
	r.Equal("One", utils[0].Name)
	r.Equal("util-one", utils[0].Slug)
	r.Equal("Two", utils[1].Name)
	r.Equal("util-two", utils[1].Slug)
}

func TestUtilityValidate(t *testing.T) {
	tests := map[string]struct {
		util        Utility
		defaultSlug string
		wantErr     bool
		wantSlug    string
	}{
		"fills default slug": {
			util: Utility{
				Name:        "One",
				Description: "desc",
			},
			defaultSlug: "util-one",
			wantSlug:    "util-one",
		},
		"keeps provided slug": {
			util: Utility{
				Name:        "One",
				Description: "desc",
				Slug:        "custom",
			},
			defaultSlug: "ignored",
			wantSlug:    "custom",
		},
		"missing name": {
			util: Utility{
				Description: "desc",
			},
			wantErr: true,
		},
		"missing description": {
			util: Utility{
				Name: "One",
			},
			wantErr: true,
		},
	}

	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			r := require.New(t)
			err := tc.util.Validate(tc.defaultSlug)
			if tc.wantErr {
				r.Error(err)
				return
			}

			r.NoError(err)
			r.Equal(tc.wantSlug, tc.util.Slug)
		})
	}
}
