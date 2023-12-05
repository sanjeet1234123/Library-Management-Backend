package main

import (
	"context"
	"database/sql"
	"example/API-task.go/models"
	"net/http"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func createUser(c *gin.Context) {
	db := c.MustGet("db").(*sql.DB)
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error1": err.Error()})
		return
	}

	tx, err := db.Begin()
	if err != nil {
		c.JSON(500, gin.H{"error": "Unable to begin transaction"})
		return
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var id3 string
	err3 := tx.QueryRow("SELECT role FROM users WHERE lib_id = $1", user.LibId).Scan(&id3) // checking the role of library with lib id

	// inserting admin role if admin not exists
	if err3 != nil && id3 != "admin" {
		stmt, err1 := tx.Prepare("INSERT INTO users (name, email, contact_number, role, lib_id, password) VALUES ($1, $2, $3, 'admin', $4, $5)")
		if err1 != nil {
			tx.Rollback()
			c.JSON(500, gin.H{"error": err1.Error()})
			return
		}
		defer stmt.Close()
		_, err2 := stmt.Exec(user.Name, user.Email, user.ContactNumber, user.LibId, user.Password)
		if err2 != nil {
			tx.Rollback()
			c.JSON(500, gin.H{"error2": err2.Error()})
			return
		}
	}

	// inserting reader role if admin exists
	stmt, err1 := tx.Prepare("INSERT INTO users (name, email, contact_number, role, lib_id, password) VALUES ($1, $2, $3, 'reader', $4, $5)")
	if err1 != nil {
		tx.Rollback()
		c.JSON(500, gin.H{"error": err1.Error()})
		return
	}
	defer stmt.Close()
	_, err2 := stmt.Exec(user.Name, user.Email, user.ContactNumber, user.LibId, user.Password)
	if err2 != nil {
		tx.Rollback()
		c.JSON(500, gin.H{"error2": err2.Error()})
		return
	}

	err = tx.Commit()
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(201, user)
}
func loginUser(c *gin.Context) {
	db := c.MustGet("db").(*sql.DB)

	var loginData struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	query := "SELECT id, name, email, contact_number, role, lib_id FROM users WHERE email = $1 AND password = $2"
	err := db.QueryRow(query, loginData.Email, loginData.Password).
		Scan(&user.ID, &user.Name, &user.Email, &user.ContactNumber, &user.Role, &user.LibId)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(401, gin.H{"error": "Invalid email or password"})
			return
		}

		c.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(200, gin.H{"message": "Login successful", "user": user, "role": user.Role})
}
//function to get all users
func getAllUser(c *gin.Context) {
    db := c.MustGet("db").(*sql.DB)
    email := c.Param("email")
    var role string
    err := db.QueryRow("SELECT role FROM Users WHERE email = $1", email).Scan(&role)
    if err != nil || role != "admin" {
        c.JSON(403, gin.H{"error": "Not authorized as Admin"})
        return
    }

    rows, err := db.Query("SELECT id, name, email, contact_number, role, lib_id FROM users")
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }
    defer rows.Close()

    users := []models.User{}
    for rows.Next() {
        var user models.User
        if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.ContactNumber, &user.Role, &user.LibId); err != nil {
            c.JSON(500, gin.H{"error": err.Error()})
            return
        }
        users = append(users, user)
    }

    c.JSON(200, users)
}

//updating  a user by admin
func updateUserByID(c *gin.Context) {

	id := c.Param("id")

	var user models.User

	db := c.MustGet("db").(*sql.DB)

	if err1 := c.ShouldBindJSON(&user); err1 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error1": "Invalid data"})
		return
	}



	stmt, err := db.Prepare("UPDATE users SET role=$1 WHERE id=$2")

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer stmt.Close()

	result, err := stmt.Exec(user.Role, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	c.JSON(http.StatusOK, user.Role)
}

