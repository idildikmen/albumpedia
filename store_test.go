package main

import (
	"database/sql"
	"testing"

	// The "testify/suite" package is used to make the test suite
	"github.com/stretchr/testify/suite"
)

type StoreSuite struct {
	suite.Suite
	/*
		The suite is defined as a struct, with the store and db as its
		attributes. Any variables that are to be shared between tests in a
		suite should be stored as attributes of the suite instance
	*/
	store *dbStore
	db    *sql.DB
}

func (s *StoreSuite) SetupSuite() {
	/*
		The database connection is opened in the setup, and
		stored as an instance variable,
		as is the higher level `store`, that wraps the `db`
	*/
	var testalbum Albums

	connString := "sqlite-database-album-test.db"
	db, err := sql.Open("sqlite3", connString)
	if err != nil {
		s.T().Fatal(err)
	}

	db.Exec("DROP TABLE IF EXISTS albums")
	db.Exec("DROP TABLE OIF EXISTS sqlite_sequence")
	createTestTable(db)

	for i := 0; i < len(testalbum.Albums); i++ {
		album := testalbum.Albums[i]
		_, err1 := db.Exec("INSERT INTO testalbum(title, artist, class, genre, year) VALUES ($1,$2,$3,$4,$5)", album.Title, album.Artist,
			album.Class, album.Genre, album.Year)
		_ = err1
	}

	s.db = db
	s.store = &dbStore{db: db}
}

func (s *StoreSuite) SetupTest() {
	/*
		We delete all entries from the table before each test runs, to ensure a
		consistent state before our tests run. In more complex applications, this
		is sometimes achieved in the form of migrations
	*/
	_, err := s.db.Exec("DELETE FROM testalbum")
	//defer s.db.Close() // Defer Closing the database

	if err != nil {
		s.T().Fatal(err)
	}
}

func (s *StoreSuite) TearDownSuite() {
	// Close the connection after all tests in the suite finish
	s.db.Close()
}

// This is the actual "test" as seen by Go, which runs the tests defined below
func TestStoreSuite(t *testing.T) {

	s := new(StoreSuite)
	suite.Run(t, s)
}

func (s *StoreSuite) TestCreateBird() {
	// Create a bird through the store `CreateBird` method
	s.store.CreateTestAlbum(&Album{
		Artist: "Beyonce",
		Title:  "Halo",
		Price:  "33.99",
	})

	// Query the database for the entry we just created
	res, err := s.db.Query(`SELECT COUNT(*) FROM testalbum WHERE artist='Beyonce' AND title='Halo' AND price='33.99'`)
	if err != nil {
		s.T().Fatal(err)
	}

	// Get the count result
	var count int

	for res.Next() {
		err := res.Scan(&count)
		if err != nil {
			s.T().Error(err)
		}
	}

	// Assert that there must be one entry with the properties of the bird that we just inserted (since the database was empty before this)
	if count != 1 {
		s.T().Errorf("incorrect count, wanted 1, got %d", count)
	}
}

func (s *StoreSuite) TestGetBird() {
	// Insert a sample bird into the `birds` table
	_, err := s.db.Exec(`INSERT INTO testalbum (title, artist,price) VALUES('Halo','Beyonce','33.99')`)
	if err != nil {
		s.T().Fatal(err)
	}

	// Get the list of birds through the stores `GetBirds` method
	testalbum, err := s.store.GetTestAlbums()
	if err != nil {
		s.T().Fatal(err)
	}

	// Assert that the count of birds received must be 1
	nAlbums := len(testalbum)
	if nAlbums != 1 {
		s.T().Errorf("incorrect count, wanted 1, got %d", nAlbums)
	}

	// Assert that the details of the bird is the same as the one we inserted
	expectedAlbum := Album{"Halo", "Beyonce", "33.99"}
	if *testalbum[0] != expectedAlbum {
		s.T().Errorf("incorrect details, expected %v, got %v", expectedAlbum, *testalbum[0])
	}
}
