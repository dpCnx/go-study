package test

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestShouldUpdateStats(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"name"}).AddRow("d").AddRow("p")

	mock.ExpectQuery("select * from user1").WillReturnRows(rows)

	r, err := db.Query("select * from user1")
	if err != nil {
		t.Errorf("db quary err: %v", err)
	}

	defer r.Close()

	for r.Next() {
		var n string
		if err := r.Scan(&n); err != nil {
			t.Errorf("r next err: %v", err)
		}

		t.Log(n)
	}

}

func TestSomething(t *testing.T) {

	assert.Equal(t, 123, 123, "they should be equal")


}
