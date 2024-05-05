package entity

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCEPValido(t *testing.T) {

	local := NewLocalidadeViaCep()
	localidade, err := local.ColetaLocalidade("04112080")
	assert.Nil(t, err)
	assert.NotNil(t, local)
	assert.NotNil(t, localidade)

	assert.Equal(t, "São Paulo", localidade.Localidade)

	localidade, err = local.ColetaLocalidade("04112-080")
	assert.Nil(t, err)
	assert.NotNil(t, local)
	assert.NotNil(t, localidade)
}

func TestCEPInvalido(t *testing.T) {

	local := NewLocalidadeViaCep()
	localidade, err := local.ColetaLocalidade("0411208A")
	assert.NotNil(t, err)
	assert.Nil(t, localidade)
	assert.Equal(t, http.StatusUnprocessableEntity, err.StatusCode)
	assert.Equal(t, "invalid zipcode", err.Err.Error())

	localidade, err = local.ColetaLocalidade("041120800")
	assert.NotNil(t, err)
	assert.Nil(t, localidade)
	assert.Equal(t, http.StatusUnprocessableEntity, err.StatusCode)
	assert.Equal(t, "invalid zipcode", err.Err.Error())
}

func TestCEPNaoEncontrado(t *testing.T) {

	local := NewLocalidadeViaCep()
	localidade, err := local.ColetaLocalidade("04112089")
	assert.NotNil(t, err)
	assert.Nil(t, localidade)
	assert.Equal(t, http.StatusNotFound, err.StatusCode)
	assert.Equal(t, "can not find zipcode", err.Err.Error())

}
