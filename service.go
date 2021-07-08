package main

import (
	"database/sql"
	"gorm.io/gorm"
	"strconv"
	"time"
)

func exportStatePowerOffFromZapsi(db *gorm.DB) {
	start := time.Now()
	logInfo("MAIN", "Exporting poweroff states started")
	var workplaceStates []WorkplaceState
	db.Where("StateID = 3").Where("dte is not null and dte > ?", time.Now().Add(-24*time.Hour)).Find(&workplaceStates)
	var fisProductionOfflines []FisProduction
	db.Where("stav like 'v'").Find(&fisProductionOfflines)
	cachedFisProductionOfflines := make(map[int]FisProduction)
	for _, fisProductionOffline := range fisProductionOfflines {
		cachedFisProductionOfflines[int(fisProductionOffline.ZapsiId.Int32)] = fisProductionOffline
	}
	noOfExportedPoweroffs := 0
	for _, workplaceState := range workplaceStates {
		_, stateExported := cachedFisProductionOfflines[workplaceState.OID]
		if !stateExported {
			logInfo("MAIN", "Exporting new state with OID "+strconv.Itoa(workplaceState.OID))
			var fisProduction FisProduction
			fisProduction.ZapsiId = sql.NullInt32{Int32: int32(workplaceState.OID), Valid: true}
			fisProduction.DatumCasOd = workplaceState.DTS
			fisProduction.DatumCasDo = workplaceState.DTE.Time
			var workplace Workplace
			db.Where("OID = ?", workplaceState.WorkplaceID).Find(&workplace)
			fisProduction.IFS = sql.NullString{String: workplace.Code, Valid: true}
			fisProduction.Stav = sql.NullString{String: "v", Valid: true}
			db.Create(&fisProduction)
			noOfExportedPoweroffs++
		}
	}
	logInfo("MAIN", "Exported "+strconv.Itoa(noOfExportedPoweroffs)+" new poweroff states in "+time.Since(start).String())
}

func exportIdlesFromZapsi(db *gorm.DB) {
	start := time.Now()
	logInfo("MAIN", "Exporting idles started")
	var terminalInputIdles []TerminalInputIdle
	db.Where("dte is not null and dte > ?", time.Now().Add(-24*time.Hour)).Find(&terminalInputIdles)
	var fisProductionIdles []FisProduction
	db.Where("stav like 'p'").Find(&fisProductionIdles)
	cachedFisProductionIdles := make(map[int]FisProduction)
	for _, fisProductionIdle := range fisProductionIdles {
		cachedFisProductionIdles[int(fisProductionIdle.ZapsiId.Int32)] = fisProductionIdle
	}
	noOfExportedIdles := 0
	for _, terminalInputIdle := range terminalInputIdles {
		_, idleExported := cachedFisProductionIdles[terminalInputIdle.OID]
		if !idleExported {
			logInfo("MAIN", "Exporting new idle with OID "+strconv.Itoa(terminalInputIdle.OID))
			var fisProduction FisProduction
			fisProduction.ZapsiId = sql.NullInt32{Int32: int32(terminalInputIdle.OID), Valid: true}
			fisProduction.DatumCasOd = terminalInputIdle.DTS
			fisProduction.DatumCasDo = terminalInputIdle.DTE.Time
			var zapsiUser User
			db.Where("OID = ?", terminalInputIdle.UserID).Find(&zapsiUser)
			if zapsiUser.OID > 1 {
				var fisUser FisUser
				db.Where("IDZ = ?", zapsiUser.Login).Find(&fisUser)
				if fisUser.IDZ > 0 {
					idz, _ := strconv.Atoi(zapsiUser.Login)
					fisProduction.IDZ = sql.NullInt32{Int32: int32(idz), Valid: true}
				} else {
					logInfo("MAIN", "Order with ID "+strconv.Itoa(terminalInputIdle.OID)+" has user with login "+zapsiUser.Login+" that is not in fis table, error")
					fisProduction.IDZ = sql.NullInt32{Valid: true}
					fisProduction.Chyba = sql.NullString{String: zapsiUser.Login, Valid: true}
				}
			}
			var workplace Workplace
			db.Where("DeviceId = ?", terminalInputIdle.DeviceID).Find(&workplace)
			fisProduction.IFS = sql.NullString{String: workplace.Code, Valid: true}
			var terminalInputOrderIdle TerminalInputOrderIdle
			db.Where("TerminalInputIdleID = ?", terminalInputIdle.OID).Find(&terminalInputOrderIdle)
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
					fisProduction.Chyba = sql.NullString{String: fisProduction.Chyba.String + "," + zapsiOrder.Barcode, Valid: true}
				}
			}
			fisProduction.Stav = sql.NullString{String: "p", Valid: true}
			var idle Idle
			db.Where("OID = ?", terminalInputIdle.IdleID).Find(&idle)
			var idleType IdleType
			db.Where("OID = ?", idle.IdleTypeID).Find(&idleType)
			fisProduction.Prostoj = sql.NullString{String: idle.Name, Valid: true}
			fisProduction.TypProstoje = sql.NullString{String: idleType.Name, Valid: true}
			db.Create(&fisProduction)
			noOfExportedIdles++
		}
	}
	logInfo("MAIN", "Exported "+strconv.Itoa(noOfExportedIdles)+" new idles in "+time.Since(start).String())
}

