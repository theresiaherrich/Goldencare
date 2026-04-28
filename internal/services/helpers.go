package service

import (
	"fmt"
	"time"
)

func hitungStatusPemberian(jam string, diberikanPada *time.Time, now time.Time) string {
	if diberikanPada != nil {
		return "diberikan"
	}
	var hour, minute int
	if _, err := fmt.Sscanf(jam, "%d:%d", &hour, &minute); err != nil {
		return "menunggu"
	}
	jadwalTime := time.Date(now.Year(), now.Month(), now.Day(), hour, minute, 0, 0, now.Location())
	if now.After(jadwalTime.Add(30 * time.Minute)) {
		return "terlewat"
	}
	return "menunggu"
}

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