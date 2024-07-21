package version

import "testing"

func TestGet(t *testing.T) {
	v := Get()
	if v.BuildVersion != "N/A" {
		t.Errorf("Wrong value version %v", v.BuildVersion)
	}

	if v.GitCommit != "N/A" {
		t.Errorf("Wrong value GitCommit %v", v.GitCommit)
	}

	if v.BuildDate != "N/A" {
		t.Errorf("Wrong value BuildDate %v", v.BuildDate)
	}

}
