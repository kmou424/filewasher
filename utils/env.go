package utils

import (
	"fmt"
	"github.com/kmou424/filewasher/types"
	"github.com/ohanakogo/exceptiongo"
	"syscall"
)

func RequireEnv(env string) string {
	value, exists := syscall.Getenv(env)
	if !exists {
		exceptiongo.QuickThrowMsg[types.EnvVarNotFoundException](fmt.Sprintf(`env "%s" not found`, env))
	}
	return value
}
