package main

import "errors"

var (
	ErrSeedNotFound = errors.New("database not seeded")
	ErrTimedOut     = errors.New("connexion timed out before finding the page/element")
)
