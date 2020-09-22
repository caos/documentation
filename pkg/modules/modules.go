package modules

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Modules struct {
	basePath string
}

func New(basePath string) *Modules {
	path := strings.TrimPrefix(basePath, filepath.Join(os.ExpandEnv("$GOPATH"), "src")+"/")
	return &Modules{
		basePath: getModuleForPath(path),
	}
}

func (m *Modules) GetPathForImport(importPath string) string {
	if importPath == "" {
		return ""
	}

	if strings.HasPrefix(importPath, m.basePath) {
		return filepath.Join(os.ExpandEnv("$GOPATH"), "src", importPath)
	}

	return filepath.Join(os.ExpandEnv("$GOPATH"), "pkg", "mod", getModuleForPath(importPath))
}

func getModuleForPath(path string) string {
	levels := strings.Split(path, "/")

	for i := len(levels); i > 0; i-- {
		path := filepath.Join(levels[0:i]...)
		mod, err := checkGoList(path)
		if err == nil {
			return strings.Replace(mod, " ", "@", 1)
		}
	}
	return ""
}

func checkGoList(path string) (string, error) {
	cmd := exec.Command("sh", "-c", "go list -m "+path)
	resultPath, err := cmd.CombinedOutput()
	strPath := string(resultPath)
	return strings.TrimSuffix(strPath, "\n"), err
}
