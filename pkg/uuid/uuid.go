package uuid

import "github.com/google/uuid"

type UUID = uuid.UUID

func New() UUID {
	return uuid.New()
}

func IsValidUUIDv7(id string) bool {
	_, err := uuid.Parse(id)
	return err == nil
}

func MustParse(id string) UUID {
	return uuid.MustParse(id)
}

func Nil() UUID {
	return uuid.Nil
}
