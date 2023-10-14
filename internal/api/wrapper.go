package api

import (
	"context"
	"fmt"
)

// Get all volunteers
func (na *NocoApiStruct) AllVolunteers() (*[]VolunteersResponse, error) {
	res, err := na.clientWithResponses.VolunteersDbTableRowListWithResponse(context.Background(), &VolunteersDbTableRowListParams{})

	if err != nil {
		return nil, err
	}

	return res.JSON200.List, nil
}

// Get all events
func (na *NocoApiStruct) AllEvents() (*[]EventsResponse, error) {
	res, err := na.clientWithResponses.EventsDbTableRowListWithResponse(context.Background(), &EventsDbTableRowListParams{})

	if err != nil {
		return nil, err
	}

	return res.JSON200.List, nil
}

// Get all events where end date is greater than today
func (na *NocoApiStruct) AllUpcomingEvents() (*[]EventsResponse, error) {
	where := "(End,gt,today)"
	res, err := na.clientWithResponses.EventsDbTableRowListWithResponse(context.Background(), &EventsDbTableRowListParams{
		Where: &where,
	})

	if err != nil {
		return nil, err
	}

	return res.JSON200.List, nil
}

// Get a singleVolunteerByID
func (na *NocoApiStruct) VolunteerByID(volunteerID uint) (*VolunteersResponse, error) {
	where := fmt.Sprintf("(Id,eq,%d)", volunteerID)
	res, err := na.clientWithResponses.VolunteersDbTableRowFindOneWithResponse(context.Background(), &VolunteersDbTableRowFindOneParams{
		Where: &where,
	})

	if err != nil {
		return nil, err
	}

	return res.JSON200, nil
}

// Get a single event by ID
func (na *NocoApiStruct) EventByID(eventID uint) (*EventsResponse, error) {
	where := fmt.Sprintf("(Id,eq,%d)", eventID)
	res, err := na.clientWithResponses.EventsDbTableRowFindOneWithResponse(context.Background(), &EventsDbTableRowFindOneParams{
		Where: &where,
	})

	if err != nil {
		return nil, err
	}

	return res.JSON200, nil
}