func exportOrdersFromZapsi(db *gorm.DB) {
	start := time.Now()
	logInfo("MAIN", "Exporting orders started")
	var terminalInputOrders []TerminalInputOrder
	db.Where("dte is not null and dte > ?", time.Now().Add(-24*time.Hour)).Find(&terminalInputOrders)
	var fisProductionOrders []FisProduction
	db.Where("stav like 'a'").Find(&fisProductionOrders)
	cachedFisProductionOrders := make(map[int]FisProduction)
	for _, fisProductionOrder := range fisProductionOrders {
		cachedFisProductionOrders[int(fisProductionOrder.ZapsiId.Int32)] = fisProductionOrder
	}
	noOfExportedOrders := 0
	for _, terminalInputOrder := range terminalInputOrders {
		_, orderExported := cachedFisProductionOrders[terminalInputOrder.OID]
		if !orderExported {
			logInfo("MAIN", "Exporting new order with OID "+strconv.Itoa(terminalInputOrder.OID))
			var fisProduction FisProduction
			fisProduction.ZapsiId = sql.NullInt32{Int32: int32(terminalInputOrder.OID), Valid: true}
			fisProduction.DatumCasOd = terminalInputOrder.DTS
			fisProduction.DatumCasDo = terminalInputOrder.DTE.Time
			var zapsiUser User
			db.Where("OID = ?", terminalInputOrder.UserID).Find(&zapsiUser)
			if zapsiUser.OID > 1 {
				var fisUser FisUser
				db.Where("IDZ = ?", zapsiUser.Login).Find(&fisUser)
				if fisUser.IDZ > 0 {
					idz, _ := strconv.Atoi(zapsiUser.Login)
					fisProduction.IDZ = sql.NullInt32{Int32: int32(idz), Valid: true}
				} else {
					logInfo("MAIN", "Order with ID "+strconv.Itoa(terminalInputOrder.OID)+" has user with login "+zapsiUser.Login+" that is not in fis table, error")
					fisProduction.IDZ = sql.NullInt32{Valid: true}
					fisProduction.Chyba = sql.NullString{String: zapsiUser.Login, Valid: true}
				}
			}
			var workplace Workplace
			db.Where("DeviceId = ?", terminalInputOrder.DeviceID).Find(&workplace)
			fisProduction.IFS = sql.NullString{String: workplace.Code, Valid: true}
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
					fisProduction.Chyba = sql.NullString{String: fisProduction.Chyba.String + "," + zapsiOrder.Barcode, Valid: true}
				}
			}
			fisProduction.MnozstviOK = sql.NullInt32{Int32: terminalInputOrder.Count.Int32 - terminalInputOrder.Fail.Int32, Valid: true}
			fisProduction.MnozstviNOK = sql.NullInt32{Int32: terminalInputOrder.Fail.Int32, Valid: true}
			KgOK, _ := strconv.Atoi(terminalInputOrder.Note.String)
			fisProduction.KgOK = sql.NullInt32{Int32: int32(KgOK), Valid: true}
			fisProduction.Stav = sql.NullString{String: "a", Valid: true}
			fisProduction.Takt = sql.NullFloat64{Float64: terminalInputOrder.AverageCycle, Valid: true}
			db.Create(&fisProduction)
			noOfExportedOrders++
		}
	}
	logInfo("MAIN", "Exported "+strconv.Itoa(noOfExportedOrders)+" new orders in "+time.Since(start).String())
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
	ordersMap := make(map[string]Order)
	for _, order := range orders {
		ordersMap[order.Barcode] = order
	}
	noOfUpdatedOrders := 0
	noOfInsertedOrders := 0
	noOfSameOrders := 0
	for _, fisOrder := range fisOrders {
		order, orderInZapsi := ordersMap[strconv.Itoa(fisOrder.ID)]
		if orderInZapsi {
			if order.Name != fisOrder.IDVC || order.CountRequested != fisOrder.Mnozstvi {
				logInfo("MAIN", strconv.Itoa(fisOrder.ID)+": updating order name to "+fisOrder.IDVC+" and count requested to "+strconv.Itoa(fisOrder.Mnozstvi))
				var updateOrder Order
				db.Where("OID = ?", order.OID).Find(&updateOrder)
				updateOrder.Name = fisOrder.IDVC
				updateOrder.CountRequested = fisOrder.Mnozstvi
				db.Where("OID = ?", updateOrder.OID).Save(&updateOrder)
				noOfUpdatedOrders++
			} else {
				noOfSameOrders++
			}
		} else {
			var fisProduct FisProduct
			db.Where("IDVM = ?", fisOrder.IDVM).Find(&fisProduct)
			if fisProduct.IDVM > 0 {
				logInfo("MAIN", strconv.Itoa(fisOrder.ID)+": adding new order with name "+fisOrder.IDVC+" and count requested "+strconv.Itoa(fisOrder.Mnozstvi))
				var product Product
				db.Where("barcode = ?", fisProduct.ArtNr).Find(&product)
				var newOrder Order
				newOrder.Name = fisOrder.IDVC
				newOrder.Barcode = strconv.Itoa(fisOrder.ID)
				newOrder.ProductID = product.OID
				newOrder.OrderStatusID = 1
				newOrder.CountRequested = fisOrder.Mnozstvi
				newOrder.Cavity = 1
				db.Create(&newOrder)
				noOfInsertedOrders++
			}
		}
	}
	logInfo("MAIN", "Updated "+strconv.Itoa(noOfUpdatedOrders)+" orders, inserted "+strconv.Itoa(noOfInsertedOrders)+" orders, skipped "+strconv.Itoa(noOfSameOrders)+" orders")
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
	productsMap := make(map[string]Product)
	for _, product := range products {
		productsMap[product.Barcode] = product
	}
	noOfUpdatedProducts := 0
	noOfInsertedProducts := 0
	noOfSameProducts := 0
	for _, fisProduct := range fisProducts {
		product, productInZapsi := productsMap[fisProduct.ArtNr]
		if productInZapsi {
			if product.Name != fisProduct.Nazev+" "+fisProduct.Velikost {
				logInfo("MAIN", fisProduct.ArtNr+": updating product name to "+fisProduct.Nazev)
				var updateProduct Product
				db.Where("OID = ?", product.OID).Find(&updateProduct)
				updateProduct.Name = fisProduct.Nazev + " " + fisProduct.Velikost
				db.Where("OID = ?", updateProduct.OID).Save(&updateProduct)
				noOfUpdatedProducts++
			} else {
				noOfSameProducts++
			}
		} else {
			logInfo("MAIN", fisProduct.ArtNr+": adding new product with name "+fisProduct.Nazev)
			var newProduct Product
			newProduct.Name = fisProduct.Nazev + " " + fisProduct.Velikost
			newProduct.Barcode = fisProduct.ArtNr
			newProduct.Cycle = 0.0
			newProduct.ProductStatusID = 1
			newProduct.Deleted = 0
			db.Create(&newProduct)
			noOfInsertedProducts++
		}
	}
	logInfo("MAIN", "Updated "+strconv.Itoa(noOfUpdatedProducts)+" products, inserted "+strconv.Itoa(noOfInsertedProducts)+" products, skipped "+strconv.Itoa(noOfSameProducts)+" products")
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
	usersMap := make(map[string]User)
	for _, user := range users {
		usersMap[user.Login] = user
	}
	noOfUpdatedUsers := 0
	noOfInsertedUsers := 0
	noOfSameUsers := 0
	for _, fisUser := range fisUsers {
		user, userInZapsi := usersMap[strconv.Itoa(fisUser.IDZ)]
		if userInZapsi {
			if user.Rfid.String != fisUser.Rfid {
				logInfo("MAIN", "["+strconv.Itoa(fisUser.IDZ)+"] "+fisUser.Jmeno+" "+fisUser.Prijmeni+": updating user with rfid "+fisUser.Rfid)
				var updateUser User
				db.Where("OID = ?", user.OID).Find(&updateUser)
				updateUser.Rfid = sql.NullString{String: fisUser.Rfid, Valid: true}
				db.Where("OID = ?", updateUser.OID).Save(&updateUser)
				noOfUpdatedUsers++
			} else {
				noOfSameUsers++
			}
		} else {
			logInfo("MAIN", "["+strconv.Itoa(fisUser.IDZ)+"] "+fisUser.Jmeno+" "+fisUser.Prijmeni+": adding new user with  rfid ["+fisUser.Rfid+"]")
			var newUser User
			newUser.Login = strconv.Itoa(fisUser.IDZ)
			newUser.Name = sql.NullString{String: fisUser.Prijmeni, Valid: true}
			newUser.FirstName = sql.NullString{String: fisUser.Jmeno, Valid: true}
			newUser.Rfid = sql.NullString{String: fisUser.Rfid, Valid: true}
			newUser.UserRoleID = 2
			db.Create(&newUser)
			noOfInsertedUsers++
		}
	}
	logInfo("MAIN", "Updated "+strconv.Itoa(noOfUpdatedUsers)+" users, inserted "+strconv.Itoa(noOfInsertedUsers)+" users, skipped "+strconv.Itoa(noOfSameUsers)+" users")
	logInfo("MAIN", "Importing users ended in "+time.Since(start).String())
}
