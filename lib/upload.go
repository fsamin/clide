package clide

import (
	"bytes"
	"io/ioutil"
	"path/filepath"

	"github.com/graymeta/stow"
)

//UploadFiles upload files on a stow.Location. It may creates the destination if not exists
func UploadFiles(location stow.Location, destination string, files []string, logger ProgressPrinter) ([]File, error) {
	var container stow.Container
	containers, _, err := location.Containers(destination, stow.CursorStart, 10000)
	if err != nil {
		return nil, err
	}
	for _, c := range containers {
		if c.Name() == destination {
			container = c
			break
		}
	}

	if container == nil {
		container, err = location.CreateContainer(destination)
		if err != nil {
			return nil, err
		}
	}

	res := []File{}
	for i := range files {
		if logger != nil {
			logger("Uploading %s...\n", files[i])
		}
		btes, err := ioutil.ReadFile(files[i])
		if err != nil {
			return nil, err
		}
		name := filepath.Base(files[i])
		item, err := container.Put(name, bytes.NewBuffer(btes), int64(len(btes)), map[string]interface{}{})
		if err != nil {
			return nil, err
		}

		url := URL(item)

		res = append(res, File{
			Filename: files[i],
			URL:      url,
		})

	}

	return res, nil
}
