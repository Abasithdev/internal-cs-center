package utils

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func QueryInt(ctx *gin.Context, key string, defaultVal int) int {
	value := ctx.DefaultQuery(key, "")
	if value == "" {
		return defaultVal
	}

	var i int
	_, err := fmt.Sscan(value, &i)
	if err != nil || i <= 0 {
		return defaultVal
	}
	return i
}
