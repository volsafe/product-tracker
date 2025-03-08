package controllers

import (
	"context"
	"profile/db"
)

func HealthCheck(c context.Context) error {
    dbConn, err := db.NewDB()
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