package requests

import (
	"context"
	"fmt"
	"github.com/let-sh/cli/info"
	"github.com/machinebox/graphql"
)

// create a client (safe to share across requests)
var Graphql = graphql.NewClient("https://api.let.sh.cn/query")

func CheckDeployCapability(projectName string) (hashID string, exists bool, err error) {
	req := graphql.NewRequest(`
    query($name: String!) {
        checkDeployCapability (projectName:$name) {
			hashID
			exists
        }
    }
`)

	req.Var("name", projectName)
	req.Header.Set("Authorization", "Bearer "+info.Credentials.Token)

	// run it and capture the response
	var respData struct {
		CheckDeployCapability struct {
			HashID string `json:"hashID"`
			Exists bool   `json:"exists"`
		} `json:"checkDeployCapability,omitempty"`
	}
	if err := Graphql.Run(context.Background(), req, &respData); err != nil {
		//if len(respData.Error.Errors) > 0 {
		//	return "", false, errors.New(respData.Error.Errors[0].Message)
		//}

		return "", false, err
	}

	return respData.CheckDeployCapability.HashID, respData.CheckDeployCapability.Exists, nil
}

func GetStsToken(uploadType, projectName string) (data struct {
	Host            string `json:"host"`
	AccessKeyID     string `json:"accessKeyID"`
	AccessKeySecret string `json:"accessKeySecret"`
	SecurityToken   string `json:"securityToken"`
}, err error) {
	req := graphql.NewRequest(`
query($type: String!, $name: String!) {
	stsToken(type:$type,projectName:$name) {
		host
		accessKeyID
		accessKeySecret
		securityToken
	}
}
`)

	req.Var("type", uploadType)
	req.Var("name", projectName)
	req.Header.Set("Authorization", "Bearer "+info.Credentials.Token)

	// run it and capture the response
	var respData struct {
		StsToken struct {
			Host            string `json:"host"`
			AccessKeyID     string `json:"accessKeyID"`
			AccessKeySecret string `json:"accessKeySecret"`
			SecurityToken   string `json:"securityToken"`
		} `json:"stsToken,omitempty"`
	}
	if err := Graphql.Run(context.Background(), req, &respData); err != nil {
		//if len(respData.Error.Errors) > 0 {
		//	return "", false, errors.New(respData.Error.Errors[0].Message)
		//}

		return data, err
	}

	return respData.StsToken, nil
}

func Deploy(projectType, projectName, config string, cn bool) (deployment struct {
	ID           string `json:"id"`
	TargetFQDN   string `json:"targetFQDN"`
	NetworkStage string `json:"networkStage"`
	PackerStage  string `json:"packerStage"`
	Status       string `json:"status"`
	Project      struct {
		ID string `json:"id"`
	} `json:"project"`
}, err error) {
	req := graphql.NewRequest(`
mutation($type: String!, $name: String!, $config: String, $cn: Boolean) {
	  deploy(input:{type:$type, projectName:$name, config:$config, cn:$cn}) {
		id
		targetFQDN
		networkStage
		packerStage
		status
	  }
}
`)

	req.Var("type", projectType)
	req.Var("name", projectName)
	req.Var("config", config)
	req.Var("cn", cn)
	req.Header.Set("Authorization", "Bearer "+info.Credentials.Token)

	// run it and capture the response
	var respData struct {
		Deploy struct {
			ID           string `json:"id"`
			TargetFQDN   string `json:"targetFQDN"`
			NetworkStage string `json:"networkStage"`
			PackerStage  string `json:"packerStage"`
			Status       string `json:"status"`
			Project      struct {
				ID string `json:"id"`
			} `json:"project"`
		} `json:"deploy,omitempty"`
	}
	if err := Graphql.Run(context.Background(), req, &respData); err != nil {
		//if len(respData.Error.Errors) > 0 {
		//	return "", false, errors.New(respData.Error.Errors[0].Message)
		//}

		return deployment, err
	}

	return respData.Deploy, nil
}

func GetDeploymentStatus(id string) (deployment struct {
	TargetFQDN   string `json:"targetFQDN"`
	NetworkStage string `json:"networkStage"`
	PackerStage  string `json:"packerStage"`
	Status       string `json:"status"`
	Done         bool   `json:"done"`
	ErrorMessage string `json:"errorMessage"`
}, err error) {
	req := graphql.NewRequest(`
query($id: UUID!) {
	  deployment(id:$id) {
		targetFQDN
		networkStage
		packerStage
		status
		done
		errorMessage
	  }
}
`)
	req.Var("id", id)
	req.Header.Set("Authorization", "Bearer "+info.Credentials.Token)

	// run it and capture the response
	var respData struct {
		Deployment struct {
			TargetFQDN   string `json:"targetFQDN"`
			NetworkStage string `json:"networkStage"`
			PackerStage  string `json:"packerStage"`
			Status       string `json:"status"`
			Done         bool   `json:"done"`
			ErrorMessage string `json:"errorMessage"`
		} `json:"deployment,omitempty"`
	}
	if err := Graphql.Run(context.Background(), req, &respData); err != nil {

		return deployment, err
	}

	return respData.Deployment, nil
}

