package delivery

import "time"

type DeadlineFactory struct {}

func NewDeadlineFactory() *DeadlineFactory {
	return &DeadlineFactory{}
}

func (f *DeadlineFactory) Deadline(now time.Time, transportType string) time.Time {
    switch transportType {
    case "on_foot":
        return now.Add(30 * time.Minute)
    case "scooter":
        return now.Add(15 * time.Minute)
    case "car":
        return now.Add(5 * time.Minute)
    default:
        return now.Add(30 * time.Minute)
    }
}