package artifact

import "time"

type Artifact struct {
	Version   string
	CreatedAt time.Time
	Items     map[string]ArtifactItem
}

type ArtifactMetadata struct {
	Version   string    `json:"version"`
	CreatedAt time.Time `json:"created_at"`
}

type ArtifactItem struct {
	Filename string
	Path     string
}
