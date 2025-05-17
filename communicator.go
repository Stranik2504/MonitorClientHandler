package ClientHandler

import (
	"encoding/json"
	"log"
	"net/url"

	"github.com/gorilla/websocket"
)

// Communicator обеспечивает взаимодействие с сервером по WebSocket.
//
// @field Token токен авторизации
// @field Ip IP-адрес сервера
// @field Con WebSocket-соединение
// @field Requests очередь исходящих сообщений
type Communicator struct {
	Token    string
	Ip       string
	Con      *websocket.Conn
	Requests AtomicQueue[*SentMessage]
}

// NewCommunicator создает новый экземпляр Communicator.
//
// @param token токен авторизации
// @param ip IP-адрес сервера
// @return указатель на Communicator
func NewCommunicator(token string, ip string) *Communicator {
	return &Communicator{
		Token: token,
		Ip:    ip,
		Con:   nil,
	}
}

// SendMessage отправляет сообщение серверу по WebSocket.
//
// @param message указатель на отправляемое сообщение
func (c *Communicator) SendMessage(message *SentMessage) {
	err := c.Con.WriteJSON(message)

	if err != nil {
		log.Printf("Ошибка отправки сообщения: %v", err)
	}
}

// Connect устанавливает WebSocket-соединение с сервером и отправляет стартовое сообщение.
func (c *Communicator) Connect() {
	// Формируем URL для соединения по WebSocket
	u := url.URL{Scheme: "ws", Host: c.Ip, Path: "/ws"}
	log.Printf("Подключение к серверу: %s", u.String())

	// Устанавливаем соединение
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)

	if err != nil {
		log.Fatalf("Ошибка соединения: %v", err)
	}

	// Дальнейшая логика работы с соединением
	log.Println("Соединение установлено")
	c.Con = conn

	// Отправляем сообщение на сервер
	c.SendStartMessage()
}

// Close закрывает WebSocket-соединение.
func (c *Communicator) Close() {
	// Закрываем соединение
	if c.Con != nil {
		err := c.Con.Close()

		if err != nil {
			log.Printf("Ошибка закрытия соединения: %v", err)
		} else {
			log.Println("Соединение закрыто")
		}
	}
}

// StartHandlingThread запускает горутину для обработки входящих сообщений от сервера.
func (c *Communicator) StartHandlingThread() {
	// Запускаем горутину для обработки сообщений
	go func() {
		for {
			// Читаем сообщения из соединения
			_, message, err := c.Con.ReadMessage()

			if err != nil {
				log.Printf("Ошибка чтения сообщения: %v", err)
				break
			}

			var receiveMessage ReceiveMessage
			err = json.Unmarshal(message, &receiveMessage)

			if err != nil {
				log.Printf("Ошибка декодирования сообщения: %v", err)
				continue
			}

			if receiveMessage.Type == StartContainer {
				_, mess := StartDockerContainer(receiveMessage.Data)

				c.SendMessage(&SentMessage{
					Type: Result,
					Data: mess,
				})

				continue
			}

			if receiveMessage.Type == StopContainer {
				_, mess := StopDockerContainer(receiveMessage.Data)

				c.SendMessage(&SentMessage{
					Type: Result,
					Data: mess,
				})

				continue
			}

			if receiveMessage.Type == RemoveContainer {
				_, mess := RemoveDockerContainer(receiveMessage.Data)

				c.SendMessage(&SentMessage{
					Type: Result,
					Data: mess,
				})

				continue
			}

			if receiveMessage.Type == RemoveImage {
				_, mess := RemoveDockerImage(receiveMessage.Data)

				c.SendMessage(&SentMessage{
					Type: Result,
					Data: mess,
				})

				continue
			}

			if receiveMessage.Type == RunScript {
				out, err := runScript(receiveMessage.Data)

				if err != nil {
					log.Printf("Ошибка выполнения скрипта: %v", err)
				}

				c.SendMessage(&SentMessage{Type: Result, Data: out})
				continue
			}

			if receiveMessage.Type == RunCommand {
				out, err := runCommand(receiveMessage.Data)

				if err != nil {
					log.Printf("Ошибка выполнения команды: %v", err)
				}

				c.SendMessage(&SentMessage{Type: Result, Data: out})
				continue
			}

			if receiveMessage.Type == Restart {
				err := reboot()

				if err != nil {
					log.Printf("Ошибка перезагрузки: %v", err)
				}

				c.SendMessage(&SentMessage{
					Type: Restarted,
					Data: "Ok",
				})

				continue
			}

			if c.Requests.Size() > 0 {
				message := c.Requests.Pop()
				c.SendMessage(message)
				continue
			}

			data, _ := json.Marshal(getMetric())

			c.SendMessage(&SentMessage{
				Type: SendMetric,
				Data: string(data),
			})
		}
	}()
}

// SendStartMessage формирует и отправляет стартовое сообщение серверу.
func (c *Communicator) SendStartMessage() {
	message := SentStartMessage{
		Type:             Start,
		Token:            c.Token,
		Metric:           getMetric(),
		DockerImages:     GetAllDockerImages(),
		DockerContainers: GetAllDockerContainers(),
	}

	// Отправляем сообщение на сервер
	err := c.Con.WriteJSON(message)

	if err != nil {
		log.Printf("Ошибка отправки сообщения: %v", err)
	}
}
