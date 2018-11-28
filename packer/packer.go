package packer

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"

	"github.com/lassik/airfreight"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

// An EntPackage represents one .go file while packing.
type EntPackage struct {
	name string
	maps map[string]map[string]airfreight.Ent
}

// Start packing static files into the given Go package.
//
// It's ok to have other stuff besides your static files in the same
// package. The generated code doesn't pollute the package namespace.
// It only defines the map names you give to .Map().
func Package(name string) EntPackage {
	return EntPackage{name: name, maps: map[string]map[string]airfreight.Ent{}}
}

func mapDir(entMap map[string]airfreight.Ent, rootDir, relDir string) {
	files, err := ioutil.ReadDir(path.Join(rootDir, relDir))
	check(err)
	for _, info := range files {
		rel := info.Name()
		if rel[0] == '.' {
			continue
		}
		rel = relDir + "/" + rel
		if info.IsDir() {
			mapDir(entMap, rootDir, rel)
		} else {
			bytes, err := ioutil.ReadFile(path.Join(rootDir, rel))
			check(err)
			entMap[rel] = airfreight.Ent{
				ModTime:  info.ModTime().Unix(),
				Contents: string(bytes),
			}
		}
	}
}

// Add one Go map into the package.
//
// The map will contain all regular files from the given rootDirs and
// their subdirectories. However, subdirectories and files whose names
// start with a dot are ignored.
//
// Filenames inside the map do not remember the rootDir from which
// they came, so if the rootDir "static" contains the file "foo.js",
// the filename inside the map will be just "/foo.js" instead of
// "/static/foo.js".
//
// The map keys are filenames, and the map values are airfreight.Ent
// structures.
func (p EntPackage) Map(mapName string, rootDirs ...string) EntPackage {
	p.maps[mapName] = map[string]airfreight.Ent{}
	for _, rootDir := range rootDirs {
		mapDir(p.maps[mapName], rootDir, "")
	}
	return p
}

// Write a Go source file for the package into the given Writer.
func (p EntPackage) WriteTo(w io.Writer) (int64, error) {
	var len int64
	n, err := fmt.Fprintf(w, "// @generated-by airfreight\n\n"+
		"package %s\n\n"+
		"import \"github.com/lassik/airfreight\"\n", p.name)
	len += int64(n)
	if err != nil {
		return len, err
	}
	for mapVar, mapEnts := range p.maps {
		n, err = fmt.Fprintf(w,
			"\nvar %s = map[string]airfreight.Ent{\n", mapVar)
		len += int64(n)
		if err != nil {
			return len, err
		}
		for entName, ent := range mapEnts {
			n, err = fmt.Fprintf(w, "\n\t%#v: airfreight.Ent{"+
				"ModTime: %#v,"+
				" Contents: %#v},\n",
				entName, ent.ModTime, ent.Contents)
			len += int64(n)
			if err != nil {
				return len, err
			}
		}
		n, err = fmt.Fprintf(w, "}\n")
		len += int64(n)
		if err != nil {
			return len, err
		}
	}
	return len, err
}

// Write a Go source file for the package into the given .go file.
//
// NOTE: The file will be overwritten if it exists.
func (p EntPackage) WriteFile(filename string) {
	file, err := os.Create(filename)
	check(err)
	defer file.Close()
	p.WriteTo(file)
}
