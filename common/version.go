package common

import (
	"fmt"
	"runtime"
)

var (
	// AppName 应用名称
	AppName = "ControllerManager"
	// Version 版本号
	Version = "v0.1"
	// Branch 分支
	Branch string
	// BuildDate 构建日期
	BuildDate string
	// CCGitHash 构建hash
	CCGitHash string
	// GoVersion 构建版本
	GoVersion = runtime.Version()
)

func PrintVersion() string {
	return fmt.Sprintf("%s (goVersion: %s, buildDate=%s)", Version, GoVersion, BuildDate)
}
