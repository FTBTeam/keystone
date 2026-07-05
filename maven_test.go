package keystone

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMissingVersionNotationParsing(t *testing.T) {
	notation := "com.example:my-library"
	_, err := ParseMavenNotation(notation)

	require.Error(t, err, "Expected error while parsing maven notation with missing version")
	assert.Contains(t, err.Error(), "does not contain a version", "Expected error message to contain 'does not contain a version'")
}

func TestMavenNotationParsing(t *testing.T) {
	notation := "com.example:my-library:1.0.0"
	parsed, err := ParseMavenNotation(notation)

	require.NoError(t, err, "Expected no error while parsing maven notation")
	assert.Equal(t, "com.example", parsed.GroupID, "Expected GroupID to be 'com.example'")
	assert.Equal(t, "my-library", parsed.ArtifactID, "Expected ArtifactID to be 'my-library'")
	assert.Equal(t, "1.0.0", parsed.Version, "Expected Version to be '1.0.0'")
	assert.Empty(t, parsed.Classifier, "Expected Type to be empty")
}

func TestBulkVariations(t *testing.T) {
	testCases := []struct {
		notation           string
		expectedGroup      string
		expectedArtifact   string
		expectedVersion    string
		expectedClassifier string
		expectedExtension  string
		expectError        bool
	}{
		// Check for non-dot domains prefix
		{"example:my-library:1.0.0", "example", "my-library", "1.0.0", "", "", false},
		// Check for snapshot versions
		{"com.example:my-library:1.0.0-SNAPSHOT", "com.example", "my-library", "1.0.0-SNAPSHOT", "", "", false},
		// Check for version set as range
		{"com.example:my-library:[1.0.0,2.0.0)", "com.example", "my-library", "[1.0.0,2.0.0)", "", "", false},
		// Check for classifier
		{"com.example:my-library:1.0.0:javadoc", "com.example", "my-library", "1.0.0", "javadoc", "", false},
		// Check for extension
		{"com.example:my-library:1.0.0@pom", "com.example", "my-library", "1.0.0", "", "pom", false},
		// Check for classifier + extension
		{"com.example:my-library:1.0.0:sources@jar", "com.example", "my-library", "1.0.0", "sources", "jar", false},
		// Check for error on missing data after split
		{"com.example:my-library:", "com.example", "my-library", "", "", "", true},
		// Check for error on empty extension
		{"com.example:my-library:1.0.0@", "com.example", "my-library", "1.0.0", "", "", true},
	}

	for _, tc := range testCases {
		t.Run(tc.notation, func(t *testing.T) {
			parsed, err := ParseMavenNotation(tc.notation)
			if tc.expectError {
				require.Error(t, err, "Expected error while parsing maven notation")
			} else {
				require.NoError(t, err, "Expected no error while parsing maven notation")
				assert.Equal(t, tc.expectedGroup, parsed.GroupID, "Expected GroupID to match")
				assert.Equal(t, tc.expectedArtifact, parsed.ArtifactID, "Expected ArtifactID to match")
				assert.Equal(t, tc.expectedVersion, parsed.Version, "Expected Version to match")
				assert.Equal(t, tc.expectedClassifier, parsed.Classifier, "Expected Classifier to match")
				assert.Equal(t, tc.expectedExtension, parsed.Extension, "Expected Extension to match")
			}
		})
	}
}

func TestPathCreation(t *testing.T) {
	testCases := []struct {
		notation     MavenNotation
		expectedPath string
	}{
		// Basic jar (default extension)
		{
			MavenNotation{GroupID: "net.fabricmc", ArtifactID: "fabric-loader", Version: "0.18.4"},
			"net/fabricmc/fabric-loader/0.18.4/fabric-loader-0.18.4.jar",
		},
		// Javadoc classifier
		{
			MavenNotation{GroupID: "net.fabricmc", ArtifactID: "fabric-loader", Version: "0.18.4", Classifier: "javadoc"},
			"net/fabricmc/fabric-loader/0.18.4/fabric-loader-0.18.4-javadoc.jar",
		},
		// Sources classifier
		{
			MavenNotation{GroupID: "net.fabricmc", ArtifactID: "fabric-loader", Version: "0.18.4", Classifier: "sources"},
			"net/fabricmc/fabric-loader/0.18.4/fabric-loader-0.18.4-sources.jar",
		},
		// JSON extension
		{
			MavenNotation{GroupID: "net.fabricmc", ArtifactID: "fabric-loader", Version: "0.18.4", Extension: "json"},
			"net/fabricmc/fabric-loader/0.18.4/fabric-loader-0.18.4.json",
		},
		// Launchwrapper classifier + JSON extension
		{
			MavenNotation{GroupID: "net.fabricmc", ArtifactID: "fabric-loader", Version: "0.18.4", Classifier: "launchwrapper", Extension: "json"},
			"net/fabricmc/fabric-loader/0.18.4/fabric-loader-0.18.4-launchwrapper.json",
		},
		// POM extension
		{
			MavenNotation{GroupID: "net.fabricmc", ArtifactID: "fabric-loader", Version: "0.18.4", Extension: "pom"},
			"net/fabricmc/fabric-loader/0.18.4/fabric-loader-0.18.4.pom",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.expectedPath, func(t *testing.T) {
			assert.Equal(t, tc.expectedPath, MavenNotationToPath(tc.notation), "Expected generated path to match")
		})
	}
}
