package main

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
