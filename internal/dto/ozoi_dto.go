package dto

import "errors"

type CreateOzoiInput struct {
	Title       string `json:"title" binding:"required"`
	Completed   bool   `json:"completed"`
	Description string `json:"description,omitempty" db:"description"`
}

type UpdateOzoiInput struct {
	Title       string  `json:"title"`
	Completed   *bool   `json:"completed"`
	Description *string `json:"description,omitempty" db:"description"`
}

func (i *CreateOzoiInput) Validate() error {
	if len(i.Title) > 100 {
		return errors.New("title must be 100 characters or less")
	}
	if len(i.Description) > 500 {
		return errors.New("description must be 500 characters or less")
	}
	return nil
}

func (i *UpdateOzoiInput) Validate() error {
	if i.Title == "" && i.Completed == nil && i.Description == nil {
		return errors.New("at least one field must be provided")
	}
	if len(i.Title) > 100 {
		return errors.New("title must be 100 characters or less")
	}
	if i.Description != nil && len(*i.Description) > 500 {
		return errors.New("description must be 500 characters or less")
	}
	return nil
}
