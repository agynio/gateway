package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

type moduleInfo struct {
	Dir string `json:"Dir"`
}

func main() {
	moduleDir, err := locateModuleDir("github.com/agynio/api")
	if err != nil {
		panic(fmt.Errorf("locate module: %w", err))
	}

	srcRoot := filepath.Join(moduleDir, "openapi", "team", "v1")
	if _, err := os.Stat(srcRoot); err != nil {
		panic(fmt.Errorf("stat source: %w", err))
	}

	repoRoot := moduleRoot()
	dstRoot := filepath.Join(repoRoot, "internal", "apischema", "teamv1")
	if err := os.RemoveAll(dstRoot); err != nil {
		panic(fmt.Errorf("remove destination: %w", err))
	}

	if err := copyDirectory(srcRoot, dstRoot); err != nil {
		panic(fmt.Errorf("copy spec: %w", err))
	}
}

func moduleRoot() string {
	_, current, _, ok := runtime.Caller(0)
	if !ok {
		cwd, err := os.Getwd()
		if err != nil {
			panic(fmt.Errorf("determine cwd: %w", err))
		}
		return cwd
	}
	return filepath.Clean(filepath.Join(filepath.Dir(current), "..", "..", ".."))
}

func locateModuleDir(path string) (string, error) {
	cmd := exec.Command("go", "list", "-m", "-json", path)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	var info moduleInfo
	if err := json.Unmarshal(output, &info); err != nil {
		return "", err
	}
	if info.Dir == "" {
		return "", fmt.Errorf("module %s missing dir", path)
	}
	return info.Dir, nil
}

func copyDirectory(src, dst string) error {
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dst, 0o755); err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		info, err := entry.Info()
		if err != nil {
			return err
		}

		switch {
		case info.Mode().IsDir():
			if err := copyDirectory(srcPath, dstPath); err != nil {
				return err
			}
		case info.Mode().IsRegular():
			if err := copyFile(srcPath, dstPath, info.Mode()); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unsupported file type: %s", srcPath)
		}
	}

	return nil
}

func copyFile(src, dst string, mode fs.FileMode) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}

	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, mode)
	if err != nil {
		return err
	}
	defer func() { _ = out.Close() }()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}

	return nil
}
