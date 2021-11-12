package version

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_versionCheck(t *testing.T) {
	is := assert.New(t)

	tests := []struct {
		semanticVersion string
		isVulnerable    bool
	}{
		{
			semanticVersion: "1.10.0",
			isVulnerable:    true,
		},
		{
			semanticVersion: "1.6.0",
			isVulnerable:    true,
		},
		{
			semanticVersion: "1.10.9",
			isVulnerable:    false,
		},
		{
			semanticVersion: "1.10.11",
			isVulnerable:    false,
		},
	}

	for _, test := range tests {
		// vuln, err := CheckVersion(VersionCheckUrl, "Geth/v1.10.0-")
		vuln, err := CheckVersion(VersionCheckUrl, SemanticVersionToGethVersion(test.semanticVersion))
		is.NoError(err)

		isVulnerable := vuln != nil

		is.Equal(test.isVulnerable, isVulnerable, "semantic version: %s", test.semanticVersion)
	}
}

func Test_SemanticVersionToGethVersion(t *testing.T) {
	is := assert.New(t)

	tests := []struct {
		semanticVersion string
		gethVersion     string
	}{
		{
			semanticVersion: "1.10.0",
			gethVersion:     "Geth/v1.10.0-",
		},
		{
			semanticVersion: "1.6.0",
			gethVersion:     "Geth/v1.6.0-",
		},
	}
	for _, test := range tests {
		gethVersion := SemanticVersionToGethVersion(test.semanticVersion)
		is.Equal(test.gethVersion, gethVersion)
	}
}
