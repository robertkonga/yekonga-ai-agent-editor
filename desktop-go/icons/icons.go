package icons

import (
	"embed"
	"io/fs"
	"path/filepath"
	"strings"
)

//go:embed icons/*
var assets embed.FS

// this is icon on project directory folder of file depend on the language and framework
type Icon struct {
	Name      string `json:"name"`
	Path      string `json:"path"`
	Language  string `json:"language"`
	Extension string `json:"extension"`
	Framework string `json:"framework"`
}

func ListOfIcons() ([]Icon, error) {
	var list []Icon

	err := fs.WalkDir(assets, "icons", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			ext := filepath.Ext(d.Name())
			name := strings.TrimSuffix(d.Name(), ext)
			list = append(list, Icon{
				Name:      name,
				Path:      path,
				Extension: ext,
			})
		}
		return nil
	})

	return list, err
}
