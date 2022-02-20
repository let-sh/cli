package requests

import (
	"context"
	"fmt"
	"github.com/let-sh/cli/info"
	"github.com/let-sh/cli/requests/http_client"
	"github.com/machinebox/graphql"
	"github.com/sirupsen/logrus"
)

//func init() {
//	//resolver, _ := dns.NewDoHResolver(
//	//	"https://dns.alidns.com/dns-query",
//	//	dns.DoHCache())
//	//
//	//dialer := &net.Dialer{
//	//	Resolver: resolver,
//	//	// &net.Resolver{
//	//	//	PreferGo: true,
//	//	//	Dial:&resolver,
//	//	//	func(ctx context.Context, network, address string) (net.Conn, error) {
//	//	//		d := net.Dialer{
//	//	//			Timeout: time.Duration(dnsResolverTimeoutMs) * time.Millisecond,
//	//	//		}
//	//	//		return d.DialContext(ctx, dnsResolverProto, dnsResolverIP)
//	//	//	}
//	//	//},
//	//}
//	//dialContext := func(ctx context.Context, network, addr string) (net.Conn, error) {
//	//	return dialer.DialContext(ctx, network, addr)
//	//}
//	//
//	//http.DefaultTransport.(*http.Transport).DialContext = dialContext
//	httpClient = &http.Client{}
//}

// deprecated: create a client (safe to share across requests)
func Graphql() *graphql.Client {
	return graphql.NewClient("https://api.let-sh.com/query", graphql.WithHTTPClient(http_client.NewClient()))
}

// deprecated: SetPreference
func SetPreference(name, value string) (ok bool, err error) {
	req := graphql.NewRequest(`
	mutation($name: String!, $value: String!) {
	  setPreference(name:$name, value:$value)
	}
`)

	req.Var("name", name)
	req.Var("value", value)
	req.Header.Set("Authorization", "Bearer "+info.Credentials.LoadToken())

	// run it and capture the response
	var respData struct {
		SetPreference bool `json:"setPreference,omitempty"`
	}
	if err := Graphql().Run(context.Background(), req, &respData); err != nil {
		return false, err
	}

	return respData.SetPreference, nil
}

// deprecated: GetAllPreference
func GetAllPreference() (data struct {
	Channel string `json:"channel"`
}, err error) {
	req := graphql.NewRequest(`
	query {
	  	allPreference {
			channel
		}	
	}
`)

	req.Header.Set("Authorization", "Bearer "+info.Credentials.LoadToken())

	// run it and capture the response
	var respData struct {
		AllPreference struct {
			Channel string `json:"channel"`
		} `json:"AllPreference,omitempty"`
	}
	if err := Graphql().Run(context.Background(), req, &respData); err != nil {
		return data, err
	}

	return respData.AllPreference, nil
}

// deprecated
func GetPreference(name string) (value string, err error) {
	req := graphql.NewRequest(`
	query($name: String!) {
	  preference(name:$name)
	}
`)

	req.Var("name", name)
	req.Header.Set("Authorization", "Bearer "+info.Credentials.LoadToken())

	// run it and capture the response
	var respData struct {
		Preference string `json:"preference,omitempty"`
	}
	if err := Graphql().Run(context.Background(), req, &respData); err != nil {
		return "", err
	}

	return respData.Preference, nil
}

// deprecated
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
	req.Header.Set("Authorization", "Bearer "+info.Credentials.LoadToken())

	// run it and capture the response
	var respData struct {
		CheckDeployCapability struct {
			HashID string `json:"hashID"`
			Exists bool   `json:"exists"`
		} `json:"checkDeployCapability,omitempty"`
	}
	if err := Graphql().Run(context.Background(), req, &respData); err != nil {
		//if len(respData.Error.Errors) > 0 {
		//	return "", false, errors.New(respData.Error.Errors[0].Message)
		//}

		return "", false, err
	}

	return respData.CheckDeployCapability.HashID, respData.CheckDeployCapability.Exists, nil
}

