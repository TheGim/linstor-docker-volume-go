package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/docker/go-plugins-helpers/sdk"
	"github.com/docker/go-plugins-helpers/volume"
)

func TestVolumeServerClosesDockerConnectionAfterResponse(t *testing.T) {
	plugin := httptest.NewUnstartedServer(nil)
	plugin.Config = newVolumeServer(NewLinstorDriver("", "test-node", t.TempDir()))
	plugin.Start()
	defer plugin.Close()

	resp, err := plugin.Client().Post(plugin.URL+"/VolumeDriver.Capabilities", sdk.DefaultContentTypeV1_1, strings.NewReader("{}"))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status %s", resp.Status)
	}
	if !resp.Close {
		t.Fatal("Docker connection kept alive after the response")
	}
	var capabilities volume.CapabilitiesResponse
	if err := json.NewDecoder(resp.Body).Decode(&capabilities); err != nil {
		t.Fatal(err)
	}
	if capabilities.Capabilities.Scope != "global" {
		t.Fatalf("unexpected volume scope %q", capabilities.Capabilities.Scope)
	}
}
