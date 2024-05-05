package entity

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

// https://mholt.github.io/json-to-go/ para converte json para stuc
type ViaCEP struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

type LocalidadeViaCep struct{}

func NewLocalidadeViaCep() *LocalidadeViaCep {
	return &LocalidadeViaCep{}
}

func (c *LocalidadeViaCep) ColetaLocalidade(cep string) (*ViaCEP, *RequestError) {

	req, err := http.Get("http://www.viacep.com.br/ws/" + cep + "/json/")
	if err != nil {
		return nil, &RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        errors.New("www.viacep.com inacessivel"),
		}
	}
	defer req.Body.Close()
	res, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, &RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        errors.New("www.viacep.com inacessivel"),
		}
	}

	var data ViaCEP
	err = json.Unmarshal(res, &data)
	if err != nil {
		return nil, &RequestError{
			StatusCode: http.StatusUnprocessableEntity,
			Err:        errors.New("invalid zipcode"),
		}

	}

	if len(data.Localidade) == 0 {
		return nil, &RequestError{
			StatusCode: http.StatusNotFound,
			Err:        errors.New("can not find zipcode"),
		}
	}

	return &data, nil
}
