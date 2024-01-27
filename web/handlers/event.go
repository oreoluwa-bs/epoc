package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/oreoluwa-bs/epoc/models"
)

type EventHandler struct {
	DB *sql.DB
}

func (eh *EventHandler) Create(w http.ResponseWriter, r *http.Request) {
	eventStore := models.EventStore{DB: eh.DB}

	var data models.CreateEvent

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&data)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	newEvent, err := eventStore.Create(data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	eventJson, err := json.Marshal(newEvent)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(eventJson)

}

func (eh *EventHandler) List(w http.ResponseWriter, r *http.Request) {
	eventStore := models.EventStore{DB: eh.DB}

	queryParams := r.URL.Query()

	filters := models.GetAllEventFilters{}

	startsQuery := queryParams.Get("starts_at")
	if len(startsQuery) > 0 {
		startsAt, err := time.Parse("2006-01-02T15:04:05Z07:00", startsQuery)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		filters.StartsAt = startsAt
	}

	endsQuery := queryParams.Get("ends_at")
	if len(endsQuery) > 0 {

		endsAt, err := time.Parse("2006-01-02T15:04:05Z07:00", endsQuery)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		filters.EndsAt = endsAt
	}

	events, err := eventStore.GetAll(filters)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	data, err := json.Marshal(events)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)

}

func (eh *EventHandler) GetById(w http.ResponseWriter, r *http.Request) {
	eventStore := models.EventStore{DB: eh.DB}

	eventId, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid event Id"))
		return
	}

	events, err := eventStore.GetById(eventId)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	data, err := json.Marshal(events)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (eh *EventHandler) UpdateById(w http.ResponseWriter, r *http.Request) {
	eventStore := models.EventStore{DB: eh.DB}

	eventId, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid event Id"))
		return
	}

	var body models.UpdateEvent

	dec := json.NewDecoder(r.Body)
	err = dec.Decode(&body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	body.Id = eventId

	events, err := eventStore.UpdateById(body)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	data, err := json.Marshal(events)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (eh *EventHandler) DeleteById(w http.ResponseWriter, r *http.Request) {
	eventStore := models.EventStore{DB: eh.DB}

	eventId, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid event Id"))
		return
	}

	var body models.UpdateEvent

	dec := json.NewDecoder(r.Body)
	err = dec.Decode(&body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	err = eventStore.DeleteById(eventId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Deleted"))
}
