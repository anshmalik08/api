package main

import (
    "bytes"
    "database/sql"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "strconv"
    "testing"

    "github.com/gorilla/mux"
    _ "github.com/go-sql-driver/mysql"
)

func setupTestDB() (*sql.DB, error) {
	dsn := "root:@tcp(localhost:3306)/practice"
    if err != nil {
        return nil, err
    }
    query := `
    CREATE TABLE IF NOT EXISTS users (
        id INT AUTO_INCREMENT,
        name VARCHAR(255) NOT NULL,
        phone VARCHAR(20) NOT NULL,
        age INT NOT NULL,
        PRIMARY KEY (id)
    )`
    _, err = db.Exec(query)
    if err != nil {
        return nil, err
    }
    _, err = db.Exec("TRUNCATE TABLE users") // Clear any existing data for clean tests
    if err != nil {
        return nil, err
    }
    return db, nil
}

func TestCreateUser(t *testing.T) {
    db, err := setupTestDB()
    if err != nil {
        t.Fatalf("Error setting up test database: %v", err)
    }
    defer db.Close()

    router := mux.NewRouter()
    router.HandleFunc("/users", createUser).Methods("POST")

    user := User{Name: "Ansh", Phone: "708282600", Age: 20}
    jsonUser, _ := json.Marshal(user)
    request, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonUser))
    response := httptest.NewRecorder()
    router.ServeHTTP(response, request)

    if response.Code != http.StatusCreated {
        t.Fatalf("Expected status code %d, but got %d", http.StatusCreated, response.Code)
    }

    var createdUser User
    json.NewDecoder(response.Body).Decode(&createdUser)

    if createdUser.Name != user.Name || createdUser.Phone != user.Phone || createdUser.Age != user.Age {
        t.Fatalf("Expected user %+v, but got %+v", user, createdUser)
    }

    if createdUser.ID == 0 {
        t.Fatalf("Expected user ID to be set, but got %d", createdUser.ID)
    }
}

func TestGetUser(t *testing.T) {
    db, err := setupTestDB()
    if err != nil {
        t.Fatalf("Error setting up test database: %v", err)
    }
    defer db.Close()

    router := mux.NewRouter()
    router.HandleFunc("/users", createUser).Methods("POST")
    router.HandleFunc("/users/{id}", getUser).Methods("GET")

    user := User{Name: "Ansh Malik", Phone: "7082826000", Age: 20}
    jsonUser, _ := json.Marshal(user)
    createRequest, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonUser))
    createResponse := httptest.NewRecorder()
    router.ServeHTTP(createResponse, createRequest)

    var createdUser User
    json.NewDecoder(createResponse.Body).Decode(&createdUser)

    getRequest, _ := http.NewRequest("GET", "/users/"+strconv.Itoa(createdUser.ID), nil)
    getResponse := httptest.NewRecorder()
    router.ServeHTTP(getResponse, getRequest)

    if getResponse.Code != http.StatusOK {
        t.Fatalf("Expected status code %d, but got %d", http.StatusOK, getResponse.Code)
    }

    var fetchedUser User
    json.NewDecoder(getResponse.Body).Decode(&fetchedUser)

    if fetchedUser != createdUser {
        t.Fatalf("Expected user %+v, but got %+v", createdUser, fetchedUser)
    }
}
