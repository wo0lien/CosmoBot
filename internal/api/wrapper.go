package api

import (
	"context"
)

// Get all volunteers
func (na *NocoApiStruct) GetAllVolunteers() (*[]VolunteersResponse, error) {
	res, err := na.ClientWithResponses.VolunteersDbTableRowListWithResponse(context.Background(), &VolunteersDbTableRowListParams{})

	if err != nil {
		return nil, err
	}

	return res.JSON200.List, nil
}

// Get all events
func (na *NocoApiStruct) GetAllEvents() (*[]EventsResponse, error) {
	res, err := na.ClientWithResponses.EventsDbTableRowListWithResponse(context.Background(), &EventsDbTableRowListParams{})

	if err != nil {
		return nil, err
	}

	return res.JSON200.List, nil
}

// Get all events where end date is greater than today
func (na *NocoApiStruct) GetAllUpcomingEvents() (*[]EventsResponse, error) {
	where := "(End,gt,today)"
	res, err := na.ClientWithResponses.EventsDbTableRowListWithResponse(context.Background(), &EventsDbTableRowListParams{
		Where: &where,
	})

	if err != nil {
		return nil, err
	}

	return res.JSON200.List, nil
}
