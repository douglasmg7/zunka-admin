package main

import (
// "net/http"
// "net/http/httptest"
// "testing"
)

/******************************************************************************
* Address.
*******************************************************************************/
// // Address by CEP.
// func TestGetAddressByCEPAPI(t *testing.T) {
// reqBody := strings.NewReader("31170-210")
// url := "/freightsrv/address"
// req, _ := http.NewRequest(http.MethodGet, url, reqBody)

// req.SetBasicAuth("bypass", "123456")
// req.Header.Set("Content-Type", "application/json")

// res := httptest.NewRecorder()

// router.ServeHTTP(res, req)
// if res.Code != 200 {
// t.Errorf("Returned code: %d", res.Code)
// return
// }

// address := viaCEPAddress{}

// err = json.Unmarshal(res.Body.Bytes(), &address)
// if err != nil {
// t.Errorf("Err: %s", err)
// return
// }

// want := viaCEPAddress{
// Cep:      "31170-210",
// Street:   "Rua Deputado Bernardino de Sena Figueiredo",
// District: "Cidade Nova",
// City:     "Belo Horizonte",
// State:    "MG",
// }

// if address.Cep != want.Cep || address.Street != want.Street || address.District != want.District || address.City != want.City || address.State != want.State {
// t.Errorf("got:  %+v\nwant %+v", address, want)
// }
// }
