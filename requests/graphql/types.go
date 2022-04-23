package graphql

type QueryAllPreference struct {
	AllPreference struct {
		Channel string `graphql:"channel"`
	}
}
type QueryPreference struct {
	Preference string `graphql:"preference(name: $name)"`
}
type QueryCheckDeployCapability struct {
	CheckDeployCapability struct {
		HashID string `graphql:"hashID"`
		Exists bool   `graphql:"exists"`
	} `graphql:"checkDeployCapability(projectName:$projectName,cn:$cn)"`
}
type QueryStsToken struct {
	StsToken struct {
		Host            string `graphql:"host"`
		AccessKeyID     string `graphql:"accessKeyID"`
		AccessKeySecret string `graphql:"accessKeySecret"`
		SecurityToken   string `graphql:"securityToken"`
	} `graphql:"stsToken(type:$tokenType,projectName:$projectName,cn:$cn)"`
}
type QueryDeployment struct {
	Deployment struct {
		TargetFQDN   string `graphql:"targetFQDN"`
		NetworkStage string `graphql:"networkStage"`
		PackerStage  string `graphql:"packerStage"`
		Status       string `graphql:"status"`
		Done         bool   `graphql:"done"`
		ErrorLogs    string `graphql:"errorLogs"`
	} `graphql:"deployment(id:$id)"`
}
type QueryBuildTemplate struct {
	BuildTemplate struct {
		ContainsStatic   bool     `graphql:"containsStatic"`
		ContainsDynamic  bool     `graphql:"containsDynamic"`
		RequireCompiling bool     `graphql:"requireCompiling"`
		LocalCompiling   bool     `graphql:"localCompiling"`
		CompileCommands  []string `graphql:"compileCommands"`
		DistDir          string   `graphql:"distDir"`
	} `graphql:"buildTemplate(type:$type)"`
}
type MutationSetPreference struct {
	SetPreference bool `graphql:"setPreference(name: $name, value: $value)"`
}
type MutationDeploy struct {
	Deploy struct {
		ID           string `graphql:"id" json:"id"`
		TargetFQDN   string `graphql:"targetFQDN" json:"targetFQDN"`
		NetworkStage string `graphql:"networkStage" json:"networkStage"`
		PackerStage  string `graphql:"packerStage" json:"packerStage"`
		Status       string `graphql:"status" json:"status"`
		Project      struct {
			ID string `graphql:"id" json:"id"`
		} `graphql:"project" json:"project"`
	} `graphql:"deploy(input:{type:$type, projectName:$name, config:$config, channel:$channel, cn:$cn})"`
}

type MutationDeployWithCheckRunID struct {
	Deploy struct {
		ID           string `graphql:"id" json:"id"`
		TargetFQDN   string `graphql:"targetFQDN" json:"targetFQDN"`
		NetworkStage string `graphql:"networkStage" json:"networkStage"`
		PackerStage  string `graphql:"packerStage" json:"packerStage"`
		Status       string `graphql:"status" json:"status"`
		Project      struct {
			ID string `graphql:"id" json:"id"`
		} `graphql:"project" json:"project"`
	} `graphql:"deploy(input:{type:$type, projectName:$name, config:$config, channel:$channel, cn:$cn, checkRunID:$checkRunID})"`
}

type MutationCancelDeployment struct {
	CancelDeployment bool `graphql:"cancel(deploymentID:$deploymentID)""`
}

type QueryDeployments struct {
	Deployments struct {
		Edges []struct {
			Node struct {
				ID string `graphql:"id"`
			} `graphql:"node"`
		} `graphql:"edges"`
	} `graphql:"deployments(first:$first,projectName:$projectName,orderBy:{direction:DESC,field:UPDATED_AT})"`
}

type MutationStartDevelopment struct {
	StartDevelopment struct {
		RemotePort    int    `graphql:"remotePort" json:"remotePort,omitempty"`
		RemoteAddress string `graphql:"remoteAddress" json:"remoteAddress,omitempty"`
		Fqdn          string `graphql:"fqdn" json:"fqdn,omitempty"`
	} `graphql:"startDevelopment(projectID:$projectID)"`
}

type QueryProject struct {
	Project struct {
		ID   string `graphql:"id"`
		Name string `graphql:"name"`
	} `graphql:"project(name:$projectName)"`
}

type MutationStopDevelopment struct {
	StopDevelopment bool `graphql:"stopDevelopment(projectID:$projectID)""`
}
type MutationLink struct {
	Link bool `graphql:"link(projectID:$projectID,hostname:$hostname)"`
}
type MutationUnlink struct {
	Unlink bool `graphql:"unlink(projectID:$projectID,hostname:$hostname)"`
}