func createLibrary(c *gin.Context) {

	db := c.MustGet("db").(*sql.DB)                                     //
	var libraryAdmin models.LibraryAdmin


	if err := c.ShouldBindJSON(&libraryAdmin); err != nil {
		c.JSON(400, gin.H{"error": "Invalid data"})
		return
	}

	db.Query("")

	tx, err := db.Begin()
	if err != nil {
		c.JSON(500, gin.H{"error": "Unable to begin transaction"})
		return
	}

	
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	stmt, err := tx.Prepare("INSERT INTO library (name) VALUES ($1)")
	if err != nil {
		tx.Rollback()
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer stmt.Close()

	result, err := stmt.Exec(libraryAdmin.Name)
	if err != nil {
		tx.Rollback()
		c.JSON(500, gin.H{"error2": err.Error()})
		return
	}
	// Fetch the last inserted ID
	id, _ := result.LastInsertId()
	libraryAdmin.ID = int(id)


	err = tx.Commit()
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to commit transaction"})
		return
	}
	c.JSON(201, libraryAdmin)
}

func getAllLibrary(c *gin.Context) {
	db := c.MustGet("db").(*sql.DB)

	rows, err := db.Query("SELECT id, name FROM library")
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	libraryAdmins := []models.LibraryAdmin{} //an empty slice with type models.libra....to store the retrieved library data.
	for rows.Next() {
		var libraryAdmin models.LibraryAdmin
		if err := rows.Scan(&libraryAdmin.ID, &libraryAdmin.Name); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		libraryAdmins = append(libraryAdmins, libraryAdmin)
	}

	c.JSON(200, libraryAdmins)
}

func createBook(c *gin.Context) {
	db := c.MustGet("db").(*sql.DB)
	var book models.Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(400, gin.H{"error1": err.Error()})
		return
	}

	tx, err := db.Begin()
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to start transaction"})
		return
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()


	//check if book already exists
	row := tx.QueryRow("SELECT totalcopies, availablecopies FROM bookinventory WHERE isbn = $1 AND libid = $2", book.ISBN, book.LibID)
	var existingTotalCopies, existingAvailableCopies int
	if err := row.Scan(&existingTotalCopies, &existingAvailableCopies); err == nil {
        
	//yes then update the total and available copies
		_, err := tx.Exec("UPDATE bookinventory SET totalcopies = $1, availablecopies = $2 WHERE isbn = $3 AND libid = $4", existingTotalCopies+book.TotalCopies, existingAvailableCopies+book.AvailableCopies, book.ISBN, book.LibID)
		if err != nil {
			tx.Rollback()
			c.JSON(500, gin.H{"error2": err.Error()})
			return
		}
		err = tx.Commit()
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to commit transaction"})
			return
		}
		c.JSON(200, gin.H{"message": "Book updated successfully"})
		return
	}

	// Book does not exist, insert new entry
	stmt, err := tx.Prepare("INSERT INTO bookinventory (isbn, libid, title, authors, publisher, version, totalcopies, availablecopies) VALUES ($1, $2, $3, $4, $5, $6 ,$7, $8)")
	if err != nil {
		tx.Rollback()
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(book.ISBN, book.LibID, book.Title, book.Author, book.Publisher, book.Version, book.TotalCopies, book.AvailableCopies)
	if err != nil {
		tx.Rollback()
		c.JSON(500, gin.H{"error2": err.Error()})
		return
	}

	err = tx.Commit()
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(201, book)
}


