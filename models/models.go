package models

import (
	"time"
)

type Organization struct {
	Name string
	Guid string
}

type Domain struct {
	Name string
	Guid string
}

type Route struct {
	Host string
	Guid string
}

type InstanceState string

const (
	InstanceStarting InstanceState = "starting"
	InstanceRunning                = "running"
	InstanceFlapping               = "flapping"
	InstanceDown                   = "down"
)

type ApplicationInstance struct {
	State InstanceState
}

type ServiceOffering struct {
	Guid     string
	Label    string
	Provider string
	Version  string
	Plans    []ServicePlan
}

type ServicePlan struct {
	Name            string
	Guid            string
	ServiceOffering ServiceOffering
}

type ServiceBinding struct {
	Url     string
	Guid    string
	AppGuid string
}

type Space struct {
	Name             string
	Guid             string
	Applications     []Application
	ServiceInstances []ServiceInstance
}

type Application struct {
	Name             string
	Guid             string
	State            string
	Instances        int
	RunningInstances int
	Memory           int
	Urls             []string
	BuildpackUrl     string
	Stack            Stack
}

type ServiceInstance struct {
	Name             string
	Guid             string
	ServiceBindings  []ServiceBinding
	ServicePlan      ServicePlan
	ApplicationNames []string
}

type Stack struct {
	Name string
	Guid string
}

type EventFields struct {
	Guid        string
	Name        string
	Timestamp   time.Time
	Description string
	ActorName   string
}
