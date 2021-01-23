package log

import (
	"github.com/getsentry/sentry-go"
	"github.com/let-sh/cli/cmd"
)

func init() {
	sentry.Init(sentry.ClientOptions{
		Dsn: "https://f201c9f3cd0e473e98ad25cde46053dc@o310861.ingest.sentry.io/5604834",
	})
	sentry.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetExtra("version", cmd.Version)
	})
}
