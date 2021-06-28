package main

import "time"

type TerminalInputFail struct {
	OID      int       `gorm:"column:OID"`
	DT       time.Time `gorm:"column:DT"`
	FailID   int       `gorm:"column:FailID"`
	UserID   int       `gorm:"column:UserID"`
	DeviceID int       `gorm:"column:DeviceID"`
	Note     string    `gorm:"column:Note"`
}

func (TerminalInputFail) TableName() string {
	return "terminal_input_fail"
}

type Fail struct {
	OID        int    `gorm:"column:OID"`
	Name       string `gorm:"column:Name"`
	Barcode    string `gorm:"column:Barcode"`
	FailTypeID int    `gorm:"column:FailTypeID"`
}

func (Fail) TableName() string {
	return "fail"
}
