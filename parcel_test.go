package main

import (
	"database/sql"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var (
	randSource = rand.NewSource(time.Now().UnixNano())
	randRange  = rand.New(randSource)
)

func getTestParcel() Parcel {
	return Parcel{
		Client:    1000 + randRange.Intn(1000), // случайный клиент от 1000 до 1999
		Status:    ParcelStatusRegistered,
		Address:   fmt.Sprintf("test-%d", randRange.Intn(10000)), // случайный адрес
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
}
func TestAddGetDelete(t *testing.T) {
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {

		}
	}(db)
	store := NewParcelStore(db)
	parcel := getTestParcel()
	id, err := store.Add(parcel)
	require.NoError(t, err)
	require.Greater(t, id, 0)
	p, err := store.Get(id)
	require.NoError(t, err)
	require.Equal(t, parcel.Client, p.Client)
	require.Equal(t, parcel.Status, p.Status)
	require.Equal(t, parcel.Address, p.Address)
	require.NotEmpty(t, p.CreatedAt)
	err = store.Delete(id)
	require.NoError(t, err)
	p, err = store.Get(id)
	require.Error(t, err)
}
func TestSetAddress(t *testing.T) {
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
		}
	}(db)
	store := NewParcelStore(db)
	parcel := getTestParcel()
	id, err := store.Add(parcel)
	require.NoError(t, err)
	require.Greater(t, id, 0)
	newAddress := "new test address"
	err = store.SetAddress(id, newAddress)
	require.NoError(t, err)
	p, err := store.Get(id)
	require.NoError(t, err)
	require.Equal(t, newAddress, p.Address)
}
func TestSetStatus(t *testing.T) {
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
		}
	}(db)
	store := NewParcelStore(db)
	parcel := getTestParcel()
	id, err := store.Add(parcel)
	require.NoError(t, err)
	require.Greater(t, id, 0)
	err = store.SetStatus(id, ParcelStatusSent)
	require.NoError(t, err)
	p, err := store.Get(id)
	require.NoError(t, err)
	require.Equal(t, ParcelStatusSent, p.Status)
}
func TestGetByClient(t *testing.T) {
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
		}
	}(db)
	store := NewParcelStore(db)
	client := 123456
	parcels := []Parcel{
		getTestParcel(),
		getTestParcel(),
		getTestParcel(),
	}
	parcelMap := map[int]Parcel{}
	for i := range parcels {
		parcels[i].Client = client
		id, err := store.Add(parcels[i])
		require.NoError(t, err)
		parcels[i].Number = id
		parcelMap[id] = parcels[i]
	}
	storedParcels, err := store.GetByClient(client)
	require.NoError(t, err)
	require.Equal(t, len(parcels), len(storedParcels))
	for _, p := range storedParcels {
		expected, ok := parcelMap[p.Number]
		require.True(t, ok)
		require.Equal(t, expected.Client, p.Client)
		require.Equal(t, expected.Address, p.Address)
		require.Equal(t, expected.Status, p.Status)
		require.NotEmpty(t, p.CreatedAt)
	}
}
