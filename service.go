package main

import (
	"gorm.io/gorm"
	"time"
)

func exportStatePowerOffFromZapsi(db *gorm.DB) {
	start := time.Now()
	logInfo("MAIN", "Exporting poweroff states started")

	logInfo("MAIN", "Exporting poweroff states ended in "+time.Since(start).String())
}

func exportIdlesFromZapsi(db *gorm.DB) {
	start := time.Now()
	logInfo("MAIN", "Exporting idles started")

	logInfo("MAIN", "Exporting idles ended in "+time.Since(start).String())
}

func exportOrdersFromZapsi(db *gorm.DB) {
	start := time.Now()
	logInfo("MAIN", "Exporting orders started")

	logInfo("MAIN", "Exporting orders ended in "+time.Since(start).String())
}

func importOrdersToZapsi(db *gorm.DB) {
	start := time.Now()
	logInfo("MAIN", "Importing orders started")

	logInfo("MAIN", "Importing orders ended in "+time.Since(start).String())
}

func importProductsToZapsi(db *gorm.DB) {
	start := time.Now()
	logInfo("MAIN", "Importing products started")

	logInfo("MAIN", "Importing products ended in "+time.Since(start).String())
}

func importUsersToZapsi(db *gorm.DB) {
	start := time.Now()
	logInfo("MAIN", "Importing users started")

	logInfo("MAIN", "Importing users ended in "+time.Since(start).String())
}
