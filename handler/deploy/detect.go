package deploy

import (
	"encoding/json"
	"github.com/let-sh/cli/log"
	"github.com/let-sh/cli/types"
	"github.com/pelletier/go-toml"
	"github.com/rogpeppe/go-internal/modfile"
	"io/ioutil"
	"os"
	"strings"
)

func (c *DeployContext) DetectProjectType() (projectType string) {
	// handle rust project
	if FileExists("Cargo.toml") {
		content, err := ioutil.ReadFile("Cargo.toml")

		if err != nil {
			log.Error(err)
			return
		}
		cargo, err := toml.Load(string(content))
		if err != nil {
			log.Error(err)
			return
		}

		dependencies := cargo.Get("dependencies").(*toml.Tree)

		if dependencies.Has("rocket") {
			c.Type = "rocket"
			return "rocket"
		}

	}

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

		for k := range packageConfig.Dependencies {
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

	// handle python framework
	if FileExists("requirements.txt") || FileExists("Pipfile") {
		if FileExists("requirements.txt") {
			reqFile, _ := ioutil.ReadFile("requirements.txt")
			if strings.Contains(strings.ToLower(string(reqFile)), "flask") {
				c.Type = "flask"
				return "flask"
			}
		}

		if FileExists("Pipfile") {
			reqFile, _ := ioutil.ReadFile("Pipfile")
			if strings.Contains(strings.ToLower(string(reqFile)), "flask") {
				c.Type = "flask"
				return "flask"
			}
		}
	}

	// handle static files
	// check if static by index.html
	_, err := os.Stat("index.html")
	if !os.IsNotExist(err) {
		c.Type = "static"
		c.Static = "./"
		return "static"
	}

	// handle static site generator
	// TODO: hexo

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
