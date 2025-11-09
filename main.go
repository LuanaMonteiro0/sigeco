package main

import (
	"time"
)

type Person struct {
	ID    string 
	Name  string
	Phone string
}

type RegistryEntry struct {
	PersonID     string
	TimestampIn  time.Time
	TimestampOut time.Time
}
