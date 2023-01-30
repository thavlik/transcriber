package base

import (
	"os"
	"strconv"
	"strings"
)

func CheckEnv(name string, target *string) string {
	if v, ok := os.LookupEnv(name); ok && target != nil {
		*target = v
		return v
	}
	return ""
}

func CheckEnvInt(name string, target *int) int {
	if v, ok := os.LookupEnv(name); ok && target != nil {
		e, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			panic(err)
		}
		c := int(e)
		*target = c
		return c
	}
	return 0
}

func CheckEnvInt64(name string, target *int64) int64 {
	if v, ok := os.LookupEnv(name); ok && target != nil {
		e, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			panic(err)
		}
		*target = e
		return e
	}
	return 0
}

func CheckEnvBool(name string, target *bool) bool {
	v := strings.ToLower(os.Getenv(name))
	b := !(v == "" || v == "0" || v == "false" || v == "no")
	if target != nil {
		*target = b
	}
	return b
}
