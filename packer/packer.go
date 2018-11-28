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

type EntPackage struct {
	name string
	maps map[string]map[string]airfreight.Ent
}

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
		if relDir != "" {
			rel = relDir + "/" + rel
		}
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

func (p EntPackage) Map(mapName string, rootDirs ...string) EntPackage {
	p.maps[mapName] = map[string]airfreight.Ent{}
	for _, rootDir := range rootDirs {
		mapDir(p.maps[mapName], rootDir, "")
	}
	return p
}

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
			"\nvar %s = map[string]airfreight.Ent{\n\n", mapVar)
		len += int64(n)
		if err != nil {
			return len, err
		}
		for entName, ent := range mapEnts {
			n, err = fmt.Fprintf(w, "\t%#v: airfreight.Ent{"+
				"ModTime: %#v,"+
				" Contents: %#v},\n\n",
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

func (p EntPackage) WriteFile(filename string) {
	file, err := os.Create(filename)
	check(err)
	defer file.Close()
	p.WriteTo(file)
}
