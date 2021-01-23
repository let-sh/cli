package sentry

import (
	"github.com/getsentry/sentry-go"
	"github.com/let-sh/cli/cmd"
	"github.com/let-sh/cli/utils"
	"github.com/shirou/gopsutil/host"
)

func init() {
	sentry.Init(sentry.ClientOptions{
		Dsn:     "https://f201c9f3cd0e473e98ad25cde46053dc@o310861.ingest.sentry.io/5604834",
		Release: cmd.Version,
	})

	hostInfo, _ := host.Info()
	sentry.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetExtra("version", cmd.Version)
		scope.SetExtra("token", utils.Credentials.Token)
		scope.SetExtra("os", hostInfo.OS)
		scope.SetExtra("platform", hostInfo.Platform)
		scope.SetExtra("platform_family", hostInfo.PlatformFamily)
		scope.SetExtra("platform_version", hostInfo.PlatformVersion)
	})
}

func Init() {}
