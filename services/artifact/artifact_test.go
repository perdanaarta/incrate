package artifact_test

import (
	"bytes"
	"incrate/services/artifact"
	"mime/multipart"
	"os"
	"testing"
)

type ArtifactTest struct {
	Storage       string
	Version       string
	VersionLatest string
	Filename      string
	Test          *testing.T

	Service *artifact.ArtifactService
}

func TestArtifact(t *testing.T) {
	storage := "../../../.test/artifact"

	s := artifact.NewArtifactService(storage)

	test := &ArtifactTest{
		Test:          t,
		Storage:       storage,
		Version:       "1.0.0",
		VersionLatest: "1.0.1",
		Filename:      "artifact.zip",
		Service:       s,
	}
	defer test.Cleanup()

	test.NewArtifact()
	test.GetArtifact()
	test.GetArtifactLatest()
	test.StoreArtifact()
	test.GetArtifactItem()
}

/*
Cleanup test unit
*/
func (t *ArtifactTest) Cleanup() {
	os.RemoveAll(t.Storage)
}

/*
Test making new artifact. Should not fail
*/
func (t *ArtifactTest) NewArtifact() {
	res, err := t.Service.New(t.Version)

	if err != nil {
		t.Test.Errorf("Expected %v, got %v", "nil", err.Error())
	}

	if res.Version != t.Version {
		t.Test.Errorf("Expected Version %v, but got %v", t.Version, res.Version)
		t.Test.Log(res)
	}
}

/*
Test getting artifact. Should not fail
*/
func (t *ArtifactTest) GetArtifact() {
	res, err := t.Service.Get(t.Version)

	if err != nil {
		t.Test.Errorf("Expected %v, got %v", "nil", err.Error())
	}

	if res.Version != t.Version {
		t.Test.Errorf("Expected Version %v, but got %v", t.Version, res.Version)
		t.Test.Log(res)
	}
}

/*
Test getting latest artifact version. Should not fail
*/
func (t *ArtifactTest) GetArtifactLatest() {

	// Making the latest version artifact
	func() {
		res, err := t.Service.New(t.VersionLatest)

		if err != nil {
			t.Test.Errorf("Expected %v, got %v", "nil", err.Error())
		}

		if res.Version != t.VersionLatest {
			t.Test.Errorf("Expected Version %v, but got %v", t.VersionLatest, res.Version)
			t.Test.Log(res)
		}
	}()

	// Getting latest version artifact
	func() {
		res, err := t.Service.GetLatest()

		if err != nil {
			t.Test.Errorf("Expected %v, got %v", "nil", err.Error())
		}

		if res.Version != t.VersionLatest {
			t.Test.Errorf("Expected Version %v, but got %v", t.VersionLatest, res.Version)
			t.Test.Log(res)
		}
	}()
}

/*
Test storing an artifact. Should not fail
*/
func (t *ArtifactTest) StoreArtifact() {
	artifact := func() *artifact.Artifact {
		artifact, err := t.Service.GetLatest()

		if err != nil {
			t.Test.Errorf("Expected %v, got %v", "nil", err.Error())
		}

		if artifact.Version != t.VersionLatest {
			t.Test.Errorf("Expected Version %v, but got %v", t.VersionLatest, artifact.Version)
			t.Test.Log(artifact)
		}

		return artifact
	}()

	file := func() *bytes.Reader {
		var body bytes.Buffer
		writer := multipart.NewWriter(&body)
		part, err := writer.CreateFormFile("file", "artifact.zip")
		if err != nil {
			t.Test.Errorf("Error creating form file: %v", err)
		}

		// Simulate file content
		fileContent := []byte("dummy content")
		_, err = part.Write(fileContent)
		if err != nil {
			t.Test.Errorf("Error writing file content: %v", err)
		}

		// Close the writer to finalize the multipart form
		writer.Close()

		// Use the multipart data
		file := bytes.NewReader(body.Bytes())
		return file
	}()

	if err := t.Service.Store(artifact, t.Filename, file); err != nil {
		t.Test.Errorf("Error storing form file: %v", err)
	}
}

/*
Test getting an artifact item. Should not fail
*/
func (t *ArtifactTest) GetArtifactItem() {
	artifact := func() *artifact.Artifact {
		artifact, err := t.Service.GetLatest()

		if err != nil {
			t.Test.Errorf("Expected %v, got %v", "nil", err.Error())
		}

		if artifact.Version != t.VersionLatest {
			t.Test.Errorf("Expected Version %v, but got %v", t.VersionLatest, artifact.Version)
			t.Test.Log(artifact)
		}

		return artifact
	}()

	item, exist := artifact.Items[t.Filename]
	if !exist {
		t.Test.Errorf("Expecting %v, got %v", t.Filename, item.Filename)
	}
}
