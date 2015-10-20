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
