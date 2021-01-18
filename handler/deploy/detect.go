package deploy

import (
	"encoding/json"
	"github.com/let-sh/cli/log"
	"github.com/let-sh/cli/types"
	"github.com/rogpeppe/go-internal/modfile"
	"io/ioutil"
	"os"
	"strings"
)

func DetectProjectType() (projectType string) {
	// handle golang framework
	if FileExists("go.mod") {
		b, err := ioutil.ReadFile("go.mod")

		if err != nil {
			log.Error(err)
		}

		f, err := modfile.Parse("go.mod", b, nil)
		if err != nil {
			log.Error(err)
		}
		for _, v := range f.Require {
			if v.Mod.Path == "github.com/gin-gonic/gin" {
				return "gin"
			}
			if v.Mod.Path == "github.com/go-martini/martini" {
				return "martini"
			}
		}
	}

	// handle javascript framework
	if FileExists("package.json") {
		// frontend
		jsonBytes, err := ioutil.ReadFile("package.json")
		if err != nil {
			log.Error(err)
		}
		var packageConfig types.PackageDotJson
		json.Unmarshal(jsonBytes, &packageConfig)

		for k, _ := range packageConfig.Dependencies {
			if strings.Contains(k, "@docusaurus") {
				return "docusaurus"
			}

			if strings.Contains(k, "next") {
				return "react"
			}

			if strings.Contains(k, "react") {
				return "react"
			}

			if strings.Contains(k, "@nuxt") {
				return "nuxt"
			}
			if strings.Contains(k, "@vue") {
				return "vue"
			}
			if strings.Contains(k, "@angular") {
				return "angular"
			}

		}
	}
	// handle static files

	// handle static site generator
	// hexo

	// hugo
	if FileExists("config.toml") && FileExists("themes") {
		return "hugo"
	}
	return "static"
}

func FileExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return false
	}
	return true
}
