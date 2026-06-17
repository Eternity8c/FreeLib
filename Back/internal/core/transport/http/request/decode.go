package core_http_request

import (
	"encoding/json"
	"fmt"
	"net/http"

	core_errors "github.com/Eternity8c/FreeLib/internal/core/errors"

	"github.com/go-playground/form/v4"
	"github.com/go-playground/validator/v10"
)

var (
	requestValidator = validator.New()
	formDecode       = form.NewDecoder()
)

func DecodeAndValidateJSONRequest(r *http.Request, dest any) error {
	if err := json.NewDecoder(r.Body).Decode(dest); err != nil {
		return fmt.Errorf(
			"decode json: %v: %w",
			err,
			core_errors.ErrInvalidArgumment,
		)
	}

	if err := requestValidator.Struct(dest); err != nil {
		return fmt.Errorf(
			"request validator: %v: %w",
			err,
			core_errors.ErrInvalidArgumment,
		)
	}

	return nil
}

func DecodeAndValidateFormData(r *http.Request, dest any) error {
	if err := r.ParseMultipartForm(32 * 1024 * 1024); err != nil {
		return fmt.Errorf("parse form: %w", err)
	}

	if err := formDecode.Decode(dest, r.Form); err != nil {
		return fmt.Errorf("decode form: %w", err)
	}

	if err := requestValidator.Struct(dest); err != nil {
		return fmt.Errorf(
			"request validator: %v: %w",
			err,
			core_errors.ErrInvalidArgumment,
		)
	}

	return nil
}
