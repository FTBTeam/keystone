package keystone

import (
	"errors"
	"strconv"
	"strings"
)

type MavenNotation struct {
	GroupID    string
	ArtifactID string
	Version    string
	Classifier string // e.g. "javadoc", "sources"
	Extension  string // e.g. "jar", "json", "pom" (defaults to "jar")
}

// ParseMavenNotation attempts to parse a maven notation into a MavenNotation structure
// Examples:
// - com.example:my-library:1.0.0
// - com.example:my-library:1.0.0:javadoc
// - com.example:my-library:1.0.0:sources@jar
// - com.example:my-library:1.0.0@json
func ParseMavenNotation(notation string) (MavenNotation, error) {
	// Split off extension first (@)
	var extension string
	if atIdx := strings.LastIndex(notation, "@"); atIdx != -1 {
		extension = notation[atIdx+1:]
		notation = notation[:atIdx]
		if extension == "" {
			return MavenNotation{}, errors.New("MavenNotation contains empty extension after @")
		}
	}

	parts := strings.Split(notation, ":")
	if len(parts) < 3 {
		return MavenNotation{}, errors.New("MavenNotation does not contain a version")
	}

	for partIndex, part := range parts {
		if part == "" {
			return MavenNotation{}, errors.New("MavenNotation contains empty parts at index " + strconv.Itoa(partIndex))
		}
	}

	if len(parts) > 4 {
		return MavenNotation{}, errors.New("MavenNotation contains too many parts")
	}

	parsed := MavenNotation{
		GroupID:    parts[0],
		ArtifactID: parts[1],
		Version:    parts[2],
		Extension:  extension,
	}

	if len(parts) == 4 {
		parsed.Classifier = parts[3]
	}

	return parsed, nil
}

func MavenNotationToPath(notation MavenNotation) string {
	groupPath := strings.ReplaceAll(notation.GroupID, ".", "/")
	ext := notation.Extension
	if ext == "" {
		ext = "jar"
	}

	path := groupPath + "/" + notation.ArtifactID + "/" + notation.Version + "/" + notation.ArtifactID + "-" + notation.Version
	if notation.Classifier != "" {
		path += "-" + notation.Classifier
	}
	return path + "." + ext
}
