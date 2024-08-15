package database

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/valdir-alves3000/postgraduate-challenges-go-expert/APIs/internal/entity"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func getDBProduct() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(&entity.Product{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func TestCreateNewProduct(t *testing.T) {
	db, err := getDBProduct()
	if err != nil {
		t.Error(err)
	}
	product, err := entity.NewProduct("Product 1", 12.0)

	assert.NoError(t, err)
	productDB := NewProduct(db)
	err = productDB.Create(product)
	assert.NoError(t, err)
	assert.NotEmpty(t, product.ID)
}

func TestFindAllProducts(t *testing.T) {
	db, err := getDBProduct()
	if err != nil {
		t.Error(err)
	}
	for i := 1; i < 24; i++ {
		product, err := entity.NewProduct(fmt.Sprintf("Product %d", i), rand.Float64()*100)
		assert.NoError(t, err)
		db.Create(product)
	}
	productDB := NewProduct(db)
	products, err := productDB.FindAll(1, 10, "")
	assert.NoError(t, err)
	assert.Len(t, products, 10)
	assert.Equal(t, "Product 1", products[0].Name)
	assert.Equal(t, "Product 10", products[9].Name)

	products, err = productDB.FindAll(2, 10, "asc")
	assert.NoError(t, err)
	assert.Len(t, products, 10)
	assert.Equal(t, "Product 11", products[0].Name)
	assert.Equal(t, "Product 20", products[9].Name)

	products, err = productDB.FindAll(3, 10, "any-value")
	assert.NoError(t, err)
	assert.Len(t, products, 3)
	assert.Equal(t, "Product 21", products[0].Name)
	assert.Equal(t, "Product 23", products[2].Name)
}

func TestFindProductByID(t *testing.T) {
	db, err := getDBProduct()
	if err != nil {
		t.Error(err)
	}
	product, err := entity.NewProduct("Product 1", 12.0)
	assert.NoError(t, err)
	db.Create(product)
	productDB := NewProduct(db)

	product, err = productDB.FindByID(product.ID.String())
	assert.NoError(t, err)
	assert.Equal(t, "Product 1", product.Name)
}

func TestUpdateProduct(t *testing.T) {
	db, err := getDBProduct()
	if err != nil {
		t.Error(err)
	}
	product, err := entity.NewProduct("Product 1", 12.0)
	assert.NoError(t, err)
	db.Create(product)

	productDb := NewProduct(db)
	product.Name = "Product 2"
	err = productDb.Update(product)
	assert.NoError(t, err)

	product, err = productDb.FindByID(product.ID.String())
	assert.NoError(t, err)
	assert.Equal(t, "Product 2", product.Name)
}

func TestDeleteProduct(t *testing.T) {
	db, err := getDBProduct()
	if err != nil {
		t.Error(err)
	}
	product, err := entity.NewProduct("Product 1", 12.0)
	assert.NoError(t, err)
	db.Create(product)
	productDB := NewProduct(db)

	err = productDB.Delete(product.ID.String())
	assert.NoError(t, err)

	_, err = productDB.FindByID(product.ID.String())
	assert.Error(t, err)
}
