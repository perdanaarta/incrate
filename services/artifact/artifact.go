package artifact

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"time"
)

type ArtifactService struct {
	StoragePath string
}

func NewArtifactService(storage_path string) *ArtifactService {
	return &ArtifactService{
		StoragePath: storage_path,
	}
}

func (s *ArtifactService) GetFromMetadata(metadata ArtifactMetadata) *Artifact {
	return &Artifact{
		Version:   metadata.Version,
		CreatedAt: metadata.CreatedAt,
	}
}

func (s *ArtifactService) NewMetadata(artifact_dir string, artifact *Artifact) (ArtifactMetadata, error) {
	metadata_file := filepath.Join(artifact_dir, "metadata.json")

	metadata := ArtifactMetadata{
		Version:   artifact.Version,
		CreatedAt: artifact.CreatedAt,
	}

	file, err := os.Create(metadata_file)
	if err != nil {
		return ArtifactMetadata{}, errors.New("missing metadata")
	}
	defer file.Close()

	if err := json.NewEncoder(file).Encode(&metadata); err != nil {
		return ArtifactMetadata{}, errors.New("invalid metadata file format")
	}

	return metadata, nil
}

func (s *ArtifactService) GetMetadata(artifact_dir string) (ArtifactMetadata, error) {
	metadata_file := filepath.Join(artifact_dir, "metadata.json")

	var metadata ArtifactMetadata

	file, err := os.Open(metadata_file)
	if err != nil {
		return ArtifactMetadata{}, errors.New("missing metadata")
	}
	defer file.Close()

	if err := json.NewDecoder(file).Decode(&metadata); err != nil {
		return ArtifactMetadata{}, errors.New("invalid metadata file format")
	}

	return metadata, nil
}

func (s *ArtifactService) New(version_number string) (*Artifact, error) {
	artifact_dir := filepath.Join(s.StoragePath, version_number)

	if _, err := os.Stat(artifact_dir); err == nil {
		return &Artifact{}, errors.New("artifact exist")
	}

	if err := os.MkdirAll(artifact_dir, os.ModePerm); err != nil {
		return &Artifact{}, err
	}

	artifact := &Artifact{
		Version:   version_number,
		Items:     make(map[string]ArtifactItem),
		CreatedAt: time.Now(),
	}

	s.NewMetadata(artifact_dir, artifact)

	return artifact, nil
}

func (s *ArtifactService) Get(version_number string) (*Artifact, error) {
	artifact_dir := filepath.Join(s.StoragePath, version_number)

	if _, err := os.Stat(artifact_dir); err != nil {
		return &Artifact{}, errors.New("artifact not found")
	}

	metadata, err := s.GetMetadata(artifact_dir)
	if err != nil {
		return &Artifact{}, err
	}

	artifact := s.GetFromMetadata(metadata)
	if err := s.LoadItems(artifact); err != nil {
		return artifact, err
	}

	return artifact, nil
}

func (s *ArtifactService) GetLatest() (*Artifact, error) {
	entries, err := os.ReadDir(s.StoragePath)
	if err != nil {
		return &Artifact{}, err
	}

	metadata := ArtifactMetadata{}
	for _, entry := range entries {
		artifact_dir := filepath.Join(s.StoragePath, entry.Name())
		metadata_file := filepath.Join(artifact_dir, "metadata.json")

		var tmp_metadata ArtifactMetadata

		// Get artifact metadata
		if entry.IsDir() {
			if _, err := os.Stat(metadata_file); err == nil {
				file, err := os.Open(metadata_file)
				if err != nil {
					continue
				}
				defer file.Close()

				if err := json.NewDecoder(file).Decode(&tmp_metadata); err != nil {
					continue
				}
			}
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

func (s *ArtifactService) Store(artifact *Artifact, filename string, file io.Reader) error {
	filepath := filepath.Join(s.StoragePath, artifact.Version, filename)

	dst, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		return err
	}

	return nil
}

func (s *ArtifactService) LoadItems(artifact *Artifact) error {
	artifact_dir := filepath.Join(s.StoragePath, artifact.Version)

	entries, err := os.ReadDir(artifact_dir)
	if err != nil {
		return err
	}

	if len(artifact.Items) == 0 {
		artifact.Items = make(map[string]ArtifactItem)
	}

	for _, entry := range entries {
		entry.Name()

		func(filename string) {
			if filename == "metadata.json" {
				return
			}

			item := ArtifactItem{
				Filename: filename,
				Path:     filepath.Join(artifact_dir, filename),
			}

			artifact.Items[filename] = item
		}(entry.Name())
	}

	return nil
}
