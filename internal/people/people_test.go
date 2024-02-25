package people_test

import (
	"testing"

	"github.com/jeancarloshp/rinha-backend-go/internal/people"
)

func TestPeopleDTO_Validate(t *testing.T) {
	tests := []struct {
		name    string
		p       *people.PeopleDTO
		wantErr bool
	}{
		{
			name: "Valid",
			p: &people.PeopleDTO{
				Nickname: "nickname",
				Name:     "name",
				Birth:    "2006-01-02",
				Stack:    []string{"stack"},
			},
			wantErr: false,
		},
		{
			name: "Nickname is empty",
			p: &people.PeopleDTO{
				Nickname: "",
				Name:     "name",
				Birth:    "2006-01-02",
				Stack:    []string{"stack"},
			},
			wantErr: true,
		},
		{
			name: "Nickname is too long",
			p: &people.PeopleDTO{
				Nickname: "nickname too long nickname too long nickname too long nickname too long",
				Name:     "",
				Birth:    "2006-01-02",
				Stack:    []string{"stack"},
			},
			wantErr: true,
		},
		{
			name: "Name is empty",
			p: &people.PeopleDTO{
				Nickname: "nickname",
				Name:     "",
				Birth:    "2006-01-02",
				Stack:    []string{"stack"},
			},
			wantErr: true,
		},
		{
			name: "Name is too long",
			p: &people.PeopleDTO{
				Nickname: "nickname",
				Name:     "name is too long name is too long name is too long name is too long name is too long name is too long",
				Birth:    "2006-01-02",
				Stack:    []string{"stack"},
			},
			wantErr: true,
		},
		{
			name: "Invalid Birth",
			p: &people.PeopleDTO{
				Nickname: "nickname",
				Name:     "name",
				Birth:    "15-02-2005",
				Stack:    []string{"stack"},
			},
			wantErr: true,
		},
		{
			name: "Invalid Stack",
			p: &people.PeopleDTO{
				Nickname: "nickname",
				Name:     "name",
				Birth:    "2006-01-02",
				Stack:    []string{"stack", "stack is too long stack is too long"},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.p.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("PeopleDTO.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