func getAllBook(c *gin.Context) {
	db := c.MustGet("db").(*sql.DB)

	tx, err := db.Begin()
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to start transaction"})
		return
	}

	rows, err := tx.Query("SELECT isbn, libid, title, authors, publisher, version, totalcopies, availablecopies FROM bookinventory")
	if err != nil {
		tx.Rollback()
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	books := []models.Book{}
	for rows.Next() {
		var book models.Book
		if err := rows.Scan(&book.ISBN, &book.LibID, &book.Title, &book.Author, &book.Publisher, &book.Version, &book.TotalCopies, &book.AvailableCopies); err != nil {
			tx.Rollback()
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		books = append(books, book)
	}

	tx.Commit()

	c.JSON(200, books)
}

func updateBookByISBN(c *gin.Context) {
	ISBN := c.Param("isbn")
	db := c.MustGet("db").(*sql.DB)
	var book models.Book

	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	tx, err := db.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var isbnChecker string
	err = tx.QueryRow("SELECT isbn FROM bookinventory WHERE title = $1", book.Title).Scan(&isbnChecker)
	if isbnChecker != ISBN {
		tx.Rollback()
		c.JSON(403, gin.H{"error": "please enter correct isbn in url"})
		return
	}

	//update the book inventory
	stmt, err := tx.Prepare("UPDATE bookinventory SET title=$1, authors=$2, publisher=$3, version=$4, totalcopies=$5, availablecopies=$6 WHERE isbn=$7")
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer stmt.Close() // it will close until the surrounding function returns.

	_, err = stmt.Exec(book.Title, book.Author, book.Publisher, book.Version, book.TotalCopies, book.AvailableCopies, ISBN)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	
	err = tx.Commit()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, book)
}

func deleteBookByISBN(c *gin.Context) {
	isbn := c.Param("isbn")  
	 

	db := c.MustGet("db").(*sql.DB)


	tx, err := db.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
   
	copiesToRemove := 1 

	//retrieve the current counts of the book
	row := tx.QueryRow("SELECT totalcopies, availablecopies FROM bookinventory WHERE isbn = $1", isbn)
	var existingTotalCopies, existingAvailableCopies int
	if err := row.Scan(&existingTotalCopies, &existingAvailableCopies); err != nil {
		tx.Rollback()
		c.JSON(500, gin.H{"error": "Error retrieving book details"})
		return
	}

	//checking if enough available copies are present or not
	if copiesToRemove > existingAvailableCopies {
		tx.Rollback()
		c.JSON(400, gin.H{"error": "Cannot remove books that are currently issued"})
		return
	}

	//updating the copies count
	_, err = tx.Exec("UPDATE bookinventory SET totalcopies = totalcopies - $1, availablecopies = availablecopies - $1 WHERE isbn = $2", copiesToRemove, isbn)
	if err != nil {
		tx.Rollback()
		c.JSON(500, gin.H{"error": "Error updating book inventory"})
		return
	}

	// optionally remove book entry if totalcopies is 0
	if existingTotalCopies-copiesToRemove <= 0 {
		_, err := tx.Exec("DELETE FROM bookinventory WHERE isbn = $1", isbn)
		if err != nil {
			tx.Rollback()
			c.JSON(500, gin.H{"error": "Error removing book from inventory"})
			return
		}
	}

	
	err = tx.Commit()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(200, gin.H{"message": "Book removed successfully"})
}

func getBookByISBN(c *gin.Context) {
	isbn := c.Param("isbn")

	db := c.MustGet("db").(*sql.DB)

	tx, err := db.BeginTx(context.TODO(), &sql.TxOptions{
		Isolation: sql.LevelSerializable,
		ReadOnly:  true,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}
	defer tx.Rollback()

	var book models.Book

	err = tx.QueryRow("SELECT isbn, libid, title, authors, publisher, version, totalcopies, availablecopies FROM bookinventory WHERE isbn = $1;", isbn).Scan(&book.ISBN, &book.LibID, &book.Title, &book.Author, &book.Publisher, &book.Version, &book.TotalCopies, &book.AvailableCopies)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(404, gin.H{"error": "Book not found"})
			return
		}
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, book)
}

