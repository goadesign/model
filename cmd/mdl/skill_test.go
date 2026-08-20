// These tests protect skill installation as a non-destructive repository
// operation: creation and upgrades are explicit, while local edits survive by
// default.
package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInstallSkill(t *testing.T) {
	t.Run("creates canonical skill", func(t *testing.T) {
		root := t.TempDir()

		path, changed, err := installSkill(root, false)

		require.NoError(t, err)
		assert.True(t, changed)
		assert.Equal(t, filepath.Join(root, filepath.FromSlash(diagramSkillPath)), path)
		content, err := os.ReadFile(path)
		require.NoError(t, err)
		assert.Equal(t, diagramSkill, content)
	})

	t.Run("is idempotent", func(t *testing.T) {
		root := t.TempDir()
		_, _, err := installSkill(root, false)
		require.NoError(t, err)

		path, changed, err := installSkill(root, false)

		require.NoError(t, err)
		assert.False(t, changed)
		assert.Equal(t, filepath.Join(root, filepath.FromSlash(diagramSkillPath)), path)
	})

	t.Run("preserves local changes", func(t *testing.T) {
		root := t.TempDir()
		target := filepath.Join(root, filepath.FromSlash(diagramSkillPath))
		require.NoError(t, os.MkdirAll(filepath.Dir(target), 0o755))
		require.NoError(t, os.WriteFile(target, []byte("local"), 0o644))

		_, changed, err := installSkill(root, false)

		assert.ErrorContains(t, err, "contains local changes")
		assert.False(t, changed)
		content, readErr := os.ReadFile(target)
		require.NoError(t, readErr)
		assert.Equal(t, []byte("local"), content)
	})

	t.Run("force replaces local changes", func(t *testing.T) {
		root := t.TempDir()
		target := filepath.Join(root, filepath.FromSlash(diagramSkillPath))
		require.NoError(t, os.MkdirAll(filepath.Dir(target), 0o755))
		require.NoError(t, os.WriteFile(target, []byte("local"), 0o644))

		path, changed, err := installSkill(root, true)

		require.NoError(t, err)
		assert.True(t, changed)
		assert.Equal(t, target, path)
		content, err := os.ReadFile(target)
		require.NoError(t, err)
		assert.Equal(t, diagramSkill, content)
	})
}
