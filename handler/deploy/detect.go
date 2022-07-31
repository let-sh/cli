package deploy

import (
    "encoding/json"
    "github.com/let-sh/cli/log"
    "github.com/let-sh/cli/types"
    "github.com/let-sh/cli/utils"
    "github.com/pelletier/go-toml"
    "github.com/rogpeppe/go-internal/modfile"
    "golang.org/x/exp/maps"
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
        if packageConfig.Dependencies == nil {
            packageConfig.Dependencies = map[string]string{}
        }
        if packageConfig.DevDependencies == nil {
            packageConfig.DevDependencies = map[string]string{}
        }
        maps.Copy(packageConfig.Dependencies, packageConfig.DevDependencies)
        packages := utils.Keys(packageConfig.Dependencies)
        if utils.ItemExists(packages, "@surgio/gateway") {
            c.Type = "gateway"
            return "gateway"
        }

        if utils.ItemExists(packages, "express") {
            c.Type = "express"
            return "express"
        }

        if utils.ItemExists(packages, "@docusaurus") {
            c.Type = "docusaurus"
            return "docusaurus"
        }

        if utils.ItemExists(packages, "next") {
            c.Type = "next"
            return "next"
        }

        if utils.ItemExists(packages, "vitepress") {
            c.Type = "vitepress"
            return "vitepress"
        }

        if utils.ItemExists(packages, "vuepress") {
            c.Type = "vuepress"
            return "vuepress"
        }

        if utils.ItemExists(packages, "react") {
            c.Type = "react"
            return "react"
        }

        if utils.ItemExists(packages, "@nuxt") {
            c.Type = "nuxt"
            return "nuxt"
        }
        if utils.ItemExists(packages, "@vue") || utils.ItemExists(packages, "vue") {
            c.Type = "vue"
            return "vue"
        }
        if utils.ItemExists(packages, "@angular") {
            c.Type = "angular"
            return "angular"
        }

        if utils.ItemExists(packages, "hexo") {
            c.Type = "hexo"
            return "hexo"
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

            if strings.Contains(strings.ToLower(string(reqFile)), "fastapi") {
                c.Type = "fastapi"
                return "fastapi"
            }
        }

        if FileExists("Pipfile") {
            reqFile, _ := ioutil.ReadFile("Pipfile")
            if strings.Contains(strings.ToLower(string(reqFile)), "flask") {
                c.Type = "flask"
                return "flask"
            }

            if strings.Contains(strings.ToLower(string(reqFile)), "fastapi") {
                c.Type = "fastapi"
                return "fastapi"
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

    // hugo
    if FileExists("config.toml") && FileExists("themes") {
        c.Type = "hugo"
        return "hugo"
    }

    c.Type = "unknown"
    return "unknown"
}

func FileExists(path string) bool {
    _, err := os.Stat(path)
    if err != nil {
        return false
    }
    return true
}

type m = map[string]string
