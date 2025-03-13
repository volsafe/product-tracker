package controllers

import (
	"context"
	"product-tracker/models"
	"product-tracker/storage"
)

var S *storage.Storage

func SetStorageInstance(storageInstance *storage.Storage) {
	S = storageInstance
}

type Product struct {
	Name              string  `json:"name"`
	Description       string  `json:"description"`
	Price             float64 `json:"price"`
	EnergyConsumption float64 `json:"energy_consumption"`
}

func InsertProduct(c context.Context, p Product) error {
	err := S.InsertProduct(c, &models.Product{
		Name:              p.Name,
		Description:       p.Description,
		Price:             p.Price,
		EnergyConsumption: p.EnergyConsumption,
	})
	if err != nil {
		return err
	}
	return nil
}

func GetProductsByName(c context.Context, name string) ([]Product, error) {
	products, err := S.GetProductsByName(c, name)
	if err != nil {
		return nil, err
	}
	var p []Product
	for _, product := range products {
		p = append(p, Product{
			Name:              product.Name,
			Description:       product.Description,
			Price:             product.Price,
			EnergyConsumption: product.EnergyConsumption,
		})
	}
	return p, nil
}

func GetProducts(c context.Context) ([]Product, error) {
	products, err := S.GetProducts(c)
	if err != nil {
		return nil, err
	}
	var p []Product
	for _, product := range products {
		p = append(p, Product{
			Name:              product.Name,
			Description:       product.Description,
			Price:             product.Price,
			EnergyConsumption: product.EnergyConsumption,
		})
	}
	return p, nil
}
