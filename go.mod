module navidrome-achievements

go 1.25

// These two requires point at the Go PDK packages vendored inside a
// checkout of github.com/navidrome/navidrome. Clone navidrome next to
// this plugin (or adjust the relative path) before building:
//
//   git clone https://github.com/navidrome/navidrome ../navidrome
//
require github.com/navidrome/navidrome/plugins/pdk/go v0.0.0

require github.com/extism/go-pdk v1.1.3

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	github.com/stretchr/testify v1.11.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/navidrome/navidrome/plugins/pdk/go => ../navidrome/plugins/pdk/go
