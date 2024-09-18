package artifact

import "time"

type Artifact struct {
	Version   string
	CreatedAt time.Time
}

type ArtifactMetadata struct {
	Version   string    `json:"version"`
	CreatedAt time.Time `json:"created_at"`
}
