package vo

import (
	"github.com/eskokado/startup-auth-go/backend/pkg/msgerror"
	"github.com/google/uuid"
)

type ID struct {
	value uuid.UUID
}

func NewID() ID {
	return ID{value: uuid.New()}
}

func ParseID(s string) (ID, error) {
	if s == "" {
		return ID{}, msgerror.AnErrEmptyID
	}

	parsed, err := uuid.Parse(s)
	if err != nil {
		return ID{}, msgerror.AnErrInvalidID
	}

	return ID{value: parsed}, nil
}

func (i ID) String() string {
	return i.value.String()
}

func (i ID) Equal(other ID) bool {
	return i.value == other.value
}