// deprecated
func GetStsToken(uploadType, projectName string, cn bool) (data struct {
	Host            string `json:"host"`
	AccessKeyID     string `json:"accessKeyID"`
	AccessKeySecret string `json:"accessKeySecret"`
	SecurityToken   string `json:"securityToken"`
}, err error) {
	req := graphql.NewRequest(`
query($type: String!, $name: String!, $cn: Boolean!) {
	stsToken(type:$type,projectName:$name,cn:$cn) {
		host
		accessKeyID
		accessKeySecret
		securityToken
	}
}
`)

	req.Var("type", uploadType)
	req.Var("name", projectName)
	req.Var("cn", cn)
	req.Header.Set("Authorization", "Bearer "+info.Credentials.LoadToken())

	// run it and capture the response
	var respData struct {
		StsToken struct {
			Host            string `json:"host"`
			AccessKeyID     string `json:"accessKeyID"`
			AccessKeySecret string `json:"accessKeySecret"`
			SecurityToken   string `json:"securityToken"`
		} `json:"stsToken,omitempty"`
	}
	if err := Graphql().Run(context.Background(), req, &respData); err != nil {
		//if len(respData.Error.Errors) > 0 {
		//	return "", false, errors.New(respData.Error.Errors[0].Message)
		//}

		return data, err
	}

	return respData.StsToken, nil
}

// deprecated
func DeployWithCheckRunID(projectType, projectName, config, channel string, cn bool,
	checkRunID int64) (deployment struct {
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
mutation($type: String!, $name: String!, $config: String, $channel: String!, $cn: Boolean, $checkRunID: Int64) {
	  deploy(input:{type:$type, projectName:$name, config:$config, channel:$channel, cn:$cn, checkRunID:$checkRunID}) {
		id
		targetFQDN
		networkStage
		packerStage
		status
		project {
			id
		}
	  }
}
`)

	req.Var("type", projectType)
	req.Var("name", projectName)
	req.Var("config", config)
	req.Var("channel", channel)
	req.Var("cn", cn)
	req.Var("checkRunID", checkRunID)
	req.Header.Set("Authorization", "Bearer "+info.Credentials.LoadToken())
	logrus.Debugln(projectType, projectName, config, cn)
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

	if err := Graphql().Run(context.Background(), req, &respData); err != nil {
		//if len(respData.Error.Errors) > 0 {
		//	return "", false, errors.New(respData.Error.Errors[0].Message)
		//}

		return deployment, err
	}
	logrus.Debugln(respData)
	return respData.Deploy, nil
}

// deprecated
func Deploy(projectType, projectName, config, channel string, cn bool) (deployment struct {
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
mutation($type: String!, $name: String!, $config: String, $channel: String!, $cn: Boolean) {
	  deploy(input:{type:$type, projectName:$name, config:$config, channel:$channel, cn:$cn}) {
		id
		targetFQDN
		networkStage
		packerStage
		status
		project {
			id
		}
	  }
}
`)

	req.Var("type", projectType)
	req.Var("name", projectName)
	req.Var("config", config)
	req.Var("channel", channel)
	req.Var("cn", cn)
	req.Header.Set("Authorization", "Bearer "+info.Credentials.LoadToken())
	logrus.Debugln(projectType, projectName, config, cn)
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

	if err := Graphql().Run(context.Background(), req, &respData); err != nil {
		//if len(respData.Error.Errors) > 0 {
		//	return "", false, errors.New(respData.Error.Errors[0].Message)
		//}

		return deployment, err
	}
	logrus.Debugln(respData)
	return respData.Deploy, nil
}

type Web3 struct {
	IpfsCID string `json:"ipfsCID"`
	ArTID   string `json:"arTID"`
}

// deprecated
func GetDeploymentStatus(id string) (deployment struct {
	TargetFQDN   string `json:"targetFQDN"`
	NetworkStage string `json:"networkStage"`
	PackerStage  string `json:"packerStage"`
	Status       string `json:"status"`
	Done         bool   `json:"done"`
	ErrorLogs    string `json:"errorLogs"`
	Web3         *Web3  `json:"web3,omitempty"`
}, err error) {
	req := graphql.NewRequest(`
query($id: UUID!) {
	  deployment(id:$id) {
		targetFQDN
		networkStage
		packerStage
		status
		done
		errorLogs
		web3 {
          ipfsCID
          arTID
		}
	  }
}
`)
	req.Var("id", id)
	req.Header.Set("Authorization", "Bearer "+info.Credentials.LoadToken())

	// run it and capture the response
	var respData struct {
		Deployment struct {
			TargetFQDN   string `json:"targetFQDN"`
			NetworkStage string `json:"networkStage"`
			PackerStage  string `json:"packerStage"`
			Status       string `json:"status"`
			Done         bool   `json:"done"`
			ErrorLogs    string `json:"errorLogs"`
			Web3         *Web3  `json:"web3,omitempty"`
		} `json:"deployment,omitempty"`
	}
	if err := Graphql().Run(context.Background(), req, &respData); err != nil {

		return deployment, err
	}

	return respData.Deployment, nil
}

