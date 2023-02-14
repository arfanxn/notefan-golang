package config

type fsDisk struct {
	Root string // the root of the filesystem
	URL  string // the URL to access file
}

// FSDisks (FileSystemDisks) represents a disk configuration
var FSDisks = map[string]fsDisk{
	"public": {
		Root: "./public/",
		URL:  "",
	},
}
