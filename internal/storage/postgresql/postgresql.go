package postgresql

import (
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/lakshsetia/learn-RESTAPI/internal/config"
	"github.com/lakshsetia/learn-RESTAPI/internal/types"
	_ "github.com/lib/pq"
)

type Postgresql struct {
	DB *sql.DB
}
func New(config *config.Config) (*Postgresql, error) {
	user, password, dbname, host, port := config.Database.Postgresql.User, config.Database.Postgresql.Password, config.Database.Postgresql.DBName, config.Database.Postgresql.Host, config.Database.Postgresql.Port 
	if user == "" || password == "" || dbname == "" || host == "" || port == "" {
		return nil, fmt.Errorf("postgresql connection parameters not specified")
	}
	connectionStr := fmt.Sprintf("user=%v password=%v dbname=%v host=%v port=%v sslmode=disable", user, password, dbname, host, port)
	slog.Info(connectionStr)
	db, err := sql.Open("postgres", connectionStr)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
	id BIGSERIAL PRIMARY KEY,
	name VARCHAR(100) NOT NULL,
	email VARCHAR(150) NOT NULL,
	age INTEGER CHECK (age >= 0) 	
	)`)
	if err != nil {
		return nil, err
	}
	return &Postgresql{
		DB: db,
	}, nil
}
func (p *Postgresql) GetUsers() ([]types.User, error) {
	rows, err := p.DB.Query("SELECT id, name, email, age FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	users := make([]types.User, 0)
	for rows.Next() {
		var user types.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Age); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}
func (p *Postgresql) CreateUser(name string, email string, age int) error {
	if _, err := p.DB.Exec("INSERT INTO users (name, email, age) VALUES ($1, $2, $3)", name, email, age); err != nil {
		return err
	}
	return nil
}
func (p *Postgresql) GetUserById(id int) (types.User, error) {
	var user types.User
	err := p.DB.QueryRow("SELECT id, name, email, age FROM users WHERE id=$1", id).Scan(&user.ID, &user.Name, &user.Email, &user.Age)
	if err == sql.ErrNoRows {
		return user, fmt.Errorf("no user found with id=%v", id)
	}
	if err != nil {
		return user, err
	}
	return user, nil
}
func (p *Postgresql) UpdateUserById(id int, name string, email string, age int) error {
	if _, err := p.DB.Exec("UPDATE users SET name=$1, email=$2, age=$3 WHERE id=$4", name, email, age, id); err != nil {
		return err
	}
	return nil
}
func (p *Postgresql) DeleteUserById(id int) error {
	if _, err := p.DB.Exec("DELETE FROM users WHERE id=$1", id); err != nil {
		return err
	}
	return nil
}