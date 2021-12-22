package mist

import "context"

func GetWm(name string) interface{} {
	return defaultManager.ms[defaultManager.mm[name]]
}

// 全局context,使用注意时注意暴露安全 TODO: hide from global
func FullContext() context.Context {
	return defaultManager.rootCtx
}
