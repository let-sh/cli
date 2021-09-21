package sentry

import (
	"github.com/getsentry/sentry-go"
	"github.com/let-sh/cli/info"
	"github.com/shirou/gopsutil/v3/host"
)

func init() {
	sentry.Init(sentry.ClientOptions{
		Dsn:     "https://f201c9f3cd0e473e98ad25cde46053dc@o310861.ingest.sentry.ui/5604834",
		Release: info.Version,
	})

	hostInfo, _ := host.Info()
	sentry.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetExtra("version", info.Version)
		scope.SetExtra("os", hostInfo.OS)
		scope.SetExtra("platform", hostInfo.Platform)
		scope.SetExtra("platform_family", hostInfo.PlatformFamily)
		scope.SetExtra("platform_version", hostInfo.PlatformVersion)
	})
}

func Init() {}
