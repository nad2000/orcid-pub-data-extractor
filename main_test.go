package main

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
	dir, err := ioutil.TempDir("", "orcid")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	os.Args = []string{
		"-t=empl",
		"-o",
		dir,
		"./data/ORCID-API-2.0_activities_xml.tar.gz",
	}

	main()

	assert.Equal(t, 1, 1)
	assert.FileExists(t, path.Join(
		dir, "activities", "636", "0000-0001-6480-3636",
		"employments", "0000-0001-6480-3636_employments_772888.xml"))

	os.Args = []string{
		`-r='(?i)ma..yn\s*dav.*'`,
		"-o",
		dir,
		"./data/ORCID-API-2.0_activities_xml.tar.gz",
	}

	main()

	assert.Equal(t, 1, 1)
	assert.FileExists(t, path.Join(
		dir, "activities", "191", "0000-0003-4991-5191",
		"fundings", "0000-0003-4991-5191_fundings_212385.xml"))

}
