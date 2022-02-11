package data

type Project struct {
	Version string `json:"version"`
	Author  string `json:"author"`
	Name    string `json:"name"`
}

func Home() *Project {
	return &Project{
		Version: "1.0",
		Author:  "Herbie",
		Name:    "Training-notebook api",
	}
}
