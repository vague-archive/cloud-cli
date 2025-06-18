package share

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/vaguevoid/cloud-cli/internal/api"
	"github.com/vaguevoid/cloud-cli/internal/lib/crypto"
	"github.com/vaguevoid/cloud-cli/internal/lib/httpx"
)

//=================================================================================================
// DEPLOY COMMAND
//=================================================================================================

const (
	UploadConcurrency = 8
)

type DeployCommand struct {
	API       *api.Client
	Org       string
	Game      string
	Label     string
	Path      string
	OnStarted func(deployID int64, manifest []DeployEntry, incremental []DeployEntry)
	OnUpload  func(deployID int64, path string)
}

type DeployResult struct {
	DeployID int64
	Slug     string
	URL      string
	Manifest []DeployEntry
}

func Deploy(cmd *DeployCommand) (*DeployResult, error) {
	if cmd.API == nil {
		return nil, fmt.Errorf("missing api client")
	} else if cmd.Org == "" {
		return nil, fmt.Errorf("missing organization")
	} else if cmd.Game == "" {
		return nil, fmt.Errorf("missing game")
	} else if cmd.Path == "" {
		return nil, fmt.Errorf("missing path")
	}

	return cmd.execute()
}

type DeployEntry struct {
	Path          string `json:"path"`
	Blake3        string `json:"blake3"`
	ContentLength int    `json:"contentLength"`
}

//=================================================================================================
// PRIVATE IMPLEMENTATION
//=================================================================================================

func (cmd *DeployCommand) execute() (*DeployResult, error) {

	info, err := os.Stat(cmd.Path)
	if err != nil {
		return nil, fmt.Errorf("directory not found %s", cmd.Path)
	} else if !info.IsDir() {
		return nil, fmt.Errorf("%s is not a directory", cmd.Path)
	}

	fullManifest, err := cmd.buildManifest()
	if err != nil {
		return nil, err
	}

	deployID, incrementalManifest, err := cmd.startDeploy(fullManifest)
	if err != nil {
		return nil, err
	}

	if cmd.OnStarted != nil {
		cmd.OnStarted(deployID, fullManifest, incrementalManifest)
	}

	err = cmd.incrementalUpload(deployID, incrementalManifest)
	if err != nil {
		return nil, err
	}

	result, err := cmd.activateDeploy(deployID)
	if err != nil {
		return nil, err
	}

	result.Manifest = fullManifest
	return result, nil
}

//-------------------------------------------------------------------------------------------------

func (cmd *DeployCommand) buildManifest() ([]DeployEntry, error) {
	manifest := make([]DeployEntry, 0)

	disallowed := func(path string) bool {
		name := filepath.Base(path)
		return strings.HasSuffix(name, ".ssh") ||
			strings.HasSuffix(name, ".git") ||
			strings.HasSuffix(name, ".env")
	}

	err := filepath.Walk(cmd.Path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || disallowed(path) {
			return nil
		}
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		contentLength := info.Size()

		relPath, err := filepath.Rel(cmd.Path, path)
		if err != nil {
			return err
		}

		manifest = append(manifest, DeployEntry{
			Path:          relPath,
			Blake3:        crypto.Blake3(f),
			ContentLength: int(contentLength),
		})

		return nil
	})

	if err != nil {
		return nil, err
	}

	return manifest, nil
}

//-------------------------------------------------------------------------------------------------

func (cmd *DeployCommand) startDeploy(fullManifest []DeployEntry) (int64, []DeployEntry, error) {
	route := cmd.API.Route(cmd.Org, cmd.Game, "deploy", cmd.Label)
	resp, err := cmd.API.PostJSON(route, fullManifest)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusAccepted {
		deployID, err := strconv.ParseInt(resp.Header.Get(httpx.HeaderXDeployID), 10, 64)
		if err != nil {
			return 0, nil, fmt.Errorf("missing or invalid deployID: %s", err)
		}
		var incrementalManifest []DeployEntry
		err = json.NewDecoder(resp.Body).Decode(&incrementalManifest)
		return deployID, incrementalManifest, err
	} else {
		body, _ := io.ReadAll(resp.Body)
		return 0, nil, fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(body))
	}
}

//-------------------------------------------------------------------------------------------------

func (cmd *DeployCommand) incrementalUpload(deployID int64, incrementalManifest []DeployEntry) error {
	semaphore := make(chan struct{}, UploadConcurrency)
	errorChannel := make(chan error, len(incrementalManifest))
	var wg sync.WaitGroup

	for _, entry := range incrementalManifest {
		path := entry.Path // capture loop variable
		semaphore <- struct{}{}
		wg.Add(1)
		if cmd.OnUpload != nil {
			cmd.OnUpload(deployID, path)
		}
		go func() {
			defer wg.Done()
			defer func() { <-semaphore }()
			route := cmd.API.Route(cmd.Org, cmd.Game, "deploy", deployID, "upload", path)
			fullPath := filepath.Join(cmd.Path, path)
			resp, err := cmd.API.PostFILE(route, fullPath)
			if err != nil {
				errorChannel <- err
			} else if resp.StatusCode != http.StatusOK {
				errorChannel <- fmt.Errorf("failed to upload to %s: status code %d", route, resp.StatusCode)
			}
		}()
	}

	wg.Wait()
	close(errorChannel)

	var errors []error
	for err := range errorChannel {
		errors = append(errors, err)
	}
	if len(errors) > 0 {
		return fmt.Errorf("%v", errors)
	}

	return nil
}

//-------------------------------------------------------------------------------------------------

func (cmd *DeployCommand) activateDeploy(deployID int64) (*DeployResult, error) {
	route := cmd.API.Route(cmd.Org, cmd.Game, "deploy", deployID, "activate")
	resp, err := cmd.API.Post(route, nil)
	if err != nil {
		return nil, err
	} else if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to activate %s: status code %d", route, resp.StatusCode)
	}

	var result DeployResult
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

//-------------------------------------------------------------------------------------------------
