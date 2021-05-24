package ginhelper

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
	"strings"
)

func FileModified(c *gin.Context, file string) (bool, string) {
	eTag := GenerateFileETag(file)

	if match := c.Request.Header.Get("If-None-Match"); match != "" {
		if strings.Contains(match, eTag) {
			return false, eTag
		}
		return true, eTag
	}

	return true, eTag
}

func GenerateFileETag(file string) string {
	stat, err := os.Stat(file);
	if err != nil {
		return ""
	}
	if stat == nil {
		return ""
	}
	return fmt.Sprintf("%d", stat.ModTime().Unix())
}
