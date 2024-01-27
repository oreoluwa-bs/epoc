package models

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type Event struct {
	Id          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	StartsAt    time.Time `json:"starts_at"`
	EndsAt      time.Time `json:"ends_at"`
}

type EventStore struct {
	DB *sql.DB
}

type CreateEvent struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	StartsAt    time.Time `json:"starts_at"`
	EndsAt      time.Time `json:"ends_at"`
}

func ValidateCreateEvent(event CreateEvent) error {
	if len(event.Name) < 1 {
		return errors.New("event name is required")
	}
	if !event.StartsAt.Before(event.EndsAt) {
		return errors.New("start date must be before the end date")
	}

	return nil
}

func (es *EventStore) Create(event CreateEvent) (Event, error) {
	if err := ValidateCreateEvent(event); err != nil {
		return Event{}, err
	}

	query := `INSERT INTO events (name, description, starts_at, ends_at) VALUES (?, ?, ?, ?)`

	result, err := es.DB.Exec(query, event.Name, event.Description, event.StartsAt, event.EndsAt)
	if err != nil {
		return Event{}, err
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		return Event{}, err
	}

	newEvent := Event{
		Id:          int(lastID),
		Name:        event.Name,
		Description: event.Description,
		StartsAt:    event.StartsAt,
		EndsAt:      event.EndsAt,
	}

	return newEvent, nil

}

type GetAllEventFilters struct {
	StartsAt time.Time `json:"starts_at"`
	EndsAt   time.Time `json:"ends_at"`
}

func (es *EventStore) GetAll(filters GetAllEventFilters) ([]Event, error) {

	query := `SELECT id, name, description, starts_at, ends_at FROM events WHERE id IS NOT NULL`

	args := []interface{}{}
	if !filters.StartsAt.IsZero() {
		query += " AND starts_at >= ?"
		args = append(args, filters.StartsAt)
	}

	if !filters.EndsAt.IsZero() {
		query += " AND ends_at <= ?"
		args = append(args, filters.EndsAt)
	}

	fmt.Println(query, args)

	rows, err := es.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}

	var events = make([]Event, 0)

	for rows.Next() {
		event := Event{}
		err := rows.Scan(&event.Id, &event.Name, &event.Description, &event.StartsAt, &event.EndsAt)
		if err != nil {
			return nil, err
		}

		events = append(events, event)
	}

	if rows.Err() != nil {
		return nil, err

	}

	return events, nil

}

func (es *EventStore) GetById(id int) (Event, error) {

	query := `SELECT id, name, description, starts_at, ends_at FROM events WHERE id=?`

	fmt.Println(query, id)

	row := es.DB.QueryRow(query, id)

	event := Event{}
	err := row.Scan(&event.Id, &event.Name, &event.Description, &event.StartsAt, &event.EndsAt)
	if err != nil {
		return Event{}, err
	}

	if row.Err() != nil {
		return Event{}, err

	}

	return event, nil
}

type UpdateEvent struct {
	Id          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	StartsAt    time.Time `json:"starts_at"`
	EndsAt      time.Time `json:"ends_at"`
}

func validateUpdateEvent(event Event) error {
	if len(event.Name) < 1 {
		return errors.New("event name is required")
	}

	if !event.StartsAt.Before(event.EndsAt) {
		return errors.New("start date must be before the end date")
	}

	return nil
}

func (es *EventStore) UpdateById(eventData UpdateEvent) (Event, error) {

	_, err := es.GetById(eventData.Id)
	if err != nil {
		return Event{}, err
	}

	if err := validateUpdateEvent(Event(eventData)); err != nil {
		return Event{}, err
	}

	query := `UPDATE events SET name=?, description=?, starts_at=?, ends_at=? WHERE id=?`

	args := []interface{}{
		eventData.Name,
		eventData.Description,
		eventData.StartsAt,
		eventData.EndsAt,
		eventData.Id,
	}

	fmt.Println(query, args)

	statement, err := es.DB.Prepare(query)
	if err != nil {
		return Event{}, err
	}

	_, err = statement.Exec(args...)
	if err != nil {
		return Event{}, err
	}

	return Event(eventData), nil
}

func (es *EventStore) DeleteById(id int) error {

	existing, err := es.GetById(id)
	if err != nil {
		return err
	}

	query := `DELETE FROM events WHERE id=?`

	args := []interface{}{
		existing.Id,
	}

	fmt.Println(query, args)

	statement, err := es.DB.Prepare(query)
	if err != nil {
		return err
	}

	_, err = statement.Exec(args...)
	if err != nil {
		return err
	}

	return nil
}
