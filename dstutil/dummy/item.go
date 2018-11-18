package dummy

import (
	"go/format"
	"io"
	"path/filepath"
	"strings"

	"gopkg.in/src-d/go-billy.v4"
)

type Item interface {
	create(dir, name string, fs billy.Filesystem)
}
type Src string
type Dir map[string]Item

func (s Src) create(dir, name string, fs billy.Filesystem) {
	src := string(s)
	if strings.HasSuffix(name, ".go") {
		// if it's a go file, format first
		b, err := format.Source([]byte(s))
		if err != nil {
			panic(err)
		}
		src = string(b)
	}
	if err := fs.MkdirAll(dir, 0777); err != nil {
		panic(err)
	}
	f, err := fs.Create(filepath.Join(dir, name))
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if _, err := io.Copy(f, strings.NewReader(src)); err != nil {
		panic(err)
	}
}
func (d Dir) create(dir, name string, fs billy.Filesystem) {
	dpath := filepath.Join(dir, name)
	if err := fs.MkdirAll(dpath, 0777); err != nil {
		panic(err)
	}
	for name, item := range d {
		item.create(dpath, name, fs)
	}
}