// deprecated
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
	req.Header.Set("Authorization", "Bearer "+info.Credentials.LoadToken())

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
	if err := Graphql().Run(context.Background(), req, &respData); err != nil {
		return buildTemplate, err
	}

	return respData.BuildTemplate, nil
}

// deprecated
func CancelDeployment(deploymentID string) (
	cancelDeploymentResult bool, err error) {
	req := graphql.NewRequest(`
mutation($deploymentID: UUID!) {
	cancelDeployment(deploymentID:$deploymentID)
}
`)
	req.Var("deploymentID", deploymentID)
	req.Header.Set("Authorization", "Bearer "+info.Credentials.LoadToken())

	// run it and capture the response
	var respData struct {
		CancelDeployment bool `json:"buildTemplate,omitempty"`
	}
	if err := Graphql().Run(context.Background(), req, &respData); err != nil {

		return cancelDeploymentResult, err
	}

	return respData.CancelDeployment, nil
}

// deprecated
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
	req.Header.Set("Authorization", "Bearer "+info.Credentials.LoadToken())

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
	if err := Graphql().Run(context.Background(), req, &respData); err != nil {
		return respData.Deployments, err
	}

	return respData.Deployments, nil
}

// deprecated
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
	req.Header.Set("Authorization", "Bearer "+info.Credentials.LoadToken())

	// run it and capture the response
	var respData struct {
		StartDevelopment struct {
			RemotePort    int    `json:"remotePort,omitempty"`
			RemoteAddress string `json:"remoteAddress,omitempty"`
			Fqdn          string `json:"fqdn,omitempty"`
		} `json:"startDevelopment,omitempty"`
	}
	if err := Graphql().Run(context.Background(), req, &respData); err != nil {
		return startDevelopmentResult, err
	}

	return respData.StartDevelopment, nil
}

// deprecated
func QueryProject(projectName string) (
	projectInfo struct {
		ID   string `json:"id,omitempty"`
		Name string `json:"name,omitempty"`
	}, err error) {
	req := graphql.NewRequest(`
query($projectName: String!) {
  project(name:$projectName){
    id
    name
  }
}
`)
	req.Var("projectName", projectName)
	req.Header.Set("Authorization", "Bearer "+info.Credentials.LoadToken())

	// run it and capture the response
	var respData struct {
		ID   string `json:"id,omitempty"`
		Name string `json:"name,omitempty"`
	}
	if err := Graphql().Run(context.Background(), req, &respData); err != nil {
		return projectInfo, err
	}

	return respData, nil
}

// deprecated
func StopDevelopment(projectID string) (
	stopDevelopmentResult bool, err error) {
	req := graphql.NewRequest(`
mutation($projectID: UUID!) {
	stopDevelopment(projectID:$projectID) 
}
`)
	req.Var("projectID", projectID)
	req.Header.Set("Authorization", "Bearer "+info.Credentials.LoadToken())

	// run it and capture the response
	var respData struct {
		StopDevelopment bool `json:"stopDevelopment,omitempty"`
	}
	if err := Graphql().Run(context.Background(), req, &respData); err != nil {

		return stopDevelopmentResult, err
	}

	return respData.StopDevelopment, nil
}

// deprecated
func Link(projectID string, hostname string) (
	linkResult bool, err error) {
	req := graphql.NewRequest(`
mutation($projectID: UUID!,$hostname: String!) {
	link(projectID:$projectID,hostname:$hostname)
}
`)
	req.Var("projectID", projectID)
	req.Var("hostname", hostname)
	req.Header.Set("Authorization", "Bearer "+info.Credentials.LoadToken())

	// run it and capture the response
	var respData struct {
		Link bool `json:"link,omitempty"`
	}
	if err := Graphql().Run(context.Background(), req, &respData); err != nil {

		return linkResult, err
	}

	return respData.Link, nil
}

// deprecated
func Unlink(projectID string, hostname string) (
	unlinkResult bool, err error) {
	req := graphql.NewRequest(`
mutation($projectID: UUID!,$hostname: String!) {
	unlink(projectID:$projectID,hostname:$hostname)
}
`)
	req.Var("projectID", projectID)
	req.Var("hostname", hostname)
	req.Header.Set("Authorization", "Bearer "+info.Credentials.LoadToken())

	// run it and capture the response
	var respData struct {
		Unlink bool `json:"unlink,omitempty"`
	}
	if err := Graphql().Run(context.Background(), req, &respData); err != nil {

		return unlinkResult, err
	}

	return respData.Unlink, nil
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
//	if err := Graphql().Run(context.Background(), req, &respData); err != nil {
//		return respData.Deployments, err
//	}
//
//	return respData.Deployments, nil
//}
