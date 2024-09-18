package artifact

import (
	"bytes"
	"encoding/json"
	"errors"
	"incrate/services/storage"
	"io"
	"path/filepath"
	"time"
)

type ArtifactService struct {
	TempStorage     string
	StorageProvider storage.StorageProvider
}

func NewArtifactService(storageProvider storage.StorageProvider) *ArtifactService {
	return &ArtifactService{
		StorageProvider: storageProvider,
	}
}

func (s *ArtifactService) GetFromMetadata(metadata ArtifactMetadata) *Artifact {
	return &Artifact{
		Version:   metadata.Version,
		CreatedAt: metadata.CreatedAt,
	}
}

func (s *ArtifactService) NewMetadata(artifact *Artifact) (ArtifactMetadata, error) {
	metadata := ArtifactMetadata{
		Version:   artifact.Version,
		CreatedAt: artifact.CreatedAt,
	}

	var buf bytes.Buffer

	if err := json.NewEncoder(&buf).Encode(&metadata); err != nil {
		return metadata, errors.New("invalid metadata file format")
	}

	metadataPath := filepath.Join(artifact.Version, "metadata.json")
	if err := s.StorageProvider.Store(metadataPath, &buf); err != nil {
		return metadata, errors.New("fail to store metadata to storage")
	}

	return metadata, nil
}

func (s *ArtifactService) GetMetadata(artifactVersion string) (metadata ArtifactMetadata, err error) {
	metadataPath := filepath.Join(artifactVersion, "metadata.json")
	buf, err := s.StorageProvider.Get(metadataPath)
	if err != nil {
		return ArtifactMetadata{}, errors.New("missing metadata")
	}

	if err := json.NewDecoder(buf).Decode(&metadata); err != nil {
		return ArtifactMetadata{}, errors.New("invalid metadata file format")
	}

	return
}

func (s *ArtifactService) New(versionNumber string, description string) (*Artifact, error) {
	artifact := &Artifact{
		Version:     versionNumber,
		Description: description,
		CreatedAt:   time.Now(),
	}

	if _, err := s.StorageProvider.Get(versionNumber); err == nil {
		return &Artifact{}, errors.New("artifact exist")
	}

	if _, err := s.NewMetadata(artifact); err != nil {
		return &Artifact{}, err
	}

	return artifact, nil
}

func (s *ArtifactService) Get(versionNumber string) (*Artifact, error) {
	metadata, err := s.GetMetadata(versionNumber)
	if err != nil {
		return &Artifact{}, errors.New("artifact not found")
	}

	return s.GetFromMetadata(metadata), nil
}

func (s *ArtifactService) GetFile(versionNumber string, filename string) (io.Reader, error) {
	return s.StorageProvider.Get(filepath.Join(versionNumber, filename))
}

func (s *ArtifactService) GetLatest() (*Artifact, error) {
	entries, err := s.StorageProvider.List("")
	if err != nil {
		return &Artifact{}, err
	}

	metadata := ArtifactMetadata{}
	for _, entry := range entries {
		tmp_metadata, err := s.GetMetadata(entry)
		if err != nil {
			continue
		}

		// Juggle metadata until the latest version selected
		if metadata.Version == "" {
			metadata = tmp_metadata
		}

		if metadata.CreatedAt.Before(tmp_metadata.CreatedAt) {
			metadata = tmp_metadata
		}
	}

	return s.Get(metadata.Version)
}

func (s *ArtifactService) Store(artifact *Artifact, filename string, content io.Reader) (err error) {
	// Will add store validation in the future

	return s.StorageProvider.Store(
		filepath.Join(artifact.Version, filename),
		content,
	)
}
