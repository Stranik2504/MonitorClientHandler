package main

import (
	"log"
	"os"
	"os/exec"
	"runtime"
)

func Sum[T ~int | ~float64](numbers []T) T {
	var total T
	for _, n := range numbers {
		total += n
	}
	return total
}

func runScriptLinux(script string) (string, error) {
	tmp, err := os.CreateTemp("", "*.sh")
	if err != nil {
		return "", err
	}
	defer func(name string) {
		err := os.Remove(name)

		if err != nil {
			log.Printf("ошибка удаления временного файла: %v", err)
		}
	}(tmp.Name())

	if _, err := tmp.WriteString(script); err != nil {
		return "", err
	}

	err = tmp.Close()

	if err != nil {
		return "", err
	}

	if err := os.Chmod(tmp.Name(), 0o700); err != nil {
		return "", err
	}

	cmd := exec.Command(tmp.Name())
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func runScriptWindows(script string) (string, error) {
	tmp, err := os.CreateTemp("", "*.bat")
	if err != nil {
		return "", err
	}
	defer func() {
		if rmErr := os.Remove(tmp.Name()); rmErr != nil {
			log.Printf("ошибка удаления временного файла: %v", rmErr)
		}
	}()

	content := "@echo off\r\n" + script
	if _, err := tmp.WriteString(content); err != nil {
		_ = tmp.Close()
		return "", err
	}

	err = tmp.Close()
	if err != nil {
		return "", err
	}

	cmd := exec.Command("cmd", "/C", tmp.Name())
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func runScript(script string) (string, error) {
	if runtime.GOOS == "windows" {
		return runScriptWindows(script)
	}
	return runScriptLinux(script)
}

func runCommand(cmdStr string) (string, error) {
	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", cmdStr)
	} else {
		cmd = exec.Command("sh", "-c", cmdStr)
	}

	out, err := cmd.CombinedOutput()

	return string(out), err
}

func reboot() error {
	// Перезагрузка Linux/Unix
	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		cmd = exec.Command("shutdown", "/r", "/t", "0")
	} else {
		cmd = exec.Command("shutdown", "-r", "now")
	}

	return cmd.Run()
}
