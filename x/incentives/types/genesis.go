package types

import "fmt"

// DefaultGenesisState returns the default genesis state of the incentives module
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		NextScheduleId: 1,
		Schedules:      []Schedule{},
	}
}

// ValidateGenesis validates the given instance of the incentives module's
// genesis state.
//
// for each schedule,
//
// - the id must be smaller than the next schedule id
//
// - the id must not be duplicate
//
// - the end time must be after the start time
//
// - the total amount must be non-zero
//
// - the released amount must be equal or smaller than the total amount
func (gs GenesisState) Validate() error {
	seenIds := make(map[uint64]bool)
	for _, schedule := range gs.Schedules {
		if schedule.Id >= gs.NextScheduleId {
			return fmt.Errorf("incentives schedule id %d is not smaller than next schedule id %d", schedule.Id, gs.NextScheduleId)
		}

		if seenIds[schedule.Id] {
			return fmt.Errorf("incentives schedule has duplicate id %d", schedule.Id)
		}

		if !schedule.EndTime.After(schedule.StartTime) {
			return fmt.Errorf("incentives schedule %d end time is not after start time", schedule.Id)
		}

		if schedule.TotalAmount.Empty() {
			return fmt.Errorf("incentives schedule %d has zero total amount", schedule.Id)
		}

		if !schedule.TotalAmount.IsAllGTE(schedule.ReleasedAmount) {
			return fmt.Errorf("incentives schedule %d total amount is not all greater or equal than released amount", schedule.Id)
		}

		seenIds[schedule.Id] = true
	}

	return nil
}
