package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
	"time"
)

type BrasilApi struct {
	Cep          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
	Service      string `json:"service"`
}

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

type CepHandler struct {
}

func NewCepHandler() *CepHandler {
	return &CepHandler{}
}

func (h *CepHandler) GetAddress(w http.ResponseWriter, r *http.Request) {
	c1 := make(chan *BrasilApi)
	c2 := make(chan *ViaCEP)

	cepParam := chi.URLParam(r, "cep")
	if cepParam == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// brasilapi
	go func() {
		address, err := getAddressFromBrasilApi(cepParam)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		c1 <- address
	}()

	// viacep
	go func() {
		address, err := getAddressFromViaCep(cepParam)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		c2 <- address
	}()

	select {
	case address := <-c1: // brasilapi
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(*address)
		fmt.Printf("Received from BrasilApi the complete address: %s\n", prettyPrint(address))

	case address := <-c2: // viacep
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(*address)
		fmt.Printf("Received from ViaCEP the complete address: %s\n", prettyPrint(address))

	case <-time.After(time.Second * 1):
		println("request timeout")
		w.WriteHeader(http.StatusRequestTimeout)
	}
}

func getAddressFromBrasilApi(cep string) (*BrasilApi, error) {
	resp, err := http.Get("https://brasilapi.com.br/api/cep/v1/" + cep)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var c BrasilApi
	err = json.Unmarshal(body, &c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func getAddressFromViaCep(cep string) (*ViaCEP, error) {
	resp, err := http.Get("https://viacep.com.br/ws/" + cep + "/json/")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var c ViaCEP
	err = json.Unmarshal(body, &c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func prettyPrint(v interface{}) string {
	b, _ := json.MarshalIndent(v, "", "  ")
	return string(b)
}
