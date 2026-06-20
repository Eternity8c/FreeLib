package domain_test

import (
	"testing"

	"github.com/Eternity8c/FreeLib/internal/core/domain"
	"github.com/stretchr/testify/assert"
)

func TestCreateBookRequest_Validate(t *testing.T) {
	cases := []struct {
		name        string
		book        domain.Book
		wantErr     bool
		errContains string
	}{
		{
			name: "valid request",
			book: domain.Book{
				Title:  "Анна Каренина",
				Author: "Лев Толстой",
				Genre:  "Классика",
			},
			wantErr: false,
		},
		{
			name: "empty title",
			book: domain.Book{
				Title:  "",
				Author: "Лев Толстой",
				Genre:  "Классика",
			},
			wantErr:     true,
			errContains: "len `title`: 0",
		},
		{
			name: "empty author",
			book: domain.Book{
				Title:  "Анна Каренина",
				Author: "",
				Genre:  "Классика",
			},
			wantErr:     true,
			errContains: "len `author`: 0",
		},
		{
			name: "empty genre",
			book: domain.Book{
				Title:  "Анна Каренина",
				Author: "Лев Толстой",
				Genre:  "",
			},
			wantErr:     true,
			errContains: "len `genre`: 0",
		},
		{
			name: "len `author` 1",
			book: domain.Book{
				Title:  "Анна Каренина",
				Author: "Л",
				Genre:  "Классика",
			},
			wantErr:     true,
			errContains: "len `author`: 1",
		},
		{
			name: "len `author` 2",
			book: domain.Book{
				Title:  "Анна Каренина",
				Author: "Ле",
				Genre:  "Классика",
			},
			wantErr:     true,
			errContains: "len `author`: 2",
		},
		{
			name: "first letter author not upper",
			book: domain.Book{
				Title:  "Анна Каренина",
				Author: "лев Толстой",
				Genre:  "Классика",
			},
			wantErr:     true,
			errContains: "first letter author not upper: л",
		},
		{
			name: "first letter genre not upper",
			book: domain.Book{
				Title:  "Анна Каренина",
				Author: "Лев Толстой",
				Genre:  "классика",
			},
			wantErr:     true,
			errContains: "first letter genre not upper: к",
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.book.Validate()

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContains)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
