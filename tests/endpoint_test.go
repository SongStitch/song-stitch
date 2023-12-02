package tests

import (
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"
)

var testURLs []string

func init() {
	customDomain := os.Getenv("TEST_DOMAIN")
	if customDomain == "" {
		customDomain = "songstitch.art"
	}

	testURLs = []string{
		fmt.Sprintf(
			"https://%s/collage?username=theden_sh&method=album&period=overall&artist=true&album=true&playcount=true&rows=15&columns=15",
			customDomain,
		),
		fmt.Sprintf(
			"https://%s/collage?username=theden_sh&method=artist&period=overall&artist=true&playcount=true&rows=10&columns=10",
			customDomain,
		),
		fmt.Sprintf(
			"https://%s/collage?username=theden_sh&method=track&period=overall&track=true&artist=true&album=true&playcount=true&rows=5&columns=5",
			customDomain,
		),
	}
}

func testEndpoint(t *testing.T, url string) {
	start := time.Now()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatalf("error creating request: %v", err)
	}

	req.Header.Set("Cache-Control", "no-cache")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	defer resp.Body.Close()

	elapsed := time.Since(start)

	if resp.StatusCode == http.StatusOK {
		t.Logf("success: received %d for %s (time: %s)", resp.StatusCode, url, elapsed)
	} else {
		t.Errorf("error: received %d for %s (time: %s)", resp.StatusCode, url, elapsed)
	}
}

func TestEndpoint(t *testing.T) {
	for _, url := range testURLs {
		testEndpoint(t, url)
	}
}
