package controllers

import (
	"context"
	"product-tracker/db"
	"product-tracker/config"
)

func HealthCheck(c context.Context) error {
    cfg := config.GetConfig()
    dbConn, err := db.NewDB(cfg.GetDSN(), cfg.Database.MaxConns, cfg.Database.MaxIdle, cfg.Database.Timeout)
    if err != nil {
        return err
    }
    defer dbConn.Close()

    err = dbConn.Ping(c)
    if err != nil {
		return err
	}
	return nil
}