package modules

import (
	"github.com/caos/documentation/pkg/modules/pack"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var cachedModules []*Module

type Module struct {
	basePath       string
	mod            string
	cachedImports  map[string]string
	cachedPackages map[string]*pack.Package
}

func New(basePath, mod string) *Module {
	return &Module{
		basePath:       basePath,
		mod:            mod,
		cachedPackages: map[string]*pack.Package{},
		cachedImports:  map[string]string{},
	}
}

func CachedModule(basePath string) *Module {
	if cachedModules == nil {
		cachedModules = make([]*Module, 0)
	}

	baseRepoPrefix := filepath.Join(os.ExpandEnv("$GOPATH"), "src") + "/"
	modPrefix := filepath.Join(os.ExpandEnv("$GOPATH"), "pkg", "mod") + "/"

	if strings.HasPrefix(basePath, baseRepoPrefix) {
		path := strings.TrimPrefix(basePath, baseRepoPrefix)
		mod, _ := getModuleForPath(path)
		modPath := baseRepoPrefix + mod

		for _, module := range cachedModules {
			if module.basePath == modPath {
				return module
			}
		}
		module := New(modPath, mod)
		cachedModules = append(cachedModules, module)
		return module
	}

	mod, _ := getModuleForPath(basePath)
	modPath := modPrefix + mod

	for _, module := range cachedModules {
		if module.basePath == modPath {
			return module
		}
	}

	module := New(modPath, mod)
	cachedModules = append(cachedModules, module)
	return module
}

func (m *Module) CachePackage(path string) (ret *pack.Package) {

	localPath := strings.TrimPrefix(path, m.basePath)

	for _, cachedPackage := range m.cachedPackages {
		if cachedPackage.ImportPath == localPath {
			return cachedPackage
		}
	}

	ret = pack.New(path, localPath)
	m.cachedPackages[localPath] = ret
	return ret
}

func (m *Module) GetPathForImport(importPath string) string {
	if importPath == "" {
		return ""
	}

	cached, found := m.cachedImports[importPath]
	if found {
		return cached
	}

	resultPath := ""
	if strings.HasPrefix(importPath, m.mod) {
		resultPath = filepath.Join(os.ExpandEnv("$GOPATH"), "src", importPath)
	} else {
		mod, relativePath := getModuleForPath(importPath)
		resultPath = filepath.Join(os.ExpandEnv("$GOPATH"), "pkg", "mod", mod, relativePath)
	}

	m.cachedImports[importPath] = resultPath
	return resultPath
}

func (m *Module) GetModulePath() string {
	return m.basePath
}

func getModuleForPath(path string) (string, string) {
	levels := strings.Split(path, "/")

	for i := len(levels); i > 0; i-- {
		localPath := filepath.Join(levels[0:i]...)
		mod, err := checkGoList(localPath)
		if err == nil {
			return strings.Replace(mod, " ", "@", 1), filepath.Join(levels[i:len(levels)]...)
		}
	}
	return "", ""
}

func checkGoList(path string) (string, error) {
	cmd := exec.Command("sh", "-c", "go list -m "+path)
	resultPath, err := cmd.CombinedOutput()
	strPath := string(resultPath)
	return strings.TrimSuffix(strPath, "\n"), err
}
