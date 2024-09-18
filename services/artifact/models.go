package artifact

import "time"

type Artifact struct {
	Version     string
	Description string
	CreatedAt   time.Time
}

type ArtifactMetadata struct {
	Version     string    `json:"version"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}
