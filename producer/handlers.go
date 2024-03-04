package main

import (
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"go-kafka/producer/data"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type handler func(w http.ResponseWriter, r *http.Request) error

func (app *application) appHandler(f handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			app.logger.Error(err.Error())
			writeJSON(w, http.StatusInternalServerError, struct {
				Error string `json:"error"`
			}{Error: err.Error()}, nil)

		}
	}
}

func (app *application) produce(w http.ResponseWriter, r *http.Request) error {
	var input struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		return err
	}
	defer r.Body.Close()

	a := &data.Actor{
		ID:        rand.Intn(1000),
		Name:      input.Name,
		CreatedAt: time.Now().UTC(),
	}

	js, err := json.Marshal(a)
	if err != nil {
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: "actor",
		Key:   sarama.StringEncoder(strconv.Itoa(a.ID)),
		Value: sarama.StringEncoder(js),
	}
	partition, offset, err := app.producer.SendMessage(msg)
	if err != nil {
		app.logger.Error(err.Error())
	}

	app.logger.Info(fmt.Sprintf("Sent to partion %v and the offset is %v", partition, offset))

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/produce/%s", a.Name))
	return writeJSON(w, http.StatusCreated, a, nil)
}