func createRequest(c *gin.Context) {
	db := c.MustGet("db").(*sql.DB)
	var requestEvents models.RequestEvents

	if err := c.ShouldBindJSON(&requestEvents); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	tx, err := db.Begin()
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to start transaction"})
		return
	}


	stmt, err := tx.Prepare("INSERT INTO requestevents (bookid, readerid, requestdate, approvaldate, approverid, requesttype) VALUES ($1, $2, $3, $4, $5, $6)")
	if err != nil {
		tx.Rollback()
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer stmt.Close()

	result, err := stmt.Exec(requestEvents.BookID, requestEvents.ReaderID, requestEvents.RequestDate, requestEvents.ApprovalDate, requestEvents.ApproverID, requestEvents.RequestType)
	if err != nil {
		tx.Rollback()
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	id, _ := result.LastInsertId()
	requestEvents.ReqID = int(id)

	err = tx.Commit()
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(201, requestEvents)
}
func getAllRequest(c *gin.Context) {
	db := c.MustGet("db").(*sql.DB)

	tx, err := db.Begin()
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to start transaction"})
		return
	}

	rows, err := tx.Query("SELECT * FROM requestevents")
	if err != nil {
		tx.Rollback()
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	requestEvents := []models.RequestEvents{}
	for rows.Next() {
		var requestEvent models.RequestEvents
		if err := rows.Scan(&requestEvent.ReqID, &requestEvent.BookID, &requestEvent.ReaderID, &requestEvent.RequestDate, &requestEvent.ApprovalDate, &requestEvent.ApproverID, &requestEvent.RequestType); err != nil {
			tx.Rollback() 
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		requestEvents = append(requestEvents, requestEvent)
	}

	tx.Commit()

	c.JSON(200, requestEvents)
}

func updateRequestByReqID(c *gin.Context) {
    db := c.MustGet("db").(*sql.DB)
    ReqID := c.Param("reqid")
    var requestEvents models.RequestEvents

    if err := c.ShouldBindJSON(&requestEvents); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
        return
    }

    tx, err := db.Begin()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
        return
    }

    stmt, err := tx.Prepare("UPDATE requestevents SET approvaldate=$1, approverid=$2, requesttype=$3 WHERE reqid=$4")
    if err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    defer stmt.Close()

    result, err := stmt.Exec(requestEvents.ApprovalDate, requestEvents.ApproverID, requestEvents.RequestType, ReqID)
    if err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // If no rows affected, rollback
    if rowsAffected == 0 {
        tx.Rollback()
        c.JSON(http.StatusNotFound, gin.H{"error": "Request not found"})
        return
    }

    tx.Commit()
    c.JSON(http.StatusOK, ReqID)
}
func createIssue(c *gin.Context) {
	db := c.MustGet("db").(*sql.DB)

	var issueRegistry models.IssueRegistery

	if err := c.ShouldBindJSON(&issueRegistry); err != nil {
		c.JSON(400, gin.H{"error1": err.Error()})
		return
	}

	tx, err := db.Begin()
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to start transaction"})
		return
	}

	var availableCopies int
	row := tx.QueryRow("SELECT availablecopies FROM bookinventory WHERE isbn = $1", issueRegistry.ISBN)
	if err := row.Scan(&availableCopies); err != nil {
		tx.Rollback()
		c.JSON(500, gin.H{"error": "Error retrieving book details"})
		return
	}
	// checking book is there for issue or not
	if availableCopies <= 0 {
		tx.Rollback()
		c.JSON(400, gin.H{"error": "No available copies for the requested ISBN"})
		return
	}
    //yes then updating the available copies-1
	_, err = tx.Exec("UPDATE bookinventory SET availablecopies = availablecopies - 1 WHERE isbn = $1", issueRegistry.ISBN)
	if err != nil {
		tx.Rollback()
		c.JSON(500, gin.H{"error": "Error updating book inventory"})
		return
	}
    
	stmt, err := tx.Prepare("INSERT INTO issueregistry ( isbn, readerid, issueapproverid, issuestatus, issuedate, expectedreturndate, returndate, returnapproverid) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)")
	if err != nil {
		tx.Rollback()
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer stmt.Close()
    
	result, err := stmt.Exec(issueRegistry.ISBN, issueRegistry.ReaderID, issueRegistry.IssueApproverID, issueRegistry.IssueStatus, issueRegistry.IssueDate, issueRegistry.ExpectedReturnDate, issueRegistry.ReturnDate, issueRegistry.ReturnApproverID)
	if err != nil {
		tx.Rollback()
		c.JSON(500, gin.H{"error2": err.Error()})
		return
	}

	id, _ := result.LastInsertId()
	issueRegistry.IssueID = int(id)
	tx.Commit()
	c.JSON(201, issueRegistry)
}
