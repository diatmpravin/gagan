package models

type Organization struct {
	Name string
	Guid string
}

type Space struct {
	Name string
	Guid string
}

type Application struct {
	Name      string
	Guid      string
	State     string
	Instances int
	Memory    int
	Urls      []string
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
	Label string
	Guid  string
	Plans []ServicePlan
}

type ServicePlan struct {
	Name string
	Guid string
}
