package main

import (
	"gorm.io/gorm"
	"strconv"
	"time"
)

func exportStatePowerOffFromZapsi(db *gorm.DB) {
	start := time.Now()
	logInfo("MAIN", "Exporting poweroff states started")
	var workplaceStates []WorkplaceState
	db.Where("StateID = 3").Where("dte is not null and dte > ?, ", time.Now().Add(-24*time.Hour)).Find(&workplaceStates)
	var fisProductionOfflines []FisProduction
	db.Where("stav like 'v'").Find(&fisProductionOfflines)
	var cachedFisProductionOfflines map[int]FisProduction
	for _, fisProductionOffline := range fisProductionOfflines {
		cachedFisProductionOfflines[fisProductionOffline.ZapsiId] = fisProductionOffline
	}
	for _, workplaceState := range workplaceStates {
		_, stateExported := cachedFisProductionOfflines[workplaceState.OID]
		if !stateExported {
			logInfo("MAIN", "Exporting new state with OID "+strconv.Itoa(workplaceState.OID))
			var fisProduction FisProduction
			fisProduction.ZapsiId = workplaceState.OID
			fisProduction.DatumCasOd = workplaceState.DTS
			fisProduction.DatumCasDo = workplaceState.DTE.Time
			var workplace Workplace
			db.Where("DeviceId = ?", workplaceState.WorkplaceID).Find(&workplace)
			fisProduction.IFS = workplace.Code
			fisProduction.Stav = "v"
			// TODO: enable save
			//db.Save(&fisProduction)
		}
	}
	logInfo("MAIN", "Exporting poweroff states ended in "+time.Since(start).String())
}

func exportIdlesFromZapsi(db *gorm.DB) {
	start := time.Now()
	logInfo("MAIN", "Exporting idles started")
	var terminalInputIdles []TerminalInputIdle
	db.Where("dte is not null and dte > ?, ", time.Now().Add(-24*time.Hour)).Find(&terminalInputIdles)
	var fisProductionIdles []FisProduction
	db.Where("stav like 'p'").Find(&fisProductionIdles)
	var cachedFisProductionIdles map[int]FisProduction
	for _, fisProductionIdle := range fisProductionIdles {
		cachedFisProductionIdles[fisProductionIdle.ZapsiId] = fisProductionIdle
	}
	for _, terminalInputIdle := range terminalInputIdles {
		_, idleExported := cachedFisProductionIdles[terminalInputIdle.OID]
		if !idleExported {
			logInfo("MAIN", "Exporting new idle with OID "+strconv.Itoa(terminalInputIdle.OID))
			var fisProduction FisProduction
			fisProduction.ZapsiId = terminalInputIdle.OID
			fisProduction.DatumCasOd = terminalInputIdle.DTS
			fisProduction.DatumCasDo = terminalInputIdle.DTE.Time
			var zapsiUser User
			db.Where("OID = ?", terminalInputIdle.UserID).Find(&zapsiUser)
			if zapsiUser.OID > 1 {
				var fisUser FisUser
				db.Where("IDZ = ?", zapsiUser.Login).Find(&fisUser)
				if fisUser.IDZ > 0 {
					fisProduction.IDZ = zapsiUser.Login
				} else {
					logInfo("MAIN", "Order with ID "+strconv.Itoa(terminalInputIdle.OID)+" has user with login "+zapsiUser.Login+" that is not in fis table, error")
					fisProduction.IDZ = "0"
					fisProduction.Chyba = zapsiUser.Login
				}
			}
			var workplace Workplace
			db.Where("DeviceId = ?", terminalInputIdle.DeviceID).Find(&workplace)
			fisProduction.IFS = workplace.Code
			var terminalInputOrderIdle TerminalInputOrderIdle
			db.Where("TerminalInputIdleIdD = ?", terminalInputIdle.OID).Find(&terminalInputOrderIdle)
			var terminalInputOrder TerminalInputOrder
			db.Where("OID = ?", terminalInputOrderIdle.TerminalInputOrderID).Find(&terminalInputOrder)
			var zapsiOrder Order
			db.Where("OID = ?", terminalInputOrder.OrderID).Find(&zapsiOrder)
			if zapsiOrder.OID > 0 {
				var fisOrder FisOrder
				db.Where("ID = ?", zapsiOrder.Barcode).Find(&fisOrder)
				if fisOrder.ID > 0 {
					fisProduction.IDFis, _ = strconv.Atoi(zapsiOrder.Barcode)
				} else {
					logInfo("MAIN", "Order with ID "+strconv.Itoa(terminalInputOrder.OID)+" has order with barcode "+zapsiOrder.Barcode+" that is not in fis table, error")
					fisProduction.IDFis = 0
					fisProduction.Chyba = fisProduction.Chyba + "," + zapsiOrder.Barcode
				}
			}
			fisProduction.Stav = "p"
			var idle Idle
			db.Where("OID = ?", terminalInputIdle.IdleID).Find(&idle)
			var idleType IdleType
			db.Where("OID = ?", idle.IdleTypeID).Find(&idleType)
			fisProduction.Prostoj = idle.Name
			fisProduction.TypProstoje = idleType.Name
			// TODO: enable save
			//db.Save(&fisProduction)
		}
	}
	logInfo("MAIN", "Exporting idles ended in "+time.Since(start).String())
}

