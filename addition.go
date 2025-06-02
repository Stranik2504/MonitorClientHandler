package main

import (
	"log"
	"os"
	"os/exec"
	"runtime"
)

// Sum вычисляет сумму элементов среза numbers типа T.
//
// @param numbers срез чисел типа T (int или float64)
// @return сумма всех элементов среза
func Sum[T ~int | ~float64](numbers []T) T {
	var total T
	for _, n := range numbers {
		total += n
	}
	return total
}

// runScriptLinux выполняет переданный скрипт shell в Linux/Unix-системах.
//
// @param script строка с shell-скриптом
// @return вывод скрипта и ошибка (если есть)
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

// runScriptWindows выполняет переданный скрипт в Windows через .bat файл.
//
// @param script строка с bat-скриптом
// @return вывод скрипта и ошибка (если есть)
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

// runScript выполняет скрипт в зависимости от ОС (Windows или Linux).
//
// @param script строка скрипта
// @return вывод скрипта и ошибка (если есть)
func runScript(script string) (string, error) {
	if runtime.GOOS == "windows" {
		return runScriptWindows(script)
	}
	return runScriptLinux(script)
}

// runCommand выполняет команду cmdStr в командной строке ОС.
//
// @param cmdStr строка команды
// @return вывод команды и ошибка (если есть)
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

// reboot перезагружает компьютер в зависимости от ОС.
//
// @return ошибка, если перезагрузка не удалась
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
