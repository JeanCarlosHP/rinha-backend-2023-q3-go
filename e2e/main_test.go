package e2e_test

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"
)

var (
	httpClient = &http.Client{}
	baseURL    = "http://localhost:8080"
)

func TestCreatePersonEndpoint(t *testing.T) {
	tests := []struct {
		name           string
		json           string
		expectedStatus int
	}{
		{
			name: "válido",
			json: `
				{
					"apelido" : "josé",
					"nome" : "José Roberto",
					"nascimento" : "2000-10-01",
					"stack" : ["C#", "Node", "Oracle"]
				}
			`,
			expectedStatus: 201,
		},
		{
			name: "válido com stack nula",
			json: `
				{
					"apelido" : "josé1",
					"nome" : "José Roberto",
					"nascimento" : "2000-10-01",
					"stack" : null
				}
			`,
			expectedStatus: 201,
		},
		{
			name: "apelido em uso",
			json: `
				{
					"apelido" : "josé",
					"nome" : "José Roberto",
					"nascimento" : "2000-10-01",
					"stack" : ["C#", "Node", "Oracle"]
				}
			`,
			expectedStatus: 422,
		},
		{
			name: "nome não pode ser nulo",
			json: `
				{
					"apelido" : "ana",
					"nome" : null,
					"nascimento" : "1985-09-23",
					"stack" : null
				}
			`,
			expectedStatus: 422,
		},
		{
			name: "apelido não pode ser nulo",
			json: `
				{
					"apelido" : null,
					"nome" : "Ana Barbosa",
					"nascimento" : "1985-09-23",
					"stack" : null
				}
			`,
			expectedStatus: 422,
		},
		{
			name: "nome deve ser string e não número",
			json: `
				{
					"apelido" : "apelido",
					"nome" : 1,
					"nascimento" : "1985-09-23",
					"stack" : null
				}
			`,
			expectedStatus: 400,
		},
		{
			name: "stack deve ser um array de apenas strings",
			json: `
				{
					"apelido" : "apelido",
					"nome" : 1,
					"nascimento" : "1985-09-23",
					"stack" : [1, "PHP"]
				}
			`,
			expectedStatus: 400,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := httpClient.Post(fmt.Sprintf("%s/pessoas", baseURL), "application/json", bytes.NewBuffer([]byte(tt.json)))
			if err != nil {
				t.Fatal(err)
			}

			if resp.StatusCode == tt.expectedStatus {
			} else {
				t.Errorf("Test %s failed: expected status %d, got %d", tt.name, tt.expectedStatus, resp.StatusCode)
			}
		})
	}
}

// func TestGetPersonByIdEndpoint(t *testing.T) {
// 	// Implement similar testing logic for GET /pessoas/[:id]
// }

// func TestSearchPeopleEndpoint(t *testing.T) {
// 	// Implement similar testing logic for GET /pessoas?t=[:termo da busca]
// }

// func TestCountPeopleEndpoint(t *testing.T) {
// 	// Implement similar testing logic for GET /contagem-pessoas
// }
