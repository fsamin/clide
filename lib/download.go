package clide

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/graymeta/stow"
)

//DownloadFiles download files on a stow.Location. It may creates the destination if not exists
func DownloadFiles(location stow.Location, destination string, files []string, logger ProgressPrinter) ([]File, error) {
	os.MkdirAll(destination, os.FileMode(0755))
	res := []File{}
	for _, f := range files {
		var container stow.Container
		var found bool
		containers, _, err := location.Containers(f, stow.CursorStart, 10000)
		if err != nil {
			return nil, err
		}
		for _, c := range containers {
			if c.Name() == f {
				container = c
				found = true
				break
			}
		}

		if !found {
			return nil, fmt.Errorf("Container %s not found", f)
		}

		if err := stow.Walk(container, stow.NoPrefix, 10000,
			func(item stow.Item, err error) error {
				if err != nil {
					return err
				}
				filename := filepath.Join(destination, item.Name())
				url := URL(item)
				res = append(res, File{
					Filename: filename,
					URL:      url,
				})
				r, err := item.Open()
				if err != nil {
					return err
				}
				defer r.Close()
				f, err := os.Create(filename)
				if err != nil {
					return err
				}
				defer f.Close()
				if logger != nil {
					logger("Downloading %s...\n", url)
				}
				if _, err := io.Copy(f, r); err != nil {
					return err
				}
				return nil
			}); err != nil {
			return nil, err
		}
	}

	return res, nil
}
