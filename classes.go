package main

import (
	"time"
)

// TypeReceivedMessage определяет типы входящих сообщений.
type TypeReceivedMessage int

// TypeSentMessage определяет типы исходящих сообщений.
type TypeSentMessage int

// Константы для типов входящих сообщений.
const (
	StartContainer  TypeReceivedMessage = iota // Запуск контейнера
	StopContainer                              // Остановка контейнера
	RemoveContainer                            // Удаление контейнера
	RemoveImage                                // Удаление образа
	RunScript                                  // Выполнение скрипта
	RunCommand                                 // Выполнение команды
	Restart                                    // Перезапуск
	Ok                                         // Подтверждение
)

// Константы для типов исходящих сообщений.
const (
	SendMetric             TypeSentMessage = iota // Отправка метрик
	AddedDockerImage                              // Добавлен docker-образ
	AddedDockerContainer                          // Добавлен docker-контейнер
	RemovedDockerImage                            // Удалён docker-образ
	RemovedDockerContainer                        // Удалён docker-контейнер
	UpdatedDockerContainer                        // Обновлён docker-контейнер
	Start                                         // Старт
	Result                                        // Результат
	Restarted                                     // Перезапущено
	None                                          // Нет действия
)

// SentStartMessage представляет сообщение о запуске, отправляемое клиенту.
//
// @field Type тип исходящего сообщения
// @field Token токен авторизации
// @field Metric метрика системы
// @field DockerImages список docker-образов
// @field DockerContainers список docker-контейнеров
type SentStartMessage struct {
	Type             TypeSentMessage
	Token            string
	Metric           *Metric
	DockerImages     []*DockerImage
	DockerContainers []*DockerContainer
}

// Metric содержит информацию о метриках системы.
//
// @field Cpus загрузка процессоров
// @field UseRam используемая оперативная память
// @field TotalRam всего оперативной памяти
// @field UseDisks используемое место на дисках
// @field TotalDisks общий объём дисков
// @field NetworkSend отправлено по сети
// @field NetworkReceive получено по сети
// @field Time время снятия метрик
type Metric struct {
	Cpus           []float64
	UseRam         int
	TotalRam       int
	UseDisks       []int
	TotalDisks     []int
	NetworkSend    int
	NetworkReceive int
	Time           time.Time
}

// DockerImage описывает docker-образ.
//
// @field Id идентификатор образа
// @field Name имя образа
// @field Size размер образа
// @field Hash хеш образа
type DockerImage struct {
	Id   int
	Name string
	Size float64
	Hash string
}

// DockerContainer описывает docker-контейнер.
//
// @field Id идентификатор контейнера
// @field Name имя контейнера
// @field ImageId идентификатор образа
// @field ImageHash хеш образа
// @field Status статус контейнера
// @field Recourses ресурсы контейнера
// @field Hash хеш контейнера
type DockerContainer struct {
	Id        int
	Name      string
	ImageId   int
	ImageHash string
	Status    string
	Recourses string
	Hash      string
}

// ReceiveMessage представляет входящее сообщение.
//
// @field Type тип входящего сообщения
// @field Data данные сообщения
type ReceiveMessage struct {
	Type TypeReceivedMessage
	Data string
}

// SentMessage представляет исходящее сообщение.
//
// @field Type тип исходящего сообщения
// @field Data данные сообщения
type SentMessage struct {
	Type TypeSentMessage
	Data string
}

// Config содержит конфигурацию клиента.
//
// @field Ip IP-адрес сервера
// @field Token токен авторизации
type Config struct {
	Ip    string
	Token string
}
