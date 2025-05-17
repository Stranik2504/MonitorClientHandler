package ClientHandler

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"log"
	"strings"
)

// GetAllDockerImages возвращает список всех docker-образов на хосте.
//
// @return срез указателей на DockerImage
func GetAllDockerImages() []*DockerImage {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())

	if err != nil {
		log.Fatalf("Ошибка создания Docker-клиента: %v", err)
	}

	images, err := cli.ImageList(context.Background(), image.ListOptions{})

	if err != nil {
		log.Fatalf("Ошибка получения списка образов: %v", err)
	}

	imgs := make([]*DockerImage, len(images))

	for i, img := range images {
		imgs[i] = &DockerImage{
			Id:   0,
			Name: img.RepoTags[0],
			Size: float64(img.Size),
			Hash: img.ID,
		}
	}

	return imgs
}

// GetAllDockerContainers возвращает список всех docker-контейнеров на хосте.
//
// @return срез указателей на DockerContainer
func GetAllDockerContainers() []*DockerContainer {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())

	if err != nil {
		log.Fatalf("Ошибка создания Docker-клиента: %v", err)
	}

	containers, err := cli.ContainerList(context.Background(), container.ListOptions{})

	if err != nil {
		log.Fatalf("Ошибка получения списка контейнеров: %v", err)
	}

	conts := make([]*DockerContainer, len(containers))

	for i, cont := range containers {
		conts[i] = &DockerContainer{
			Id:        0,
			Name:      strings.Join(cont.Names, ""),
			ImageId:   0,
			ImageHash: cont.ImageID,
			Status:    cont.State,
			Recourses: "",
			Hash:      cont.ID,
		}
	}

	return conts
}

// StopDockerContainer останавливает контейнер по его хешу.
//
// @param containerHash хеш контейнера
// @return true и пустая строка при успехе, иначе false и сообщение об ошибке
func StopDockerContainer(containerHash string) (bool, string) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())

	if err != nil {
		log.Fatalf("Ошибка создания Docker-клиента: %v", err)
		return false, fmt.Sprintf("Ошибка создания Docker-клиента: %v", err)
	}

	err = cli.ContainerStop(context.Background(), containerHash, container.StopOptions{})

	if err != nil {
		log.Fatalf("Ошибка остановки контейнера: %v", err)
		return false, fmt.Sprintf("Ошибка остановки контейнера: %v", err)
	}

	return true, ""
}

// RemoveDockerContainer удаляет контейнер по его хешу.
//
// @param containerHash хеш контейнера
// @return true и пустая строка при успехе, иначе false и сообщение об ошибке
func RemoveDockerContainer(containerHash string) (bool, string) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())

	if err != nil {
		log.Fatalf("Ошибка создания Docker-клиента: %v", err)
		return false, fmt.Sprintf("Ошибка создания Docker-клиента: %v", err)
	}

	err = cli.ContainerRemove(context.Background(), containerHash, container.RemoveOptions{Force: true})

	if err != nil {
		log.Fatalf("Ошибка удаления контейнера: %v", err)
		return false, fmt.Sprintf("Ошибка удаления контейнера: %v", err)
	}

	return true, ""
}

// RemoveDockerImage удаляет docker-образ по его хешу.
//
// @param imageHash хеш образа
// @return true и пустая строка при успехе, иначе false и сообщение об ошибке
func RemoveDockerImage(imageHash string) (bool, string) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())

	if err != nil {
		log.Fatalf("Ошибка создания Docker-клиента: %v", err)
		return false, fmt.Sprintf("Ошибка создания Docker-клиента: %v", err)
	}

	_, err = cli.ImageRemove(context.Background(), imageHash, image.RemoveOptions{Force: true})

	if err != nil {
		log.Fatalf("Ошибка удаления образа: %v", err)
		return false, fmt.Sprintf("Ошибка удаления образа: %v", err)
	}

	return true, ""
}

// StartDockerContainer запускает контейнер по его хешу.
//
// @param containerHash хеш контейнера
// @return true и пустая строка при успехе, иначе false и сообщение об ошибке
func StartDockerContainer(containerHash string) (bool, string) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())

	if err != nil {
		log.Fatalf("Ошибка создания Docker-клиента: %v", err)
		return false, fmt.Sprintf("Ошибка создания Docker-клиента: %v", err)
	}

	err = cli.ContainerStart(context.Background(), containerHash, container.StartOptions{})

	if err != nil {
		log.Fatalf("Ошибка запуска контейнера: %v", err)
		return false, fmt.Sprintf("Ошибка запуска контейнера: %v", err)
	}

	return true, ""
}
