package common

const (
	Lock   = 0
	UnLock = 1
)

const (
	GRPCTimeout = 350
)

const (
	Offline    = "offline"
	Health     = "health"
	NotRelease = "notrelease"
	Exception  = "exception"
)

const (
	GoPluginBinaryName   = "main"
	GoPluginBinaryNameSO = "main.so"

	GoPluginPlatform = "mac"
)

const (
	GoPluginPerformanceSwitch = false
)

const (
	Linux = "linux"
	Mac   = "mac"
)

const (
	SysPluginName = "Plugin"
)

//var	(
//	PluginEnv = map[string][]string{
//		"linux": {"GOOS=linux", "GOARCH=amd64", "go", "buld", "-buildmode=plugin"},
//	}
//)
