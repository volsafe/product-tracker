package controllers


import (
	"context"
	"product-tracker/storage"

)


var S *storage.Storage

func SetStorageInstance(storageInstance *storage.Storage) {
	S = storageInstance
}


type Product struct {
	Name  string `json:"name"`
	Quantity int `json:"quantity"`
	EnergyConsumed float64 `json:"energy_consumed"`
	Date string `json:"date"`
}


func InsertProduct(c context.Context, p Product) error {
	err := S.InsertProduct(c, storage.Product{
		Name: p.Name,
		Quantity: p.Quantity,
		EnergyConsumed: p.EnergyConsumed,
		Date: p.Date,
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
			Name: product.Name,
			Quantity: product.Quantity,
			EnergyConsumed: product.EnergyConsumed,
			Date: product.Date,
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
			Name: product.Name,
			Quantity: product.Quantity,
			EnergyConsumed: product.EnergyConsumed,
			Date: product.Date,
		})
	}
	return p, nil
}