func GetTemplate(typeArg string) (
	buildTemplate struct {
		ContainsStatic   bool     `json:"containsStatic"`
		ContainsDynamic  bool     `json:"containsDynamic"`
		RequireCompiling bool     `json:"requireCompiling"`
		LocalCompiling   bool     `json:"localCompiling"`
		CompileCommands  []string `json:"compileCommands"`
		DistDir          string   `json:"distDir"`
	}, err error) {
	req := graphql.NewRequest(`
query($type: String!) {
  buildTemplate(type:$type) {
    containsStatic
    containsDynamic
    requireCompiling
    localCompiling
    compileCommands
    distDir
  }
}
`)
	req.Var("type", typeArg)
	req.Header.Set("Authorization", "Bearer "+info.Credentials.Token)

	// run it and capture the response
	var respData struct {
		BuildTemplate struct {
			ContainsStatic   bool     `json:"containsStatic"`
			ContainsDynamic  bool     `json:"containsDynamic"`
			RequireCompiling bool     `json:"requireCompiling"`
			LocalCompiling   bool     `json:"localCompiling"`
			CompileCommands  []string `json:"compileCommands"`
			DistDir          string   `json:"distDir"`
		} `json:"buildTemplate,omitempty"`
	}
	if err := Graphql.Run(context.Background(), req, &respData); err != nil {

		return buildTemplate, err
	}

	return respData.BuildTemplate, nil
}

func CancelDeployment(deploymentID string) (
	cancelDeploymentResult bool, err error) {
	req := graphql.NewRequest(`
mutation($deploymentID: UUID!) {
	cancelDeployment(deploymentID:$deploymentID)
}
`)
	req.Var("deploymentID", deploymentID)
	req.Header.Set("Authorization", "Bearer "+info.Credentials.Token)

	// run it and capture the response
	var respData struct {
		CancelDeployment bool `json:"buildTemplate,omitempty"`
	}
	if err := Graphql.Run(context.Background(), req, &respData); err != nil {

		return cancelDeploymentResult, err
	}

	return respData.CancelDeployment, nil
}

func QueryDeployments(projectName string, count int) (
	deployments struct {
		Edges []struct {
			Node struct {
				ID string `json:"id"`
			} `json:"node"`
		} `json:"edges"`
	}, err error) {
	req := graphql.NewRequest(fmt.Sprintf(`
query {
  deployments(first:%d,projectName:"%s",orderBy:{
    direction:DESC,
    field:UPDATED_AT
  }) {
    edges{
      node {
        id
      }
    }
  }
}
`, count, projectName))
	req.Header.Set("Authorization", "Bearer "+info.Credentials.Token)

	// run it and capture the response
	var respData struct {
		Deployments struct {
			Edges []struct {
				Node struct {
					ID string `json:"id"`
				} `json:"node"`
			} `json:"edges"`
		} `json:"deployments"`
	}
	if err := Graphql.Run(context.Background(), req, &respData); err != nil {
		return respData.Deployments, err
	}

	return respData.Deployments, nil
}

func StartDevelopment(projectID string) (
	startDevelopmentResult struct {
		RemotePort    int    `json:"remotePort,omitempty"`
		RemoteAddress string `json:"remoteAddress,omitempty"`
		Fqdn          string `json:"fqdn,omitempty"`
	}, err error) {
	req := graphql.NewRequest(`
mutation($projectID: UUID!) {
	startDevelopment(projectID:$projectID) {
		remotePort
    	remoteAddress
		fqdn
	}
}
`)
	req.Var("projectID", projectID)
	req.Header.Set("Authorization", "Bearer "+info.Credentials.Token)

	// run it and capture the response
	var respData struct {
		StartDevelopment struct {
			RemotePort    int    `json:"remotePort,omitempty"`
			RemoteAddress string `json:"remoteAddress,omitempty"`
			Fqdn          string `json:"fqdn,omitempty"`
		} `json:"startDevelopment,omitempty"`
	}
	if err := Graphql.Run(context.Background(), req, &respData); err != nil {
		return startDevelopmentResult, err
	}

	return respData.StartDevelopment, nil
}

func StopDevelopment(projectID string) (
	stopDevelopmentResult bool, err error) {
	req := graphql.NewRequest(`
mutation($projectID: UUID!) {
	stopDevelopment(projectID:$projectID) 
}
`)
	req.Var("projectID", projectID)
	req.Header.Set("Authorization", "Bearer "+info.Credentials.Token)

	// run it and capture the response
	var respData struct {
		StopDevelopment bool `json:"stopDevelopment,omitempty"`
	}
	if err := Graphql.Run(context.Background(), req, &respData); err != nil {

		return stopDevelopmentResult, err
	}

	return respData.StopDevelopment, nil
}

//func QueryLogs(deploymentID string, count int) (
//	deployments struct {
//		Edges []struct {
//			Node struct {
//				ID string `json:"id"`
//			} `json:"node"`
//		} `json:"edges"`
//	}, err error) {
//	req := graphql.NewRequest(fmt.Sprintf(`
//query {
//  deployments(first:%d,projectName:"%s",orderBy:{
//    direction:DESC,
//    field:UPDATED_AT
//  }) {
//    edges{
//      node {
//        id
//      }
//    }
//  }
//}
//`, count, projectName))
//	req.Header.Set("Authorization", "Bearer "+utils.Credentials.Token)
//
//	// run it and capture the response
//	var respData struct {
//		Deployments struct {
//			Edges []struct {
//				Node struct {
//					ID string `json:"id"`
//				} `json:"node"`
//			} `json:"edges"`
//		} `json:"deployments"`
//	}
//	if err := Graphql.Run(context.Background(), req, &respData); err != nil {
//		return respData.Deployments, err
//	}
//
//	return respData.Deployments, nil
//}
