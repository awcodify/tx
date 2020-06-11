# Tx
Transactions enforce the integrity of the database and guard the data against program errors or database break-downs [[=>]](https://api.rubyonrails.org/classes/ActiveRecord/Transactions/ClassMethods.html). 

## How to use
Wrap your transaction inside `Wrap` function. It determine between `Rollback()` and `Commit()` . See example below:
```go
package main

func main() {
	// ... some codes
	tx : tx.NewTx(db)
	err := tx.Wrap(createUserAndRole)
	if err != nil {
		log.Fatal(err)
	}
	// ... some codes
}

// This our needs
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
```
