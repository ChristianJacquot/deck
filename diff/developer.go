package diff

import (
	"github.com/kong/deck/crud"
	"github.com/kong/deck/state"
	"github.com/pkg/errors"
)

func (sc *Syncer) deleteDevelopers() error {
	currentDevelopers, err := sc.currentState.Developers.GetAll()
	if err != nil {
		return errors.Wrap(err, "error fetching developers from state")
	}

	for _, developer := range currentDevelopers {
		n, err := sc.deleteDeveloper(developer)
		if err != nil {
			return err
		}
		if n != nil {
			err = sc.queueEvent(*n)
			if err != nil {
				return err
			}
		}

	}
	return nil
}

func (sc *Syncer) deleteDeveloper(developer *state.Developer) (*Event, error) {
	_, err := sc.targetState.Developers.Get(*developer.ID)
	if err == state.ErrNotFound {
		return &Event{
			Op:   crud.Delete,
			Kind: "developer",
			Obj:  developer,
		}, nil
	}
	if err != nil {
		return nil, errors.Wrapf(err, "looking up developer '%v'",
			developer.Identifier())
	}
	return nil, nil
}

func (sc *Syncer) createUpdateDevelopers() error {
	targetDevelopers, err := sc.targetState.Developers.GetAll()
	if err != nil {
		return errors.Wrap(err, "error fetching developers from state")
	}

	for _, developer := range targetDevelopers {
		n, err := sc.createUpdateDeveloper(developer)
		if err != nil {
			return err
		}
		if n != nil {
			err = sc.queueEvent(*n)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (sc *Syncer) createUpdateDeveloper(developer *state.Developer) (*Event, error) {
	developerCopy := &state.Developer{Developer: *developer.DeepCopy()}
	currentDeveloper, err := sc.currentState.Developers.Get(*developer.ID)

	if err == state.ErrNotFound {
		// developer not present, create it
		return &Event{
			Op:   crud.Create,
			Kind: "developer",
			Obj:  developerCopy,
		}, nil
	}
	if err != nil {
		return nil, errors.Wrapf(err, "error looking up developer %v",
			developer.Identifier())
	}

	// found, check if update needed
	if !currentDeveloper.EqualWithOpts(developerCopy, false, true) {
		return &Event{
			Op:     crud.Update,
			Kind:   "developer",
			Obj:    developerCopy,
			OldObj: currentDeveloper,
		}, nil
	}
	return nil, nil
}
