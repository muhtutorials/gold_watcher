package repository

import (
	"testing"
	"time"
)

func TestSQLiteRepository_Migrate(t *testing.T) {
	err := testRepo.Migrate()
	if err != nil {
		t.Error("migrate failed:", err)
	}
}

func TestSQLiteRepository_InsertHolding(t *testing.T) {
	h := Holdings{
		Amount:        1,
		PurchaseDate:  time.Now(),
		PurchasePrice: 1000,
	}

	result, err := testRepo.InsertHolding(h)
	if err != nil {
		t.Error("insert failed:", err)
	}

	if result.ID <= 0 {
		t.Error("invalid ID sent:", result.ID)
	}
}

func TestSQLiteRepository_AllHoldings(t *testing.T) {
	h, err := testRepo.AllHoldings()
	if err != nil {
		t.Error("get all holdings failed:", err)
	}

	if len(h) != 1 {
		t.Error("wrong number of rows returned; expected 1 but got", len(h))
	}
}

func TestSQLiteRepository_GetHoldingByID(t *testing.T) {
	h, err := testRepo.GetHoldingByID(1)
	if err != nil {
		t.Error("get by ID failed:", err)
	}

	if h.PurchasePrice != 1000 {
		t.Error("wrong purchase price returned; expected 1000 but got", h.PurchasePrice)
	}

	_, err = testRepo.GetHoldingByID(2)
	if err == nil {
		t.Error("get by ID returned value for nonexistent ID")
	}
}

func TestSQLiteRepository_UpdateHolding(t *testing.T) {
	h, err := testRepo.GetHoldingByID(1)
	if err != nil {
		t.Error("get by ID failed:", err)
	}

	h.PurchasePrice = 1001

	err = testRepo.UpdateHolding(1, *h)
	if err != nil {
		t.Error("update failed:", err)
	}
}

func TestSQLiteRepository_DeleteHolding(t *testing.T) {
	err := testRepo.DeleteHolding(1)
	if err != nil {
		t.Error("delete failed:", err)
		if err != errDeleteFailed {
			t.Error("wrong error returned")
		}
	}

	err = testRepo.DeleteHolding(2)
	if err == nil {
		t.Error("no error when trying to delete nonexistent record")
	}
}
