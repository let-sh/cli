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

func (c *DeployContext)DetectProjectType() (projectType string) {
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
				c.Type = "gin"
				return "gin"
			}
			if v.Mod.Path == "github.com/go-martini/martini" {
				c.Type = "martini"
				return "martini"
			}
		}

		return "go"
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
				c.Type = "docusaurus"
				return "docusaurus"
			}

			if strings.Contains(k, "next") {
				c.Type = "next"
				return "next"
			}

			if strings.Contains(k, "react") {
				c.Type = "react"
				return "react"
			}

			if strings.Contains(k, "@nuxt") {
				c.Type = "nuxt"
				return "nuxt"
			}
			if strings.Contains(k, "@vue") || strings.Contains(k, "vue") {
				c.Type = "vue"
				return "vue"
			}
			if strings.Contains(k, "@angular") {
				c.Type = "angular"
				return "angular"
			}

		}
	}

	// handle static files
	// check if static by index.html
	_, err := os.Stat("index.html")
	if !os.IsNotExist(err) {
		c.Type = "static"
		c.Static = "./"
		return
	}

	// handle static site generator
	// hexo

	// hugo
	if FileExists("config.toml") && FileExists("themes") {
		c.Type = "hugo"
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
