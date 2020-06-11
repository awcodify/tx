package tx

import (
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

var (
	db, mock, err = sqlmock.New()
)

// In this test case we create into two tables (user and role)
//  If there's no any errors, it should be committed to user and role table
//	If any error present, it should be rolled back
//	Tou can go to createAndUserRole function to get more context
func TestWrap(t *testing.T) {
	t.Run("All transactions is run without any error", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO users\\(name,email\\) VALUES(.+)").
			WithArgs("example", "example@email.com").
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec("INSERT INTO role\\(user_id,privileges\\) VALUES(.+)").
			WithArgs(1, "SELECT,INSERT").
			WillReturnResult(sqlmock.NewResult(1, 1))
		// Since all transaction is works, it will commit to the database.
		mock.ExpectCommit()

		tx := NewTx(db)
		err := tx.Wrap(createUserAndRole)

		assert.NoError(t, err)
	})

	t.Run("Failed to begin the transaction", func(t *testing.T) {
		mock.ExpectBegin().WillReturnError(fmt.Errorf("Failed to begin the transaction"))

		tx := NewTx(db)
		err := tx.Wrap(createUserAndRole)

		assert.EqualError(t, err, "Failed to begin the transaction")
	})

	t.Run("Some transaction is failed", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO users\\(name,email\\) VALUES(.+)").
			WithArgs("example", "example@email.com").
			WillReturnError(fmt.Errorf("Failed to insert into example table"))
		// Since some error present, then it should rollback
		mock.ExpectRollback()

		tx := NewTx(db)
		err := tx.Wrap(createUserAndRole)

		assert.EqualError(t, err, "Failed to insert into example table")
	})
}

func createUserAndRole(t Transaction) (err error) {
	res, err := t.Exec("INSERT INTO users(name,email) VALUES(?,?)", "example", "example@email.com")
	if err != nil {
		return
	}

	userId, err := res.LastInsertId()
	if err != nil {
		return
	}

	res, err = t.Exec("INSERT INTO role(user_id,privileges) VALUES(?,?)", userId, "SELECT,INSERT")
	if err != nil {
		return
	}

	return
}
