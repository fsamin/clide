package clide

import (
	"strings"

	"github.com/graymeta/stow"
)

//URL return the URL string of an item
func URL(item stow.Item) string {
	url := strings.Replace(item.URL().String(), "swift://", "https://", 1)
	url = strings.Replace(url, "s3://", "", 1)
	return url
}
