//go:build !tinygo

package garage

import "embed"

//go:embed css go.mod html js images template
var fs embed.FS
