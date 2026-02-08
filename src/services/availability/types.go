package availability

import (
	"fmt"
	"time"
)

type AvailabilityResponse struct {
	VideoID int    `json:"video_id"`
	Window  Window `json:"availability_window"`
}

type Window struct {
	From string `json:"from"`
	To   string `json:"to"`
}

func (r AvailabilityResponse) MapToDomain() (Availability, error) {
	const layout = "2006-01-02T15:04:05.000"
	from, err := time.Parse(layout, r.Window.From)
	if err != nil {
		return Availability{}, fmt.Errorf("invalid 'from' date %q: %w", r.Window.From, err)
	}

	to, err := time.Parse(layout, r.Window.To)
	if err != nil {
		return Availability{}, fmt.Errorf("invalid 'to' date %q: %w", r.Window.To, err)
	}
	//validate availability window makes sense
	if to.Before(from) {
		return Availability{}, fmt.Errorf(
			"invalid availability window: 'to' date %s is before 'from' date %s",
			r.Window.To,
			r.Window.From,
		)
	}
	return Availability{
		VideoID: r.VideoID,
		From:    from,
		To:      to,
	}, nil

}
