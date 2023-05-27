package todo

type Task struct {
	ID string `json: "id"`
	Order int `json: "order"`
}

type List struct {
	ID    string `json: "id"`
	Title string `json: "task"`
	Order int    `json: "order"`
	Tasks []Task `json: "task"`
}

// global variables for todo package
var (
	tasks []Task
	lists []List
)
