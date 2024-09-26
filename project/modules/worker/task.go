package worker

type Task int

const (
	TaskIssueReceipt Task = iota
	TaskAppendToTracker
	TicketBookingConfirmed
	TicketBookingCanceled
)

var TopicsMap = map[Task]string{
	TaskIssueReceipt:       "issue-receipt",
	TaskAppendToTracker:    "append-to-tracker",
	TicketBookingCanceled:  "TicketBookingCanceled",
	TicketBookingConfirmed: "TicketBookingConfirmed",
}
