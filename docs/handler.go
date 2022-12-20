// Forked from
// https://github.com/tendermint/starport/blob/v0.19.2-alpha/starport/pkg/openapiconsole/console.go

package docs

import (
	"embed"
	"html/template"
	"net/http"
)

//go:embed index.tpl
var index embed.FS

// Handler returns an http handler that servers OpenAPI console for an OpenAPI
// spec at specURL.
func Handler(title, specURL string) http.HandlerFunc {
	t, err := template.ParseFS(index, "index.tpl")
	if err != nil {
		panic(err)
	}

	return func(w http.ResponseWriter, req *http.Request) {
		err := t.Execute(w, struct {
			Title string
			URL   string
		}{
			title,
			specURL,
		})
		if err != nil {
			panic(err) // this might not be the right way to check this error.
		}
	}
}
