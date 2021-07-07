package main

import (
	"database/sql"
	"time"
)

type FisUser struct {
	IDZ      int    `gorm:"column:IDZ"`
	Prijmeni string `gorm:"column:Prijmeni"`
	Jmeno    string `gorm:"column:Jmeno"`
	Rfid     string `gorm:"column:Rfid"`
}

func (FisUser) TableName() string {
	return "fis_user"
}

type User struct {
	OID        int            `gorm:"column:OID"`
	Login      string         `gorm:"column:Login"`
	Password   sql.NullString `gorm:"column:Password"`
	Name       sql.NullString `gorm:"column:Name"`
	FirstName  sql.NullString `gorm:"column:FirstName"`
	Rfid       sql.NullString `gorm:"column:Rfid"`
	Barcode    sql.NullString `gorm:"column:Barcode"`
	Pin        sql.NullString `gorm:"column:Pin"`
	Function   sql.NullString `gorm:"column:Function"`
	UserTypeID sql.NullInt32  `gorm:"column:UserTypeID"`
	Email      sql.NullString `gorm:"column:Email"`
	Phone      sql.NullString `gorm:"column:Phone"`
	UserRoleID int            `gorm:"column:UserRoleID"`
}

func (User) TableName() string {
	return "user_new"
}

type FisProduct struct {
	IDVM     int    `gorm:"column:IDVM"`
	ArtNr    string `gorm:"column:ArtNr"`
	Nazev    string `gorm:"column:Nazev"`
	Velikost string `gorm:"column:Velikost"`
}

func (FisProduct) TableName() string {
	return "fis_product"
}

type Product struct {
	OID             int     `gorm:"column:OID"`
	Name            string  `gorm:"column:Name"`
	Barcode         string  `gorm:"column:Barcode"`
	Cycle           float64 `gorm:"column:Cycle"`
	ProductStatusID int     `gorm:"column:ProductStatusID"`
	Deleted         int     `gorm:"column:Deleted"`
}

func (Product) TableName() string {
	return "product_new"
}

type FisOrder struct {
	ID       int    `gorm:"column:ID"`
	IDVC     string `gorm:"column:IDVC"`
	IDVM     int    `gorm:"column:IDVM"`
	Mnozstvi int    `gorm:"column:Mnozstvi"`
}

func (FisOrder) TableName() string {
	return "fis_order"
}

type FisProduction struct {
	Id          int             `gorm:"column:Id"`
	IDFis       int             `gorm:"column:IDFis"`
	DatumCasOd  time.Time       `gorm:"column:DatumCasOd"`
	DatumCasDo  time.Time       `gorm:"column:DatumCasDo"`
	IDZ         sql.NullInt32   `gorm:"column:IDZ"`
	IDS         sql.NullInt32   `gorm:"column:IDS"`
	MnozstviOK  sql.NullInt32   `gorm:"column:MnozstviOK"`
	MnozstviNOK sql.NullInt32   `gorm:"column:MnozstviNOK"`
	KgOK        sql.NullInt32   `gorm:"column:KgOK"`
	KgNOK       sql.NullInt32   `gorm:"column:KgNOK"`
	ZapsiId     sql.NullInt32   `gorm:"column:ZapsiID"`
	IFS         sql.NullString  `gorm:"column:IFS"`
	Stav        sql.NullString  `gorm:"column:Stav"`
	Takt        sql.NullFloat64 `gorm:"column:Takt"`
	Prostoj     sql.NullString  `gorm:"column:Prostoj"`
	TypProstoje sql.NullString  `gorm:"column:TypProstoje"`
	Chyba       sql.NullString  `gorm:"column:Chyba"`
}

func (FisProduction) TableName() string {
	return "fis_production_new"
}

type Order struct {
	OID            int    `gorm:"column:OID"`
	Name           string `gorm:"column:Name"`
	Barcode        string `gorm:"column:Barcode"`
	ProductID      int    `gorm:"column:ProductID"`
	OrderStatusID  int    `gorm:"column:OrderStatusID"`
	CountRequested int    `gorm:"column:CountRequested"`
	Cavity         int    `gorm:"column:Cavity"`
}

