package ClientHandler

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

//
// ensureConfig проверяет наличие конфигурационного файла по указанному пути.
// Если файл отсутствует, создает его с дефолтными значениями.
//
// @param path путь к конфигурационному файлу
// @return указатель на Config и ошибка (если есть)
//
func ensureConfig(path string) (*Config, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		defaultCfg := Config{
			Ip:    "127.0.0.1:8080",
			Token: "your_token_here",
		}

		data, err := json.MarshalIndent(defaultCfg, "", "	")

		if err != nil {
			return nil, fmt.Errorf("ошибка маршалинга дефолтного конфига: %v", err)
		}

		if err := os.WriteFile(path, data, 0644); err != nil {
			return nil, fmt.Errorf("ошибка записи дефолтного конфига: %v", err)
		}

		return nil, fmt.Errorf("конфигурационный файл не найден, создан `config.json` с дефолтными значениями")
	}

	return loadConfig(path)
}

//
// loadConfig загружает конфигурацию из файла по указанному пути.
//
// @param path путь к конфигурационному файлу
// @return указатель на Config и ошибка (если есть)
//
func loadConfig(path string) (*Config, error) {
	file, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	var cfg Config

	if err := json.NewDecoder(file).Decode(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

//
// watchDocker отслеживает изменения docker-контейнеров и образов с заданным интервалом.
// При обнаружении изменений отправляет соответствующие сообщения через Communicator.
//
// @param c указатель на Communicator
// @param interval интервал проверки изменений
//
func watchDocker(c *Communicator, interval time.Duration) {
	prevCont := GetAllDockerContainers()
	prevImg  := GetAllDockerImages()
	ticker   := time.NewTicker(interval)
	defer ticker.Stop()

	prevCMap := make(map[string]*DockerContainer)
	for _, ctr := range prevCont {
		prevCMap[ctr.Hash] = ctr
	}

	prevImgMap := make(map[string]*DockerImage)
	for _, img := range prevImg {
		prevImgMap[img.Hash] = img
	}

	for range ticker.C {
		// Docker containers
		currCont := GetAllDockerContainers()

		currCMap := make(map[string]*DockerContainer)
		for _, ctr := range currCont {
			currCMap[ctr.Hash] = ctr
		}

		// Check for added, updated containers
		for id, ctr := range currCMap {
			if _, ok := prevCMap[id]; !ok {
				data, _ := json.Marshal(ctr)
				c.Requests.Add(&SentMessage{Type: AddedDockerContainer, Data: string(data)})
			}

			if prev, ok := prevCMap[id]; ok && (prev.Status != ctr.Status || prev.Recourses != ctr.Recourses) {
				data, _ := json.Marshal(ctr)
				c.Requests.Add(&SentMessage{Type: UpdatedDockerContainer, Data: string(data)})
			}
		}

		// Check for removed containers
		for id, ctr := range prevCMap {
			if _, ok := currCMap[id]; !ok {
				data, _ := json.Marshal(ctr)
				c.Requests.Add(&SentMessage{Type: RemovedDockerContainer, Data: string(data)})
			}
		}

		prevCMap = currCMap
		prevCont = currCont

		// Docker images
		currImg := GetAllDockerImages()

		currImgMap := make(map[string]*DockerImage)
		for _, img := range currImg {
			currImgMap[img.Hash] = img
		}

		// Check for added images
		for id, img := range currImgMap {
			if _, ok := prevImgMap[id]; !ok {
				data, _ := json.Marshal(img)
				c.Requests.Add(&SentMessage{Type: AddedDockerImage, Data: string(data)})
			}
		}

		// Check for removed images
		for id, img := range prevImgMap {
			if _, ok := currImgMap[id]; !ok {
				data, _ := json.Marshal(img)
				c.Requests.Add(&SentMessage{Type: RemovedDockerImage, Data: string(data)})
			}
		}

		prevImgMap = currImgMap
		prevImg = currImg
	}
}

//
// main является точкой входа в приложение.
//
func main() {
	cfg, err := ensureConfig("config.json")

	if err != nil {
		log.Println(err)
		os.Exit(0)
	}

	com := NewCommunicator(cfg.Token, cfg.Ip)

	com.Connect()
	com.StartHandlingThread()

	go watchDocker(com, 10 * time.Second)

	select {} // Держим приложение живым
}

