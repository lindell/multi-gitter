package bitbucketserver

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	bitbucketv1 "github.com/gfleury/go-bitbucket-v1"
	"github.com/stretchr/testify/assert"
)

func TestBitbucketServer_GetRepositories_SkipEmptyRepo(t *testing.T) {
	// Mock Bitbucket Server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Mock server received request: %s %s", r.Method, r.URL.Path)
		if strings.HasSuffix(r.URL.Path, "/rest/api/1.0/projects/PROJ/repos") { // GetRepositoriesWithOptions
			w.Header().Set("Content-Type", "application/json")
			// Simulate the structure used by mapstructure in the main code (bitbucketRepositoryPager)
			response := map[string]interface{}{
				"isLastPage": true,
				"values": []interface{}{
					map[string]interface{}{
						"id":      1,
						"slug":    "repo1",
						"name":    "repo1",
						"project": map[string]interface{}{"key": "PROJ"},
						"links": map[string]interface{}{
							"clone": []map[string]interface{}{
								{"name": "http", "href": fmt.Sprintf("http://%s/scm/proj/repo1.git", r.Host)},
								{"name": "ssh", "href": fmt.Sprintf("ssh://git@%s/proj/repo1.git", r.Host)},
							},
						},
					},
					map[string]interface{}{
						"id":      2,
						"slug":    "repo2",
						"name":    "repo2",
						"project": map[string]interface{}{"key": "PROJ"},
						"links": map[string]interface{}{
							"clone": []map[string]interface{}{
								{"name": "http", "href": fmt.Sprintf("http://%s/scm/proj/repo2.git", r.Host)},
								{"name": "ssh", "href": fmt.Sprintf("ssh://git@%s/proj/repo2.git", r.Host)},
							},
						},
					},
				},
			}
			// The GetRepositoriesWithOptions endpoint returns an object where the repositories are in Values
			// and pagination info is at the top level. This structure is what mapstructure expects.
			json.NewEncoder(w).Encode(response) // Encode the response map directly
		} else if strings.HasSuffix(r.URL.Path, "/rest/api/1.0/projects/PROJ/repos/repo1/branches/default") { // GetDefaultBranch for repo1
			w.Header().Set("Content-Type", "application/json")
			response := bitbucketv1.Branch{
				DisplayID: "main",
			}
			// The GetDefaultBranch endpoint returns an object where the branch is in Values
			// Forcing the structure based on mapstructure decoding in the original code
			json.NewEncoder(w).Encode(map[string]interface{}{"Values": response})
		} else if strings.HasSuffix(r.URL.Path, "/rest/api/1.0/projects/PROJ/repos/repo2/branches/default") { // GetDefaultBranch for repo2 (error)
			http.Error(w, "Simulated EOF or error fetching default branch", http.StatusInternalServerError)
		} else {
			log.Printf("Unhandled mock server request: %s %s", r.Method, r.URL.Path)
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	// Capture log output
	var logBuf bytes.Buffer
	log.SetOutput(&logBuf)
	defer log.SetOutput(os.Stderr) // Reset log output

	bbService, err := New("testuser", "testtoken", server.URL, false, false, func(rt http.RoundTripper) http.RoundTripper { return rt }, RepositoryListing{
		Projects: []string{"PROJ"},
	})
	assert.NoError(t, err)

	repos, err := bbService.GetRepositories(context.Background())
	assert.NoError(t, err)

	assert.Len(t, repos, 1, "Expected only one repository")
	if len(repos) == 1 {
		assert.Equal(t, "PROJ/repo1", repos[0].FullName(), "Expected PROJ/repo1 to be returned")
	}

	logOutput := logBuf.String()
	assert.Contains(t, logOutput, "Skipping repository PROJ/repo2 due to error fetching default branch", "Expected log message for skipped repository")

	// Verify that convertRepository was not called for repo2, and was for repo1.
	// This is implicitly tested by checking the number of returned repos and the log message.
	// A more direct way would involve more complex mocking or refactoring.
}
