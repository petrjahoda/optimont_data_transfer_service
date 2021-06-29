package main

import (
	"gorm.io/gorm"
	"strconv"
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
	var fisOrders []FisOrder
	db.Find(&fisOrders)
	logInfo("MAIN", "Downloaded "+strconv.Itoa(len(fisOrders))+" fis orders")
	var orders []Order
	db.Find(&orders)
	logInfo("MAIN", "Downloaded "+strconv.Itoa(len(orders))+" zapsi orders")
	var ordersMap map[string]Order
	for _, order := range orders {
		ordersMap[order.Barcode] = order
	}
	for _, fisOrder := range fisOrders {
		order, orderInZapsi := ordersMap[strconv.Itoa(fisOrder.ID)]
		if orderInZapsi {
			logInfo("MAIN", strconv.Itoa(fisOrder.ID)+": updating order name to "+fisOrder.IDVC+" and count requested to "+strconv.Itoa(fisOrder.Mnozstvi))
			var updateOrder Order
			db.Where("OID = ?", order.OID).Find(&updateOrder)
			updateOrder.Name = fisOrder.IDVC
			updateOrder.CountRequested = fisOrder.Mnozstvi
			db.Save(&updateOrder)
		} else {
			logInfo("MAIN", strconv.Itoa(fisOrder.ID)+": adding new order with name "+fisOrder.IDVC+" and count requested "+strconv.Itoa(fisOrder.Mnozstvi))
			var fisProduct FisProduct
			db.Where("IDVM = ?", fisOrder.IDVM).Find(&fisProduct)
			logInfo("MAIN", strconv.Itoa(fisOrder.ID)+": found fis product with ArtNr "+fisProduct.ArtNr)
			var product Product
			db.Where("barcode = ?", fisProduct.ArtNr).Find(&product)
			logInfo("MAIN", strconv.Itoa(fisOrder.ID)+": found zapsi product with Barcode "+product.Barcode)
			var newOrder Order
			newOrder.Name = fisOrder.IDVC
			newOrder.Barcode = strconv.Itoa(fisOrder.ID)
			newOrder.ProductID = product.OID
			newOrder.OrderStatusId = 1
			newOrder.CountRequested = fisOrder.Mnozstvi
			newOrder.Cavity = 1
			db.Save(&newOrder)
		}
	}
	logInfo("MAIN", "Importing orders ended in "+time.Since(start).String())
}

func importProductsToZapsi(db *gorm.DB) {
	start := time.Now()
	logInfo("MAIN", "Importing products started")
	var fisProducts []FisProduct
	db.Find(&fisProducts)
	logInfo("MAIN", "Downloaded "+strconv.Itoa(len(fisProducts))+" fis products")
	var products []Product
	db.Find(&products)
	logInfo("MAIN", "Downloaded "+strconv.Itoa(len(products))+" zapsi products")
	var productsMap map[string]Product
	for _, product := range products {
		productsMap[product.Barcode] = product
	}
	for _, fisProduct := range fisProducts {
		product, productInZapsi := productsMap[fisProduct.ArtNr]
		if productInZapsi {
			logInfo("MAIN", fisProduct.ArtNr+": updating product name to "+fisProduct.Nazev)
			var updateProduct Product
			db.Where("OID = ?", product.OID).Find(&updateProduct)
			updateProduct.Name = fisProduct.Nazev
			db.Save(&updateProduct)
		} else {
			logInfo("MAIN", fisProduct.ArtNr+": adding new product with name "+fisProduct.Nazev)
			var newProduct Product
			newProduct.Name = fisProduct.Nazev
			newProduct.Barcode = fisProduct.ArtNr
			newProduct.Cycle = 0.0
			newProduct.ProductStatusID = 1
			newProduct.Deleted = 0
			db.Save(&newProduct)
		}
	}
	logInfo("MAIN", "Importing products ended in "+time.Since(start).String())
}

func importUsersToZapsi(db *gorm.DB) {
	start := time.Now()
	logInfo("MAIN", "Importing users started")
	var fisUsers []FisUser
	db.Find(&fisUsers)
	logInfo("MAIN", "Downloaded "+strconv.Itoa(len(fisUsers))+" fis users")
	var users []User
	db.Find(&users)
	logInfo("MAIN", "Downloaded "+strconv.Itoa(len(users))+" zapsi users")
	var usersMap map[string]User
	for _, user := range users {
		usersMap[user.Login] = user
	}
	for _, fisUser := range fisUsers {
		user, userInZapsi := usersMap[strconv.Itoa(fisUser.IDZ)]
		if userInZapsi {
			logInfo("MAIN", "["+strconv.Itoa(fisUser.IDZ)+"] "+fisUser.Jmeno+" "+fisUser.Prijmeni+": updating user with rfid "+fisUser.Rfid)
			var updateUser User
			db.Where("OID = ?", user.OID).Find(&updateUser)
			updateUser.Rfid = fisUser.Rfid
			db.Save(&updateUser)
		} else {
			logInfo("MAIN", "["+strconv.Itoa(fisUser.IDZ)+"] "+fisUser.Jmeno+" "+fisUser.Prijmeni+": adding new user with  rfid ["+fisUser.Rfid+"]")
			var newUser User
			newUser.Login = strconv.Itoa(fisUser.IDZ)
			newUser.Name = fisUser.Prijmeni
			newUser.FirstName = fisUser.Jmeno
			newUser.Rfid = fisUser.Rfid
			newUser.UserRoleId = 2
			db.Save(&newUser)
		}
	}
	logInfo("MAIN", "Importing users ended in "+time.Since(start).String())
}
