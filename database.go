package main

import (
	"database/sql"
	"time"
)

type FisUser struct {
	IDZ      int
	Prijmeni string
	Jmeno    string
	Rfid     string
}

func (FisUser) TableName() string {
	return "fis_user"
}

type User struct {
	OID        int
	Login      string
	Name       string
	FirstName  string
	Rfid       string
	UserRoleId int
}

func (User) TableName() string {
	return "user"
}

type FisProduct struct {
	IDVM     int
	ArtNr    string
	Nazev    string
	Velikost string
}

func (FisProduct) TableName() string {
	return "fis_product"
}

type Product struct {
	OID             int
	Name            string
	Barcode         string
	Cycle           float64
	ProductStatusID int
	Deleted         int
}

func (Product) TableName() string {
	return "product"
}

type FisOrder struct {
	ID       int
	IDVC     string
	IDVM     int
	Mnozstvi int
}

func (FisOrder) TableName() string {
	return "fis_order"
}

type FisProduction struct {
	Id          int
	IDFis       int
	DatumCasOd  time.Time
	DatumCasDo  time.Time
	IDZ         string
	IDS         int
	MnozstviOK  int
	MnozstviNOK int
	KgOK        int
	KgNOK       int
	Prenos      int
	ZapsiId     int
	IFS         string
	Stav        string
	Takt        float64
	Prostoj     string
	TypProstoje string
	Chyba       string
}

func (FisProduction) TableName() string {
	return "fis_production"
}

type Order struct {
	OID            int
	Name           string
	Barcode        string
	ProductID      int
	OrderStatusId  int
	CountRequested int
	Cavity         int
}

func (Order) TableName() string {
	return "order"
}

type TerminalInputOrder struct {
	OID             int
	DTS             time.Time
	DTE             sql.NullTime
	OrderID         int
	UserID          sql.NullInt32
	DeviceID        int
	Interval        float64
	Count           sql.NullInt32
	Fail            sql.NullInt32
	AverageCycle    float64
	WorkerCount     int
	WorkplaceModeID sql.NullInt32
	Note            sql.NullString
	WorkshiftID     sql.NullInt32
	Cavity          int
	ExtID           sql.NullInt32
	ExtNum          sql.NullFloat64
	ExtText         sql.NullString
	ExtDT           sql.NullTime
}

func (Order) TerminalInputOrder() string {
	return "terminal_input_order"
}

type TerminalInputIdle struct {
	OID      int
	DTS      time.Time
	DTE      sql.NullTime
	IdleID   int
	UserID   sql.NullInt32
	Interval float64
	DeviceID int
	Note     sql.NullString
	ExtID    sql.NullInt32
	ExtNum   sql.NullFloat64
	ExtText  sql.NullString
	ExtDT    sql.NullTime
}

func (Order) TerminalInputIdle() string {
	return "terminal_input_idle"
}

type Workplace struct {
	OID      int
	Name     string
	DeviceID sql.NullInt32
	Code     string
}

func (Order) Workplace() string {
	return "workplace"
}

type TerminalInputOrderIdle struct {
	TerminalInputOrderID int
	TerminalInputIdleID  int
}

func (Order) TerminalInputOrderIdle() string {
	return "terminal_input_order_terminal_input_idle"
}

type Idle struct {
	OID        int
	Name       string
	IdleTypeID int
}

func (Order) Idle() string {
	return "idle"
}

type IdleType struct {
	OID  int
	Name string
}

func (Order) IdleType() string {
	return "idle_type"
}

type WorkplaceState struct {
	OID         int
	DTS         time.Time
	DTE         sql.NullTime
	StateID     int
	Interval    float64
	WorkplaceID int
}

func (Order) WorkplaceState() string {
	return "workplace_state"
}
