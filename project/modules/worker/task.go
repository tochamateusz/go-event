package worker

type Task int

const (
	TaskIssueReceipt Task = iota
	TaskAppendToTracker
)

var TopicsMap = map[Task]string{
	TaskIssueReceipt:    "issue-receipt",
	TaskAppendToTracker: "append-to-tracker",
}

type Message struct {
	Task     Task
	TicketID string
}
