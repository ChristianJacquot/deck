package solver

import (
	"github.com/kong/deck/crud"
	"github.com/kong/deck/diff"
	"github.com/kong/deck/state"
	"github.com/kong/go-kong/kong"
)

// developerCRUD implements crud.Actions interface.
type developerCRUD struct {
	client *kong.Client
}

func developerFromStuct(arg diff.Event) *state.Developer {
	developer, ok := arg.Obj.(*state.Developer)
	if !ok {
		panic("unexpected type, expected *state.developer")
	}
	return developer
}

// Create creates a Developer in Kong.
// The arg should be of type diff.Event, containing the developer to be created,
// else the function will panic.
// It returns a the created *state.Developer.
func (s *developerCRUD) Create(arg ...crud.Arg) (crud.Arg, error) {
	event := eventFromArg(arg[0])
	developer := developerFromStuct(event)
	createdDeveloper, err := s.client.Developers.Create(nil, &developer.Developer)
	if err != nil {
		return nil, err
	}
	return &state.Developer{Developer: *createdDeveloper}, nil
}

// Delete deletes a Developer in Kong.
// The arg should be of type diff.Event, containing the developer to be deleted,
// else the function will panic.
// It returns a the deleted *state.Developer.
func (s *developerCRUD) Delete(arg ...crud.Arg) (crud.Arg, error) {
	event := eventFromArg(arg[0])
	developer := developerFromStuct(event)
	err := s.client.Developers.Delete(nil, developer.ID)
	if err != nil {
		return nil, err
	}
	return developer, nil
}

// Update updates a Developer in Kong.
// The arg should be of type diff.Event, containing the developer to be updated,
// else the function will panic.
// It returns a the updated *state.Developer.
func (s *developerCRUD) Update(arg ...crud.Arg) (crud.Arg, error) {
	event := eventFromArg(arg[0])
	developer := developerFromStuct(event)

	updatedDeveloper, err := s.client.Developers.Create(nil, &developer.Developer)
	if err != nil {
		return nil, err
	}
	return &state.Developer{Developer: *updatedDeveloper}, nil
}
