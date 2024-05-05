package entity

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClimaValido(t *testing.T) {

	w := NewWeather()
	w1, err := w.ColetaTempo("Sao Paulo")
	assert.Nil(t, err)
	assert.NotNil(t, w)
	assert.NotNil(t, w1)

	assert.Equal(t, "Brazil", w1.Location.Country)
	assert.Equal(t, "Sao Paulo", w1.Location.Name)

}

func TestClimaNaoEncontrado(t *testing.T) {

	w := NewWeather()
	w1, err := w.ColetaTempo("SaoPaulo")
	assert.NotNil(t, err)
	assert.NotNil(t, w)
	assert.Nil(t, w1)

	assert.Equal(t, http.StatusNotFound, err.StatusCode)
	assert.Equal(t, "clima não encontrado para a localidade SaoPaulo", err.Err.Error())
}
