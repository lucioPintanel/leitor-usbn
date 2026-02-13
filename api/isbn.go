package api

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

// ValidateISBN verifica se a string é um ISBN válido (10 ou 13 dígitos)
func ValidateISBN(isbn string) error {
	// Remove hífens e espaços
	isbn = strings.ReplaceAll(isbn, "-", "")
	isbn = strings.ReplaceAll(isbn, " ", "")

	// Verifica se contém apenas dígitos (e pode ter X no final para ISBN-10)
	if len(isbn) == 10 {
		if !regexp.MustCompile(`^\d{9}[\dX]$`).MatchString(isbn) {
			return errors.New("ISBN-10 deve conter apenas dígitos (0-9) e pode ter X no final")
		}
	} else {
		if !regexp.MustCompile(`^\d+$`).MatchString(isbn) {
			return errors.New("ISBN deve conter apenas dígitos")
		}
	}

	// ISBN-10 ou ISBN-13
	switch len(isbn) {
	case 10:
		return validateISBN10(isbn)
	case 13:
		return validateISBN13(isbn)
	default:
		return errors.New("ISBN deve ter 10 ou 13 dígitos")
	}
}

// validateISBN10 calcula e valida o dígito verificador de ISBN-10
func validateISBN10(isbn string) error {
	if len(isbn) != 10 {
		return errors.New("ISBN-10 deve ter exatamente 10 dígitos")
	}

	sum := 0
	for i := 0; i < 9; i++ {
		digit, _ := strconv.Atoi(string(isbn[i]))
		sum += digit * (10 - i)
	}

	checksum := (11 - (sum % 11)) % 11
	expectedCheck := string(isbn[9])

	if checksum == 10 {
		if expectedCheck != "X" {
			return errors.New("dígito verificador ISBN-10 inválido")
		}
	} else {
		if expectedCheck != string(rune(checksum+'0')) {
			return errors.New("dígito verificador ISBN-10 inválido")
		}
	}

	return nil
}

// validateISBN13 calcula e valida o dígito verificador de ISBN-13
func validateISBN13(isbn string) error {
	if len(isbn) != 13 {
		return errors.New("ISBN-13 deve ter exatamente 13 dígitos")
	}

	sum := 0
	for i := 0; i < 12; i++ {
		digit, _ := strconv.Atoi(string(isbn[i]))
		if i%2 == 0 {
			sum += digit
		} else {
			sum += digit * 3
		}
	}

	checksum := (10 - (sum % 10)) % 10
	expectedCheck, _ := strconv.Atoi(string(isbn[12]))

	if checksum != expectedCheck {
		return errors.New("dígito verificador ISBN-13 inválido")
	}

	return nil
}

// NormalizeISBN remove hífens e espaços de um ISBN
func NormalizeISBN(isbn string) string {
	isbn = strings.ReplaceAll(isbn, "-", "")
	isbn = strings.ReplaceAll(isbn, " ", "")
	return isbn
}
