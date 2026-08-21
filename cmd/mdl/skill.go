// The skill installer publishes MDL's canonical diagram-editing instructions
// to the project locations discovered by supported coding agents. It preserves
// locally edited skill files unless the caller explicitly authorizes replacement.
package main

import (
	"bytes"
	_ "embed"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type (
	agentSkillLocation struct {
		name        string
		projectDirs []string
		path        string
		executables []string
	}

	skillInstallResult struct {
		agent   string
		path    string
		changed bool
	}

	skillTargetState struct {
		location agentSkillLocation
		path     string
		changed  bool
	}

	executableLookup func(string) (string, error)
)

const (
	legacyCursorDiagramSkillPath = ".cursor/skills/editing-model-diagrams/SKILL.md"
	agentsDiagramSkillPath       = ".agents/skills/editing-model-diagrams/SKILL.md"
	claudeDiagramSkillPath       = ".claude/skills/editing-model-diagrams/SKILL.md"
)

var agentSkillLocations = []agentSkillLocation{
	{
		name:        "Cursor and Codex",
		projectDirs: []string{".cursor", ".agents"},
		path:        agentsDiagramSkillPath,
		executables: []string{"cursor", "cursor-agent", "codex"},
	},
	{
		name:        "Claude Code",
		projectDirs: []string{".claude"},
		path:        claudeDiagramSkillPath,
		executables: []string{"claude"},
	},
}

//go:embed skills/editing-model-diagrams/SKILL.md
var diagramSkill []byte

// runSkillInstall discovers the coding agents available to the current project
// and installs the canonical skill in every location those agents read.
func runSkillInstall(force bool) error {
	root, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("determine current repository: %w", err)
	}

	locations, detected, err := detectAgentSkillLocations(root, exec.LookPath)
	if err != nil {
		return err
	}
	results, removedLegacy, err := installSkills(root, locations, force)
	if err != nil {
		return err
	}
	if !detected {
		fmt.Println("No supported coding agent detected; using the portable Agent Skills fallback.")
	}
	for _, result := range results {
		if result.changed {
			fmt.Printf("Installed MDL diagram skill for %s at %s\n", result.agent, result.path)
		} else {
			fmt.Printf("MDL diagram skill for %s is already current at %s\n", result.agent, result.path)
		}
	}
	if removedLegacy != "" {
		fmt.Printf("Removed duplicate legacy Cursor skill at %s\n", removedLegacy)
	}
	return nil
}

// detectAgentSkillLocations returns every supported project location whose
// agent executable or project configuration directory is present. It retains
// the shared Agent Skills location as a fallback when no agent can be detected.
func detectAgentSkillLocations(root string, lookPath executableLookup) ([]agentSkillLocation, bool, error) {
	var locations []agentSkillLocation
	for _, location := range agentSkillLocations {
		detected, err := agentSkillLocationDetected(root, location, lookPath)
		if err != nil {
			return nil, false, err
		}
		if detected {
			locations = append(locations, location)
		}
	}
	if len(locations) > 0 {
		return locations, true, nil
	}
	return agentSkillLocations[:1], false, nil
}

// agentSkillLocationDetected reports whether the repository is configured for
// an agent or one of the agent's command-line executables is available.
func agentSkillLocationDetected(root string, location agentSkillLocation, lookPath executableLookup) (bool, error) {
	for _, dir := range location.projectDirs {
		projectDir := filepath.Join(root, dir)
		info, err := os.Stat(projectDir)
		if err == nil && info.IsDir() {
			return true, nil
		}
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			return false, fmt.Errorf("inspect %s project configuration: %w", location.name, err)
		}
	}

	for _, executable := range location.executables {
		if _, err := lookPath(executable); err == nil {
			return true, nil
		}
	}
	return false, nil
}

// installSkills validates every destination before writing any copy, then
// installs the same embedded skill content into all detected agent locations.
// It removes the old Cursor-only copy after the shared copy is durable.
func installSkills(root string, locations []agentSkillLocation, force bool) ([]skillInstallResult, string, error) {
	states := make([]skillTargetState, 0, len(locations))
	for _, location := range locations {
		target := filepath.Join(root, filepath.FromSlash(location.path))
		existing, err := os.ReadFile(target)
		switch {
		case err == nil && bytes.Equal(existing, diagramSkill):
			states = append(states, skillTargetState{location: location, path: target})
		case err == nil && !force:
			return nil, "", fmt.Errorf(
				"%s contains local changes; rerun with -force to replace it",
				target,
			)
		case err != nil && !errors.Is(err, os.ErrNotExist):
			return nil, "", fmt.Errorf("read existing %s diagram skill: %w", location.name, err)
		default:
			states = append(states, skillTargetState{
				location: location,
				path:     target,
				changed:  true,
			})
		}
	}

	legacyPath, removeLegacy, err := legacyCursorSkillMigration(root, locations, force)
	if err != nil {
		return nil, "", err
	}

	results := make([]skillInstallResult, 0, len(states))
	for _, state := range states {
		if state.changed {
			if err := os.MkdirAll(filepath.Dir(state.path), 0o755); err != nil {
				return nil, "", fmt.Errorf("create %s diagram skill directory: %w", state.location.name, err)
			}
			if err := os.WriteFile(state.path, diagramSkill, 0o644); err != nil {
				return nil, "", fmt.Errorf("write %s diagram skill: %w", state.location.name, err)
			}
		}
		results = append(results, skillInstallResult{
			agent:   state.location.name,
			path:    state.path,
			changed: state.changed,
		})
	}
	if removeLegacy {
		if err := os.Remove(legacyPath); err != nil {
			return nil, "", fmt.Errorf("remove duplicate legacy Cursor diagram skill: %w", err)
		}
		return results, legacyPath, nil
	}
	return results, "", nil
}

// legacyCursorSkillMigration validates an old Cursor-only copy before writes
// and authorizes its removal only when the shared Agent Skills target is used.
func legacyCursorSkillMigration(
	root string,
	locations []agentSkillLocation,
	force bool,
) (string, bool, error) {
	usesSharedTarget := false
	for _, location := range locations {
		if location.path == agentsDiagramSkillPath {
			usesSharedTarget = true
			break
		}
	}
	if !usesSharedTarget {
		return "", false, nil
	}

	legacyPath := filepath.Join(root, filepath.FromSlash(legacyCursorDiagramSkillPath))
	existing, err := os.ReadFile(legacyPath)
	switch {
	case errors.Is(err, os.ErrNotExist):
		return "", false, nil
	case err != nil:
		return "", false, fmt.Errorf("read legacy Cursor diagram skill: %w", err)
	case bytes.Equal(existing, diagramSkill) || force:
		return legacyPath, true, nil
	default:
		return "", false, fmt.Errorf(
			"%s contains local changes; rerun with -force to migrate it",
			legacyPath,
		)
	}
}
