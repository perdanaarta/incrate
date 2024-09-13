package api

import (
	"incrate/services/artifact"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewArtifactAPI() *ArtifactAPIsHandler {
	return &ArtifactAPIsHandler{
		ArtifactService: artifact.NewArtifactService("storage"),
	}
}

type ArtifactAPIsHandler struct {
	ArtifactService *artifact.ArtifactService
}

func (h *ArtifactAPIsHandler) getArtifact(ctx *gin.Context) *artifact.Artifact {
	version_number := ctx.Param("version")

	if version_number == "latest" {
		artifact, err := h.ArtifactService.GetLatest()
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "no versions were found"})
			return nil
		}

		ctx.JSON(http.StatusOK, artifact)
		return nil
	}

	artifact, err := h.ArtifactService.Get(version_number)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "version not found"})
		return nil
	}

	return artifact
}

func (h *ArtifactAPIsHandler) Get(ctx *gin.Context) {
	version_number := ctx.Param("version")

	if version_number == "latest" {
		artifact, err := h.ArtifactService.GetLatest()
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "no versions were found"})
			return
		}

		ctx.JSON(http.StatusOK, artifact)
		return
	}

	artifact, err := h.ArtifactService.Get(version_number)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "version not found"})
		return
	}

	ctx.JSON(http.StatusOK, artifact)
}

func (h *ArtifactAPIsHandler) Create(ctx *gin.Context) {
	var request struct {
		Version string `json:"version"`
	}

	ctx.BindJSON(&request)

	artifact, err := h.ArtifactService.New(request.Version)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, artifact)
}

func (h *ArtifactAPIsHandler) UploadItem(ctx *gin.Context) {
	artifact := h.getArtifact(ctx)
	if artifact == nil {
		return
	}

	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "file error"})
		return
	}

	parts, err := file.Open()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.ArtifactService.Store(artifact, file.Filename, parts); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Writer.WriteHeader(204)
}

func (h *ArtifactAPIsHandler) DownloadItem(ctx *gin.Context) {
	filename := ctx.Param("filename")

	artifact := h.getArtifact(ctx)
	if artifact == nil {
		return
	}

	item, exist := artifact.Items[filename]
	if !exist {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
		return
	}

	ctx.FileAttachment(item.Path, item.Filename)
}
