package manager

type action string

const (
	Create action = "create"
	Update        = "update"
	Delete        = "delete"
)

type Mutation struct {
	Action      action
	Hostname    string
	ServiceName string
	Port        int
}
