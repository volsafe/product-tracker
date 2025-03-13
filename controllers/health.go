package controllers

import (
	"context"
	"product-tracker/config"
	"product-tracker/db"
)

func HealthCheck(c context.Context) error {
	cfg := config.GetConfig()
	dbConfig := &db.DBConfig{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		DbName:   cfg.Database.DbName,
	}

	dbConn, err := db.NewDB(dbConfig)
	if err != nil {
		return err
	}
	defer dbConn.Close()

	return dbConn.Ping()
}
