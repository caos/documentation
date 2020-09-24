package pack

import (
	"github.com/caos/documentation/pkg/treeelement"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Package struct {
	BasePath       string
	ImportPath     string
	CachedElements map[string]*treeelement.TreeElement
}

func New(basePath, importPath string) *Package {
	path := os.ExpandEnv(basePath)

	return &Package{
		BasePath:       path,
		ImportPath:     importPath,
		CachedElements: map[string]*treeelement.TreeElement{},
	}
}

func (p *Package) GetGoFileList() []string {
	files, err := getFilesInDirectory(p.BasePath)
	if err != nil {
		return nil
	}

	goFiles := make([]string, 0)
	for _, file := range files {
		if !strings.HasSuffix(file, ".go") {
			continue
		}
		goFiles = append(goFiles, file)
	}
	return goFiles
}

func getFilesInDirectory(dirPath string) ([]string, error) {
	files := make([]string, 0)

	infos, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	for _, info := range infos {
		if !info.IsDir() {
			files = append(files, filepath.Join(dirPath, info.Name()))
		}
	}

	return files, err
}
