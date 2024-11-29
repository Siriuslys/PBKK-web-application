package main

import (
    "database/sql"
    "fmt"
    "log"
    "os"

    _ "github.com/go-sql-driver/mysql" // Import the MySQL driver
)

// Album represents a record in the album table
type Album struct {
    ID     int64
    Title  string
    Artist string
    Price  float32
}

var db *sql.DB

func main() {
    // Capture connection properties from environment variables
    dbUser := os.Getenv("DBUSER")
    dbPass := os.Getenv("DBPASS")
    dbHost := os.Getenv("DBHOST") // For example, "127.0.0.1"
    dbName := os.Getenv("DBNAME")

    // Build the data source name (DSN)
    dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=true&loc=Asia%%2FJakarta",
        dbUser, dbPass, dbHost, dbName)

    // Initialize the database connection
    var err error
    db, err = sql.Open("mysql", dsn)
    if err != nil {
        log.Fatalf("Failed to open database: %v", err)
    }

    // Verify the database connection
    if err := db.Ping(); err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }
    fmt.Println("Connected to the database!")

    // Query albums by artist
    albums, err := albumsByArtist("John Coltrane")
    if err != nil {
        log.Fatalf("Error fetching albums: %v", err)
    }
    fmt.Printf("Albums found: %v\n", albums)

    // Query album by ID
    alb, err := albumByID(2)
    if err != nil {
        log.Fatalf("Error fetching album: %v", err)
    }
    fmt.Printf("Album found: %v\n", alb)

    // Add a new album
    albID, err := addAlbum(Album{
        Title:  "The Modern Sound of Betty Carter",
        Artist: "Betty Carter",
        Price:  49.99,
    })
    if err != nil {
        log.Fatalf("Error adding album: %v", err)
    }
    fmt.Printf("ID of added album: %v\n", albID)
}

// albumsByArtist queries for albums that have the specified artist name.
func albumsByArtist(name string) ([]Album, error) {
    // An albums slice to hold data from returned rows.
    var albums []Album

    rows, err := db.Query("SELECT id, title, artist, price FROM album WHERE artist = ?", name)
    if err != nil {
        return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
    }
    defer rows.Close()

    // Loop through rows, using Scan to assign column data to struct fields.
    for rows.Next() {
        var alb Album
        if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
            return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
        }
        albums = append(albums, alb)
    }
    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
    }
    return albums, nil
}

// albumByID queries for the album with the specified ID.
func albumByID(id int64) (Album, error) {
    // An album to hold data from the returned row.
    var alb Album

    row := db.QueryRow("SELECT id, title, artist, price FROM album WHERE id = ?", id)
    if err := row.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
        if err == sql.ErrNoRows {
            return alb, fmt.Errorf("albumByID %d: no such album", id)
        }
        return alb, fmt.Errorf("albumByID %d: %v", id, err)
    }
    return alb, nil
}

// addAlbum adds a new album to the database.
func addAlbum(alb Album) (int64, error) {
    result, err := db.Exec("INSERT INTO album (title, artist, price) VALUES (?, ?, ?)", alb.Title, alb.Artist, alb.Price)
    if err != nil {
        return 0, fmt.Errorf("addAlbum: %v", err)
    }
    id, err := result.LastInsertId()
    if err != nil {
        return 0, fmt.Errorf("addAlbum: %v", err)
    }
    return id, nil
}
