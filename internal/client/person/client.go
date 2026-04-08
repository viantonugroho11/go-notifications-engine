package person

import (
	"context"
)

type PersonClient interface {
	GetPerson(ctx context.Context, id string) (Person, error)
}