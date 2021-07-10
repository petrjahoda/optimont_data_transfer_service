package main

import (
	"github.com/kardianos/service"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strconv"
	"time"
)

const version = "2021.3.1.10"
const serviceName = "Optimont data transfer service"
const serviceDescription = "Transfers data between zapsi and optimont tables"

//const config = "user=postgres password=pj79.. dbname=system host=database port=5432 sslmode=disable application_name=system_service"
const config = "zapsi_uzivatel:zapsi@tcp(zapsidatabase:3306)/zapsi2?charset=utf8mb4&parseTime=True&loc=Local"
const downloadInSeconds = 60

var serviceRunning = false
var processRunning = false

type program struct{}

func main() {
	logInfo("MAIN", serviceName+" ["+version+"] starting...")
	logInfo("MAIN", "© "+strconv.Itoa(time.Now().Year())+" Petr Jahoda")
	serviceConfig := &service.Config{
		Name:        serviceName,
		DisplayName: serviceName,
		Description: serviceDescription,
	}
	program := &program{}
	s, err := service.New(program, serviceConfig)
	if err != nil {
		logError("MAIN", "Cannot start: "+err.Error())
	}
	err = s.Run()
	if err != nil {
		logError("MAIN", "Cannot start: "+err.Error())
	}
}

func (p *program) Start(service.Service) error {
	logInfo("MAIN", serviceName+" ["+version+"] started")
	serviceRunning = true
	go p.run()
	return nil
}

func (p *program) Stop(service.Service) error {
	serviceRunning = false
	if processRunning {
		logInfo("MAIN", serviceName+" ["+version+"] stopping...")
		time.Sleep(1 * time.Second)
	}
	logInfo("MAIN", serviceName+" ["+version+"] stopped")
	return nil
}

func (p *program) run() {
	for serviceRunning {
		processRunning = true
		start := time.Now()
		logInfo("MAIN", serviceName+" ["+version+"] running")
		db, err := gorm.Open(mysql.Open(config), &gorm.Config{})
		if err != nil {
			logError("MAIN", "Problem opening database: "+err.Error())
			sleepTime := downloadInSeconds*time.Second - time.Since(start)
			logInfo("MAIN", "Sleeping for "+sleepTime.String())
			time.Sleep(sleepTime)
			processRunning = false
			continue
		}
		sqlDB, _ := db.DB()
		importUsersToZapsi(db)
		importProductsToZapsi(db)
		importOrdersToZapsi(db)
		//exportOrdersFromZapsi(db)
		//exportIdlesFromZapsi(db)
		//exportStatePowerOffFromZapsi(db)
		sqlDB.Close()
		sleepTime := downloadInSeconds*time.Second - time.Since(start)
		logInfo("MAIN", "Sleeping for "+sleepTime.String())
		time.Sleep(sleepTime)
		processRunning = false
	}
}
