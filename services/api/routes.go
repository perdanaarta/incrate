package api

import (
	"incrate/services/log"

	"github.com/gin-gonic/gin"
)

func NewAPIsRouter() *gin.Engine {

	if !(len(log.DefaultLogWriter.MultiWriter) < 1) {
		gin.DefaultWriter = log.DefaultLogWriter.Writer
	}

	router := gin.Default()

	api := router.Group("/api/v1")

	artifactH := NewArtifactAPI()
	atf := api.Group("/version")
	{
		atf.POST("/new", artifactH.Create)
		atf.GET(":version", artifactH.Get)

		atf.POST(":version/upload", artifactH.UploadItem)
		atf.GET(":version/:filename", artifactH.DownloadItem)
	}

	return router
}
