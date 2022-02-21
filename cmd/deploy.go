/*
Copyright © 2021 Fred Liang <fred@oasis.ac>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/muesli/termenv"
	"io/ioutil"
	"os"
	"os/exec"
	"os/signal"
	"os/user"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/atotto/clipboard"
	"github.com/c2h5oh/datasize"
	"github.com/let-sh/cli/handler/deploy"
	"github.com/let-sh/cli/info"
	"github.com/let-sh/cli/log"
	"github.com/let-sh/cli/requests"
	"github.com/let-sh/cli/requests/graphql"
	"github.com/let-sh/cli/types"
	"github.com/let-sh/cli/utils"
	"github.com/let-sh/cli/utils/cache"
	"github.com/let-sh/cli/utils/s3"
	"github.com/manifoldco/promptui"
	"github.com/mholt/archiver/v3"
	c "github.com/otiai10/copy"
	ignore "github.com/sabhiram/go-gitignore"
	gql "github.com/shurcooL/graphql"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var DeploymentID string

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy your current project to let.sh",
	Long:  `Deploy your current project to let.sh with a single command line`,
	Run: func(cmd *cobra.Command, args []string) {
		// Setup our Ctrl+C handler
		SetupCloseHandler()

		// check whether user is logged in
		if info.Credentials.LoadToken() == "" {
			log.Warning("please login via `lets login` first")
			return
		}

		log.BStart("deploying")

		// check deployment directory is valid
		{
			// check not home dir
			dir, _ := os.Getwd()
			usr, _ := user.Current()
			if dir == usr.HomeDir {
				log.Error(errors.New("currently under home dir, please switch to your project dir"))
				return
			}

			// limit files to 10000
			files, _ := ioutil.ReadDir("./")
			if len(files) > 10000 {
				log.Error(errors.New("too many files in current dir, please check whether in the " +
					"correct directory"))
				return
			}
		}

		// merge config
		// cli flag > config file > auto saved config > detected config & types
		{
			// detect current project config first
			// init current project name
			dir, _ := os.Getwd()
			deploymentCtx.Name = filepath.Base(dir)

			// detect current project info
			deploymentCtx.DetectProjectType()

			// Step2: get cache config
			deploymentCtx.LoadProjectInfoCache()

			// Step3: load user config and environment variables
			deploymentCtx.LoadLetJson()
			err := deploymentCtx.LoadEnvFiles()
			if err != nil {
				log.Error(err)
				return
			}

			// Step4: merge cli flag config
			deploymentCtx.LoadCliFlag(inputProjectName, inputProjectType)
			deploymentCtx.Web3 = &inputWeb3

			// load cn
			// if user customized cn flag
			deploymentCtx.LoadRegion(cmd, inputCN)
		}

		// check project exists
		// if not exists, tell to create
		// and confirm project configuration
		if !deploymentCtx.ConfirmProject() {
			log.Warning("deploy canceled")
			return
		}

		// check if user dir is changed
		if _, ok := cache.ProjectsInfo[deploymentCtx.Name]; ok {
			pwd, _ := os.Getwd()

			logrus.Debug("cached project dir: ", cache.ProjectsInfo[deploymentCtx.Name].Dir)
			logrus.Debug("current project dir: ", pwd)

			if pwd != cache.ProjectsInfo[deploymentCtx.Name].Dir {
				log.S.StopFail()

				if !inputAssumeYes {
					// if current dir is not previous dir
					prompt := promptui.Prompt{
						Label:   "Detected your project dir changed, continue deployment?[Y/n]",
						Default: "Y",
						Validate: func(input string) error {
							if utils.ItemExists([]string{"", "n", "N", "No", "Y", "y", "yes", "Yes"}, input) {
								return nil
							}
							return errors.New("no matching input")
						},
					}

					result, err := prompt.Run()
					if err != nil {
						log.Error(err)
						return
					}

					if utils.ItemExists([]string{"n", "N", "No"}, result) {
						log.S.StopFail()
						log.Warning("Deployment canceled")
						return
					}
				}

				log.BStart("deploying")
			}
		}

		// make pre deploy request
		{
			query := struct {
				graphql.QueryCheckDeployCapability
				graphql.QueryBuildTemplate
				graphql.QueryStsToken
				graphql.QueryPreference
			}{}
			err := graphql.NewClient().Query(context.Background(), &query, map[string]interface{}{
				"projectName": gql.String(deploymentCtx.Name),
				"tokenType":   gql.String("buildBundle"),
				"type":        gql.String(deploymentCtx.Type),
				"cn":          gql.Boolean(*deploymentCtx.CN),
				"name":        gql.String("channel"),
			})

			var requestError *gql.RequestError
			var graphqlError *gql.GraphQLError
			if errors.As(err, &graphqlError) {
				log.Error(errors.New(graphqlError.GraphqlErrors[0].Message))
				return
			}

			if errors.As(err, &requestError) {
				log.Error(requestError)
				return
			}

			deploymentCtx.PreDeployRequest = query

			// check whether is dynamic project
			if deploymentCtx.LetConfig.Web3 != nil {
				if deploymentCtx.PreDeployRequest.BuildTemplate.ContainsDynamic && *deploymentCtx.LetConfig.Web3 {
					log.Warning("you cannot deploy dynamic project to web3 infra yet")
					return
				}
			}

		}

		// get project type config from api
		log.S.StopFail()

		fmt.Println("")
		fmt.Println(log.CyanBold("Detected Project Info"))
		fmt.Println("name:", termenv.String(deploymentCtx.Name).Bold().String())
		fmt.Println("type:", termenv.String(deploymentCtx.Type).Bold().String())
		fmt.Println("")

		{
			if deploymentCtx.Static == "" {
				deploymentCtx.Static = deploymentCtx.PreDeployRequest.BuildTemplate.DistDir
			}
		}

		// if contains static, upload static files to s3
		dirPath := deploymentCtx.Static
		if len(dirPath) == 0 {
			dirPath = "./"
		}
		if deploymentCtx.PreDeployRequest.BuildTemplate.ContainsStatic {
			if utils.ItemExists([]string{"static"}, deploymentCtx.Type) {
				// todo: merge static dir value source
				if err := s3.UploadDirToStaticSource(dirPath, deploymentCtx.Name, deploymentCtx.Name+"-"+deploymentCtx.PreDeployRequest.CheckDeployCapability.HashID, *deploymentCtx.CN); err != nil {
					log.Error(err)
					return
				}
			} else {
				if deploymentCtx.PreDeployRequest.BuildTemplate.LocalCompiling {
					for _, command := range deploymentCtx.PreDeployRequest.BuildTemplate.CompileCommands {
						command := strings.Split(command, " ")
						c := exec.Command(command[0], command[1:]...)
						c.Stdout = os.Stdout
						c.Stderr = os.Stderr
						err := c.Run()
						if err != nil {
							log.Error(err)
							return
						}
					}
				}

				if err := s3.UploadDirToStaticSource(deploymentCtx.Static, deploymentCtx.Name, deploymentCtx.Name+"-"+deploymentCtx.PreDeployRequest.CheckDeployCapability.HashID, *deploymentCtx.CN); err != nil {
					log.Error(err)
					return
				}
			}
		}

		// if contains dynamic, upload dynamic files to s3
		// then trigger deployment
		if deploymentCtx.PreDeployRequest.BuildTemplate.ContainsDynamic {
			//
			//// create temp dir
			//dir := os.TempDir()
			//
			//defer os.RemoveAll(dir)
			////fmt.Println(dir)
			////os.MkdirAll(dir+"/source", os.ModePerm)
			//
			//// copy current dir to temp dir
			//c.Copy("./", dir+"/"+deploymentCtx.Name+"-"+hashID)
			dirPath, _ = os.Getwd()

			// remove if not clean
			//os.Remove(dir + "/" + deploymentCtx.Name + "-" + hashID + ".tar.gz")

			// Read directory files
			var names []string
			err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
				if !info.IsDir() {
					names = append(names, path)
				}
				return nil
			})
			if err != nil {
				log.Error(err)
				return
			}

			// respect .gitignore and .letignore
			if _, err := os.Stat(filepath.Join(dirPath, ".gitignore")); err == nil {
				// match a file against a particular .gitignore
				i, _ := ignore.CompileIgnoreFile(filepath.Join(dirPath, ".gitignore"))

				tmp := []string{}
				for _, v := range names {
					if !i.MatchesPath(v) {
						tmp = append(tmp, v)
					}
				}

				names = tmp
			}

			// .letignore
			if _, err := os.Stat(filepath.Join(dirPath + ".letignore")); err == nil {
				// match a file against a particular .gitignore
				i, _ := ignore.CompileIgnoreFile(filepath.Join(dirPath + ".letignore"))

				tmp := []string{}
				for _, v := range names {
					if !i.MatchesPath(v) {
						tmp = append(tmp, v)
					}
				}
				names = tmp
			}

			// calculate files size
			if size, err := utils.GetFilesSize(names); err != nil {
				log.Error(err)
				return
			} else {
				// source code is too big
				// < 20 MB directly upload
				// 20 MB <= files < 40 MB confirm
				// >= 40 MB abort
				if uint64(size) > 20*datasize.MB.Bytes() {
					log.Warning(`your directory is too big, larger than 20 MB.
you could remove the irrelevant via .letignore or gitignore.`)
					return
				}
			}

			// copy to temp dir
			tempDir, _ := ioutil.TempDir("", "upload")
			defer os.RemoveAll(tempDir)
			for _, f := range names {
				toName := strings.Replace(f, dirPath, tempDir+"/", 1)
				err := c.Copy(f, toName)
				if err != nil {
					log.Error(err)
				}
			}

			// copy to temp dir
			tempZipDir, _ := ioutil.TempDir("", "zip")
			defer os.RemoveAll(tempZipDir)

			// switch dir
			os.Chdir(tempDir) // switch to temp directory
			err = archiver.Archive([]string{"."}, tempZipDir+"/"+deploymentCtx.Name+"-"+deploymentCtx.PreDeployRequest.CheckDeployCapability.HashID+".tar.gz")
			os.Chdir(dirPath) // switch back

			if err != nil {
				log.Error(err)
				return
			}
			s3.UploadFileToCodeSource(tempZipDir+"/"+deploymentCtx.Name+"-"+deploymentCtx.PreDeployRequest.CheckDeployCapability.HashID+".tar.gz", deploymentCtx.Name+"-"+deploymentCtx.PreDeployRequest.CheckDeployCapability.HashID+".tar.gz", deploymentCtx.Name, *deploymentCtx.CN)
		}

		logrus.WithFields(logrus.Fields{
			"json": deploymentCtx,
		}).Debugln("deploymentCtx")

		configBytes, _ := json.Marshal(deploymentCtx)

		// determine which channel to deploy
		channel := deploymentCtx.PreDeployRequest.Preference

		if inputProd == true { // if manually set to deploy to production, rewrite channel
			channel = "prod"
		}
		if inputDev == true { // if manually set to deploy to production, rewrite channel
			channel = "dev"
		}

		var deployment struct {
			ID           string `json:"id"`
			TargetFQDN   string `json:"targetFQDN"`
			NetworkStage string `json:"networkStage"`
			PackerStage  string `json:"packerStage"`
			Status       string `json:"status"`
			Project      struct {
				ID string `json:"id"`
			} `json:"project"`
		}
		var err error
		if inputCheckRunID == 0 {
			deployment, err = requests.Deploy(deploymentCtx.Type, deploymentCtx.Name, string(configBytes), channel, *deploymentCtx.CN)
		} else {
			deployment, err = requests.DeployWithCheckRunID(deploymentCtx.Type, deploymentCtx.Name, string(configBytes), channel,
				*deploymentCtx.CN, inputCheckRunID)
		}

		if err != nil {
			log.Error(err)
			return
		}

		if inputDetach {
			log.S.StopFail()
			log.Success("triggered deployment succeeded")
			return
		}

		DeploymentID = deployment.ID

		pwd, _ := os.Getwd()

		// save deployment info
		cache.SaveProjectInfo(types.Project{
			ID:           deployment.Project.ID,
			Name:         deploymentCtx.Name,
			Dir:          pwd,
			Type:         deploymentCtx.Type,
			ServeCommand: cache.ProjectsInfo[deploymentCtx.Name].ServeCommand,
		})

		log.BStart("deploying")

		// awaiting deployment result
		for {
			currentStatus, err := requests.GetDeploymentStatus(deployment.ID)

			// handle logging
			if err != nil {
				log.Error(err)
			}

			// stages:
			// * Queuing
			// * Building
			switch currentStatus.Status {
			case "Queuing":
				log.BUpdate("queuing")
			case "Running":
				if currentStatus.PackerStage == "Build" {
					log.BUpdate("building")
				}
			}
			//log.BUpdate(" NetworkStage: " + currentStatus.NetworkStage + ", PackerStage: " + currentStatus.PackerStage + ", Status: " + currentStatus.Status)

			if currentStatus.Done {
				if currentStatus.Status == "Failed" {
					log.Error(errors.New("build logs: " + currentStatus.ErrorLogs))
					break
				}
				// write review url to clipboard

				writeClipBoardError := clipboard.WriteAll("https://" + currentStatus.TargetFQDN)

				log.S.StopFail()
				//fmt.Println(
				//	color.New(color.Bold).Sprint("Preview: ")+color.New(color.Underline).Sprint("https://"+currentStatus.TargetFQDN)+func() string {
				//		if writeClipBoardError == nil {
				//			return color.New().Sprint("  (📋 Copied!)")
				//		}
				//		return ""
				//	}(), "\n"+
				//		color.New(color.Bold).Sprint("Details: ")+color.New(color.Underline).Sprint("https://let.sh/console/project/"+deploymentCtx.Name+"/details"),
				//)

				// if web3
				if currentStatus.Web3 != nil {
					log.CyanBold("Web3 Info:")
					fmt.Println("IPFS:   ", termenv.String("https://ipfs.io/ipfs/"+currentStatus.Web3.IpfsCID).
						Underline().Bold().
						String())
					fmt.Println("Arweave:", termenv.String("https://arweave.net/"+currentStatus.Web3.ArTID).
						Underline().
						Bold().
						String())
					fmt.Println("\n")
				}

				fmt.Println(

					termenv.String("URL:   ").String(), termenv.String("https://"+currentStatus.
						TargetFQDN).Underline().Bold().String()+func() string {
						if writeClipBoardError == nil {
							p := termenv.ColorProfile()
							return termenv.String("  (📋 Copied!)").Foreground(p.Color("#808080")).String()
						}
						return ""
					}(),
					"\n"+termenv.String("Details: ").String()+termenv.String("https://let."+
						"sh/console/projects/"+deploymentCtx.Name+"/details").Bold().Underline().String(),
					//color.New(color.Bold).Sprint("Preview: ")+color.New(color.Underline).Sprint("https://"+currentStatus.TargetFQDN)+func() string {
					//	if writeClipBoardError == nil {
					//		return color.New().Sprint("  (📋 Copied!)")
					//	}
					//	return ""
					//}(),
					//"\n"+
					//	color.New(color.Bold).Sprint("Details: ")+color.New(color.Underline).Sprint("https://let.sh/console/project/"+deploymentCtx.Name+"/details"),
				)

				break

			}
			// interval
			time.Sleep(time.Second)
		}
		return
	},
}

var deploymentCtx deploy.DeployContext
var inputProjectName string
var inputProjectType string
var inputCN bool
var inputStaticDir string
var inputProd bool
var inputDev bool
var inputDetach bool      // return immediately after submitted deployment
var inputAssumeYes bool   // assume the answer to all prompts is yes
var inputCheckRunID int64 // github check run id
var inputWeb3 bool        // deploy to web3

func init() {
	rootCmd.AddCommand(deployCmd)

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	deployCmd.Flags().StringVarP(&inputProjectName, "project", "p", "", "current project name")
	deployCmd.Flags().StringVarP(&inputProjectType, "type", "t", "", "current project type, e.g. react")
	deployCmd.Flags().StringVarP(&inputStaticDir, "static", "", "", "static dir name (if deploy type is static)")

	deployCmd.Flags().BoolVarP(&inputDetach, "detach", "", false, "return immediately after submitted deployment")
	deployCmd.Flags().BoolVarP(&inputAssumeYes, "assume-yes", "y", false,
		"assume the answer to all prompts is yes")

	deployCmd.Flags().BoolVarP(&inputProd, "prod", "", false, "deploy in production channel, will assign linked domain")
	deployCmd.Flags().BoolVarP(&inputDev, "dev", "", false, "deploy in development channel")

	deployCmd.Flags().BoolVarP(&inputWeb3, "web3", "", false, "deploy in web3 infra, store files on arweave, "+
		"visit via ipfs")

	deployCmd.Flags().BoolVarP(&inputCN, "cn", "", true, "deploy in mainland of china")
	deployCmd.Flags().MarkHidden("cn")

	deployCmd.Flags().Int64VarP(&inputCheckRunID, "check-run-id", "", 0, "github check run id")
	deployCmd.Flags().MarkHidden("check-run-id")

}

func SetupCloseHandler() {
	channel := make(chan os.Signal)
	signal.Notify(channel, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-channel
		if len(DeploymentID) > 0 {
			succeeded, err := requests.CancelDeployment(DeploymentID)
			if err != nil {
				log.S.StopFail()
				log.Error(err)
				os.Exit(0)
			}
			if succeeded {
				log.S.StopFail()
				log.Warning("Deployment canceled")
				os.Exit(0)
			} else {
				log.S.StopFail()
				log.Warning("Deployment cancellation failed")
				os.Exit(0)
			}
		} else {
			log.S.StopFail()
			log.Warning("Deployment canceled")
			os.Exit(0)
		}

	}()
}
