package main

import (
	"time"
)

type TypeReceivedMessage int
type TypeSentMessage int

const (
	StartContainer TypeReceivedMessage = iota
	StopContainer
	RemoveContainer
	RemoveImage
	RunScript
	RunCommand
	Restart
	Ok
)

const (
	SendMetric TypeSentMessage = iota
	AddedDockerImage
	AddedDockerContainer
	RemovedDockerImage
	RemovedDockerContainer
	UpdatedDockerContainer
	Start
	Result
	Restarted
	None
)

type SentStartMessage struct {
	Type             TypeSentMessage
	Token            string
	Metric           *Metric
	DockerImages     []*DockerImage
	DockerContainers []*DockerContainer
}

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

type DockerImage struct {
	Id   int
	Name string
	Size float64
	Hash string
}

type DockerContainer struct {
	Id        int
	Name      string
	ImageId   int
	ImageHash string
	Status    string
	Recourses string
	Hash      string
}

type ReceiveMessage struct {
	Type TypeReceivedMessage
	Data string
}

type SentMessage struct {
	Type TypeSentMessage
	Data string
}

type Config struct {
	Ip    string
	Token string
}
