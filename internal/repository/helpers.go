package repository

import (
	"fmt"
	"time"
)

func hariIndonesia(d time.Weekday) string {
	return map[time.Weekday]string{
		time.Sunday:    "Minggu",
		time.Monday:    "Senin",
		time.Tuesday:   "Selasa",
		time.Wednesday: "Rabu",
		time.Thursday:  "Kamis",
		time.Friday:    "Jumat",
		time.Saturday:  "Sabtu",
	}[d]
}

func errNotFound(entity string) error {
	return fmt.Errorf("%s not found", entity)
}