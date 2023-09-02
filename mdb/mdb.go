package mdb

import (
	"database/sql"
	"log"
	"time"

	sqlite3 "github.com/mattn/go-sqlite3"
)

type EmailEntry struct {
	Id          int64
	Email       string
	confirmedAt *time.Time
	optOut      bool
}

func TryCreate(db *sql.DB) {
	_, err := db.Exec(`
		create table emails(
			id integer primary key,
			email text unique,
			confirmed_at integer,
			opt_out integer
		)
	`)
	if err != nil {
		if sqlError, ok := err.(sqlite3.Error); ok {
			if sqlError.Code != 1 { // Error Code : 1 ==> Table Exists
				log.Fatal(sqlError)
			}
		} else {
			log.Fatal(err)
		}
	}
}

func getEmailEntryFromRow(row *sql.Rows) (*EmailEntry, error) {
	var id int64
	var email string
	var confirmedAt int64
	var optOut bool

	err := row.Scan(&id, &email, &confirmedAt, &optOut)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	t := time.Unix(confirmedAt, 0)
	return &EmailEntry{Id: id, Email: email, confirmedAt: &t, optOut: optOut}, nil
}

func CreateEmail(db *sql.DB, email string) error {
	_, err := db.Exec(
		`Insert into emails(email,confirmed_at,opt_out) values (?,0,false)`, email)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func GetEmail(db *sql.DB, email string) (*EmailEntry, error) {
	rows, err := db.Query(
		`select id,email,confirmed_at,opt_out from emails where email = ?`, email)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		return getEmailEntryFromRow(rows)
	}
	return nil, nil
}

func UpdateEmail(db *sql.DB, entry EmailEntry) error {
	t := entry.confirmedAt.Unix()

	_, err := db.Exec(`
	Insert into emails(email,confirmed_at,opt_out) values (?,?,?)
	on conflict(email) do update set
	confirmed_at=?
	opt_out=?
	`, entry.Email, t, entry.confirmedAt, entry.optOut, t, entry.optOut)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil

}

func DeleteEmail(db *sql.DB, email string) error {
	_, err := db.Exec(`
	update emails set opt_out = true where email=?
	`, email)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil

}

type BatchEmailQueryParams struct {
	Page  int
	Count int
}

func GetEmailBatch(db *sql.DB, params BatchEmailQueryParams) ([]EmailEntry, error) {
	var empty []EmailEntry
	rows, err := db.Query(`
	select id,email,confirmed_at,opt_out 
	from emails
	where opt_out=false
	order by id ASC
	Limit ? Offset ?
	`, params.Count, (params.Page-1)*params.Count)

	if err != nil {
		log.Println(err)
		return empty, err
	}
	defer rows.Close()
	emails := make([]EmailEntry, 0, params.Count)
	for rows.Next() {
		email, err := getEmailEntryFromRow(rows)
		if err != nil {
			return nil, err
		}
		emails = append(emails, *email)
	}
	return emails, nil
}
