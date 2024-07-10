package home2git

import (
	"testing"
)

func TestHome2Git(t *testing.T) {
	p, err := HomepageToGit("https://abseil.io/", "abseil")
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	if p.GitURL != "https://github.com/abseil/abseil-cpp" {
		t.Errorf("Error: %v", p.GitURL)
		t.FailNow()
	}
}
