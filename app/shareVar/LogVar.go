package shareVar

import "github.com/sirupsen/logrus"

var Logger *logrus.Logger

type LayerName string

const (
	Controller LayerName = "Controller"
	Service    LayerName = "Service"
	Repository LayerName = "Repository"
)

type NameEntity string

const (
	Project    NameEntity = "Project"
	SemgrepCli NameEntity = "SemgrepCli"
	GitCli     NameEntity = "GitCli"
	AsocAPI    NameEntity = "AsocAPI"
)

type LogBasicKeyField string

const (
	Request     LogBasicKeyField = "request"
	Response    LogBasicKeyField = "response"
	TimeElapsed LogBasicKeyField = "time_elapsed"
)
