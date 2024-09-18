package api

import (
	"fmt"
	"incrate/config"
	"incrate/services/artifact"
	"incrate/services/storage"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewArtifactAPI() *ArtifactAPIsHandler {
	c := config.New()
	return &ArtifactAPIsHandler{
		ArtifactService: artifact.NewArtifactService(storage.NewFileStorageProvider(c.Artifact.Storage)),
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
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
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

	// file, err := ctx.FormFile("file")
	// if err != nil {
	// 	ctx.JSON(http.StatusBadRequest, gin.H{"error": "file error"})
	// 	return
	// }

	// parts, err := file.Open()
	// if err != nil {
	// 	ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	return
	// }

	filename := ctx.GetHeader("X-Filename")
	if filename == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "filename header not provided"})
		return
	}

	if err := h.ArtifactService.Store(artifact, filename, ctx.Request.Body); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusAccepted, gin.H{"status": "file upload succeeded"})
}

func (h *ArtifactAPIsHandler) DownloadItem(ctx *gin.Context) {
	filename := ctx.Param("filename")

	artifact := h.getArtifact(ctx)
	if artifact == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "artifact not found"})
		return
	}

	file, err := h.ArtifactService.GetFile(artifact.Version, filename)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// item, exist := artifact.Items[filename]
	// if !exist {
	// 	ctx.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
	// 	return
	// }

	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	ctx.Header("Content-Type", "application/octet-stream")
	if _, err := io.Copy(ctx.Writer, file); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "fail to serve file from storage"})
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
}