func exportOrdersFromZapsi(db *gorm.DB) {
	start := time.Now()
	logInfo("MAIN", "Exporting orders started")
	var terminalInputOrders []TerminalInputOrder
	db.Where("dte is not null and dte > ?, ", time.Now().Add(-24*time.Hour)).Find(&terminalInputOrders)
	var fisProductionOrders []FisProduction
	db.Where("stav like 'a'").Find(&fisProductionOrders)
	var cachedFisProductionOrders map[int]FisProduction
	for _, fisProductionOrder := range fisProductionOrders {
		cachedFisProductionOrders[fisProductionOrder.ZapsiId] = fisProductionOrder
	}
	for _, terminalInputOrder := range terminalInputOrders {
		_, orderExported := cachedFisProductionOrders[terminalInputOrder.OID]
		if !orderExported {
			logInfo("MAIN", "Exporting new order with OID "+strconv.Itoa(terminalInputOrder.OID))
			var fisProduction FisProduction
			fisProduction.ZapsiId = terminalInputOrder.OID
			fisProduction.DatumCasOd = terminalInputOrder.DTS
			fisProduction.DatumCasDo = terminalInputOrder.DTE.Time
			var zapsiUser User
			db.Where("OID = ?", terminalInputOrder.UserID).Find(&zapsiUser)
			if zapsiUser.OID > 1 {
				var fisUser FisUser
				db.Where("IDZ = ?", zapsiUser.Login).Find(&fisUser)
				if fisUser.IDZ > 0 {
					fisProduction.IDZ = zapsiUser.Login
				} else {
					logInfo("MAIN", "Order with ID "+strconv.Itoa(terminalInputOrder.OID)+" has user with login "+zapsiUser.Login+" that is not in fis table, error")
					fisProduction.IDZ = "0"
					fisProduction.Chyba = zapsiUser.Login
				}
			}
			var workplace Workplace
			db.Where("DeviceId = ?", terminalInputOrder.DeviceID).Find(&workplace)
			fisProduction.IFS = workplace.Code
			var zapsiOrder Order
			db.Where("OID = ?", terminalInputOrder.OrderID).Find(&zapsiOrder)
			if zapsiOrder.OID > 0 {
				var fisOrder FisOrder
				db.Where("ID = ?", zapsiOrder.Barcode).Find(&fisOrder)
				if fisOrder.ID > 0 {
					fisProduction.IDFis, _ = strconv.Atoi(zapsiOrder.Barcode)
				} else {
					logInfo("MAIN", "Order with ID "+strconv.Itoa(terminalInputOrder.OID)+" has order with barcode "+zapsiOrder.Barcode+" that is not in fis table, error")
					fisProduction.IDFis = 0
					fisProduction.Chyba = fisProduction.Chyba + "," + zapsiOrder.Barcode
				}
			}
			fisProduction.MnozstviOK = int(terminalInputOrder.Count.Int32 - terminalInputOrder.Fail.Int32)
			fisProduction.MnozstviNOK = int(terminalInputOrder.Fail.Int32)
			fisProduction.KgOK, _ = strconv.Atoi(terminalInputOrder.Note.String)
			fisProduction.Stav = "a"
			fisProduction.Takt = terminalInputOrder.AverageCycle
			// TODO: enable save
			//db.Save(&fisProduction)
		}
	}
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
	noOfUpdatedOrders := 0
	noOfInsertedOrders := 0
	for _, fisOrder := range fisOrders {
		order, orderInZapsi := ordersMap[strconv.Itoa(fisOrder.ID)]
		if orderInZapsi {
			logInfo("MAIN", strconv.Itoa(fisOrder.ID)+": updating order name to "+fisOrder.IDVC+" and count requested to "+strconv.Itoa(fisOrder.Mnozstvi))
			var updateOrder Order
			db.Where("OID = ?", order.OID).Find(&updateOrder)
			updateOrder.Name = fisOrder.IDVC
			updateOrder.CountRequested = fisOrder.Mnozstvi
			// TODO: enable save
			//db.Save(&updateOrder)
			noOfUpdatedOrders++
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
			// TODO: enable save
			//db.Save(&newOrder)
			noOfInsertedOrders++
		}
	}
	logInfo("MAIN", "Updated "+strconv.Itoa(noOfUpdatedOrders)+" orders, inserted "+strconv.Itoa(noOfInsertedOrders)+" orders")
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
	noOfUpdatedProducts := 0
	noOfInsertedProducts := 0
	for _, fisProduct := range fisProducts {
		product, productInZapsi := productsMap[fisProduct.ArtNr]
		if productInZapsi {
			logInfo("MAIN", fisProduct.ArtNr+": updating product name to "+fisProduct.Nazev)
			var updateProduct Product
			db.Where("OID = ?", product.OID).Find(&updateProduct)
			updateProduct.Name = fisProduct.Nazev
			// TODO: enable save
			//db.Save(&updateProduct)
			noOfUpdatedProducts++
		} else {
			logInfo("MAIN", fisProduct.ArtNr+": adding new product with name "+fisProduct.Nazev)
			var newProduct Product
			newProduct.Name = fisProduct.Nazev
			newProduct.Barcode = fisProduct.ArtNr
			newProduct.Cycle = 0.0
			newProduct.ProductStatusID = 1
			newProduct.Deleted = 0
			// TODO: enable save
			//db.Save(&newProduct)
			noOfInsertedProducts++
		}
	}
	logInfo("MAIN", "Updated "+strconv.Itoa(noOfUpdatedProducts)+" products, inserted "+strconv.Itoa(noOfInsertedProducts)+" products")
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
	noOfUpdatedUsers := 0
	noOfInsertedUsers := 0
	for _, fisUser := range fisUsers {
		user, userInZapsi := usersMap[strconv.Itoa(fisUser.IDZ)]
		if userInZapsi {
			logInfo("MAIN", "["+strconv.Itoa(fisUser.IDZ)+"] "+fisUser.Jmeno+" "+fisUser.Prijmeni+": updating user with rfid "+fisUser.Rfid)
			var updateUser User
			db.Where("OID = ?", user.OID).Find(&updateUser)
			updateUser.Rfid = fisUser.Rfid
			// TODO: enable save
			//db.Save(&updateUser)
			noOfUpdatedUsers++
		} else {
			logInfo("MAIN", "["+strconv.Itoa(fisUser.IDZ)+"] "+fisUser.Jmeno+" "+fisUser.Prijmeni+": adding new user with  rfid ["+fisUser.Rfid+"]")
			var newUser User
			newUser.Login = strconv.Itoa(fisUser.IDZ)
			newUser.Name = fisUser.Prijmeni
			newUser.FirstName = fisUser.Jmeno
			newUser.Rfid = fisUser.Rfid
			newUser.UserRoleId = 2
			// TODO: enable save
			//db.Save(&newUser)
			noOfInsertedUsers++
		}
	}
	logInfo("MAIN", "Updated "+strconv.Itoa(noOfUpdatedUsers)+" users, inserted "+strconv.Itoa(noOfInsertedUsers)+" users")
	logInfo("MAIN", "Importing users ended in "+time.Since(start).String())
}
