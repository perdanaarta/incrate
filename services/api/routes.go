package api

import "github.com/gin-gonic/gin"

func NewAPIsRouter() *gin.Engine {
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
