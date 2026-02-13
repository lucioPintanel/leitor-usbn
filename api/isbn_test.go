package api

import (
	"testing"
)

func TestValidateISBN(t *testing.T) {
	tests := []struct {
		name    string
		isbn    string
		wantErr bool
	}{
		// ISBN-10 válidos (Clean Code by Robert C. Martin)
		{
			name:    "ISBN-10 válido sem formatação",
			isbn:    "0132350882",
			wantErr: false,
		},
		{
			name:    "ISBN-10 válido com hífens",
			isbn:    "0-13-235088-2",
			wantErr: false,
		},

		// Inválidos
		{
			name:    "ISBN muito curto",
			isbn:    "123456789",
			wantErr: true,
		},
		{
			name:    "ISBN muito longo",
			isbn:    "12345678901234",
			wantErr: true,
		},
		{
			name:    "ISBN com letras inválidas",
			isbn:    "978-ABC-235088-2",
			wantErr: true,
		},
		{
			name:    "ISBN-10 com checksum errado",
			isbn:    "0132350881",
			wantErr: true,
		},
		{
			name:    "ISBN vazio",
			isbn:    "",
			wantErr: true,
		},
		{
			name:    "ISBN com espaços e símbolos",
			isbn:    "0132@50882",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateISBN(tt.isbn)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateISBN(%q) error = %v, wantErr %v", tt.isbn, err, tt.wantErr)
			}
		})
	}
}

func TestNormalizeISBN(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "978-0-13-235088-2",
			expected: "9780132350882",
		},
		{
			input:    "0-13-235088-2",
			expected: "0132350882",
		},
		{
			input:    "978 0 13 235088 2",
			expected: "9780132350882",
		},
		{
			input:    "9780132350882",
			expected: "9780132350882",
		},
	}

	for _, tt := range tests {
		result := NormalizeISBN(tt.input)
		if result != tt.expected {
			t.Errorf("NormalizeISBN(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}
