// These tests protect skill installation as a non-destructive repository
// operation: creation and upgrades are explicit, while local edits survive by
// default.
package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDetectAgentSkillLocations(t *testing.T) {
	tests := []struct {
		name        string
		projectDirs []string
		executables []string
		wantPaths   []string
		wantFound   bool
	}{
		{
			name:        "detects Cursor executable",
			executables: []string{"cursor-agent"},
			wantPaths:   []string{cursorDiagramSkillPath},
			wantFound:   true,
		},
		{
			name:        "detects Codex project directory",
			projectDirs: []string{".agents"},
			wantPaths:   []string{codexDiagramSkillPath},
			wantFound:   true,
		},
		{
			name:        "detects Claude Code executable",
			executables: []string{"claude"},
			wantPaths:   []string{claudeDiagramSkillPath},
			wantFound:   true,
		},
		{
			name:        "installs for every detected agent",
			projectDirs: []string{".cursor"},
			executables: []string{"codex", "claude"},
			wantPaths: []string{
				cursorDiagramSkillPath,
				codexDiagramSkillPath,
				claudeDiagramSkillPath,
			},
			wantFound: true,
		},
		{
			name:      "falls back to Cursor",
			wantPaths: []string{cursorDiagramSkillPath},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := t.TempDir()
			for _, dir := range tt.projectDirs {
				require.NoError(t, os.MkdirAll(filepath.Join(root, dir), 0o755))
			}

			locations, found, err := detectAgentSkillLocations(root, executableLookupFor(tt.executables))

			require.NoError(t, err)
			assert.Equal(t, tt.wantFound, found)
			assert.Equal(t, tt.wantPaths, skillLocationPaths(locations))
		})
	}
}

func TestInstallSkills(t *testing.T) {
	locations := []agentSkillLocation{
		agentSkillLocations[0],
		agentSkillLocations[2],
	}

	t.Run("creates every canonical skill", func(t *testing.T) {
		root := t.TempDir()

		results, err := installSkills(root, locations, false)

		require.NoError(t, err)
		require.Len(t, results, 2)
		for _, result := range results {
			assert.True(t, result.changed)
			content, readErr := os.ReadFile(result.path)
			require.NoError(t, readErr)
			assert.Equal(t, diagramSkill, content)
		}
	})

	t.Run("is idempotent", func(t *testing.T) {
		root := t.TempDir()
		_, err := installSkills(root, locations, false)
		require.NoError(t, err)

		results, err := installSkills(root, locations, false)

		require.NoError(t, err)
		require.Len(t, results, 2)
		assert.False(t, results[0].changed)
		assert.False(t, results[1].changed)
	})

	t.Run("preserves all targets when one contains local changes", func(t *testing.T) {
		root := t.TempDir()
		cursorTarget := filepath.Join(root, filepath.FromSlash(cursorDiagramSkillPath))
		claudeTarget := filepath.Join(root, filepath.FromSlash(claudeDiagramSkillPath))
		require.NoError(t, os.MkdirAll(filepath.Dir(claudeTarget), 0o755))
		require.NoError(t, os.WriteFile(claudeTarget, []byte("local"), 0o644))

		results, err := installSkills(root, locations, false)

		assert.ErrorContains(t, err, "contains local changes")
		assert.Empty(t, results)
		_, statErr := os.Stat(cursorTarget)
		assert.ErrorIs(t, statErr, os.ErrNotExist)
		content, readErr := os.ReadFile(claudeTarget)
		require.NoError(t, readErr)
		assert.Equal(t, []byte("local"), content)
	})

	t.Run("force replaces local changes", func(t *testing.T) {
		root := t.TempDir()
		target := filepath.Join(root, filepath.FromSlash(cursorDiagramSkillPath))
		require.NoError(t, os.MkdirAll(filepath.Dir(target), 0o755))
		require.NoError(t, os.WriteFile(target, []byte("local"), 0o644))

		results, err := installSkills(root, agentSkillLocations[:1], true)

		require.NoError(t, err)
		require.Len(t, results, 1)
		assert.True(t, results[0].changed)
		assert.Equal(t, target, results[0].path)
		content, err := os.ReadFile(target)
		require.NoError(t, err)
		assert.Equal(t, diagramSkill, content)
	})
}

// executableLookupFor builds a deterministic PATH lookup for discovery tests.
func executableLookupFor(executables []string) executableLookup {
	available := make(map[string]struct{}, len(executables))
	for _, executable := range executables {
		available[executable] = struct{}{}
	}
	return func(file string) (string, error) {
		if _, ok := available[file]; ok {
			return filepath.Join("/bin", file), nil
		}
		return "", exec.ErrNotFound
	}
}

// skillLocationPaths returns relative destinations for concise discovery assertions.
func skillLocationPaths(locations []agentSkillLocation) []string {
	paths := make([]string, 0, len(locations))
	for _, location := range locations {
		paths = append(paths, location.path)
	}
	return paths
}
