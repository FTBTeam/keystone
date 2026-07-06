package keystone

import "os"

func GetEnvDefault(key string, defaultValue string) string {
	if env, ok := os.LookupEnv(key); ok {
		return env
	}
	return defaultValue
}
