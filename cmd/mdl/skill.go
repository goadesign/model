// The skill installer publishes MDL's canonical diagram-editing instructions
// into the current repository. It preserves locally edited skill files unless
// the caller explicitly authorizes replacement.
package main

import (
	"bytes"
	_ "embed"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

const diagramSkillPath = ".cursor/skills/editing-model-diagrams/SKILL.md"

//go:embed skills/editing-model-diagrams/SKILL.md
var diagramSkill []byte

func runSkillInstall(force bool) error {
	root, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("determine current repository: %w", err)
	}

	path, changed, err := installSkill(root, force)
	if err != nil {
		return err
	}
	if changed {
		fmt.Printf("Installed MDL diagram skill at %s\n", path)
	} else {
		fmt.Printf("MDL diagram skill is already current at %s\n", path)
	}
	return nil
}

func installSkill(root string, force bool) (string, bool, error) {
	target := filepath.Join(root, filepath.FromSlash(diagramSkillPath))
	existing, err := os.ReadFile(target)
	switch {
	case err == nil && bytes.Equal(existing, diagramSkill):
		return target, false, nil
	case err == nil && !force:
		return "", false, fmt.Errorf(
			"%s contains local changes; rerun with -force to replace it",
			target,
		)
	case err != nil && !errors.Is(err, os.ErrNotExist):
		return "", false, fmt.Errorf("read existing diagram skill: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return "", false, fmt.Errorf("create diagram skill directory: %w", err)
	}
	if err := os.WriteFile(target, diagramSkill, 0o644); err != nil {
		return "", false, fmt.Errorf("write diagram skill: %w", err)
	}
	return target, true, nil
}