func (Order) TableName() string {
	return "order_new"
}

type TerminalInputOrder struct {
	OID             int             `gorm:"column:OID"`
	DTS             time.Time       `gorm:"column:DTS"`
	DTE             sql.NullTime    `gorm:"column:DTE"`
	OrderID         int             `gorm:"column:OrderID"`
	UserID          sql.NullInt32   `gorm:"column:UserID"`
	DeviceID        int             `gorm:"column:DeviceID"`
	Interval        float64         `gorm:"column:Interval"`
	Count           sql.NullInt32   `gorm:"column:Count"`
	Fail            sql.NullInt32   `gorm:"column:Fail"`
	AverageCycle    float64         `gorm:"column:AverageCycle"`
	WorkerCount     int             `gorm:"column:WorkerCount"`
	WorkplaceModeID sql.NullInt32   `gorm:"column:WorkplaceModeID"`
	Note            sql.NullString  `gorm:"column:Note"`
	WorkshiftID     sql.NullInt32   `gorm:"column:WorkshiftID"`
	Cavity          int             `gorm:"column:Cavity"`
	ExtID           sql.NullInt32   `gorm:"column:ExtID"`
	ExtNum          sql.NullFloat64 `gorm:"column:ExtNum"`
	ExtText         sql.NullString  `gorm:"column:ExtTest"`
	ExtDT           sql.NullTime    `gorm:"column:ExtDT"`
}

func (TerminalInputOrder) TableName() string {
	return "terminal_input_order"
}

type TerminalInputIdle struct {
	OID      int             `gorm:"column:OID"`
	DTS      time.Time       `gorm:"column:DTS"`
	DTE      sql.NullTime    `gorm:"column:DTE"`
	IdleID   int             `gorm:"column:IdleID"`
	UserID   sql.NullInt32   `gorm:"column:UserID"`
	Interval float64         `gorm:"column:Interval"`
	DeviceID int             `gorm:"column:DeviceID"`
	Note     sql.NullString  `gorm:"column:Note"`
	ExtID    sql.NullInt32   `gorm:"column:ExtID"`
	ExtNum   sql.NullFloat64 `gorm:"column:ExtNum"`
	ExtText  sql.NullString  `gorm:"column:ExtText"`
	ExtDT    sql.NullTime    `gorm:"column:ExtDT"`
}

func (TerminalInputIdle) TableName() string {
	return "terminal_input_idle"
}

type Workplace struct {
	OID      int           `gorm:"column:OID"`
	Name     string        `gorm:"column:Name"`
	DeviceID sql.NullInt32 `gorm:"column:DeviceID"`
	Code     string        `gorm:"column:Code"`
}

func (Workplace) TableName() string {
	return "workplace"
}

type TerminalInputOrderIdle struct {
	TerminalInputOrderID int `gorm:"column:TerminalInputOrderID"`
	TerminalInputIdleID  int `gorm:"column:TerminalInputIdleID"`
}

func (TerminalInputOrderIdle) TableName() string {
	return "terminal_input_order_terminal_input_idle"
}

type Idle struct {
	OID        int    `gorm:"column:OID"`
	Name       string `gorm:"column:Name"`
	IdleTypeID int    `gorm:"column:IdleTypeID"`
}

func (Idle) TableName() string {
	return "idle"
}

type IdleType struct {
	OID  int    `gorm:"column:OID"`
	Name string `gorm:"column:Name"`
}

func (IdleType) TableName() string {
	return "idle_type"
}

type WorkplaceState struct {
	OID         int          `gorm:"column:OID"`
	DTS         time.Time    `gorm:"column:DTS"`
	DTE         sql.NullTime `gorm:"column:DTE"`
	StateID     int          `gorm:"column:StateID"`
	Interval    float64      `gorm:"column:Interval"`
	WorkplaceID int          `gorm:"column:WorkplaceID"`
}

func (WorkplaceState) TableName() string {
	return "workplace_state"
}
