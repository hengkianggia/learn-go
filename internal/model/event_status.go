package model

import "fmt"

func (e EventStatus) IsValid() error {
	switch e {
	case Draft, Published, Cancelled:
		return nil
	}
	return fmt.Errorf("invalid event status: %s", e)
}
