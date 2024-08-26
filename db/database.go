package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// Connect handles connection to PostgreSQL, database creation, and table creation
func Connect() (*sql.DB, error) {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// dbHost := os.Getenv("DB_HOST")
	// dbPort := os.Getenv("DB_PORT")
	// dbUser := os.Getenv("DB_USER")
	// dbPassword := os.Getenv("DB_PASSWORD")
	// dbName := os.Getenv("DB_NAME")

	// // Construct the connection string
	// conntr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
	// 	dbHost, dbPort, dbUser, dbPassword, dbName)
	// productio connection
	connStr := os.Getenv("db_url")

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	// Create database if it does not exist

	//defer db.Close()

	tableSchemas := []string{
		// 		`GRANT ALL PRIVILEGES ON DATABASE postgres TO vehiclemandb_user;
		// `,

		`CREATE SEQUENCE IF NOT EXISTS users_user_id_seq
`,
		`CREATE SEQUENCE IF NOT EXISTS  subscriptions_subscription_id_seq`,
		`CREATE SEQUENCE IF NOT EXISTS  expenses_expense_id_seq`,
		`CREATE SEQUENCE IF NOT EXISTS incomes_income_id_seq`,
		`CREATE SEQUENCE IF NOT EXISTS lorries_lorry_id_seq`,

		`CREATE TABLE IF NOT EXISTS users
(
    user_id integer NOT NULL DEFAULT nextval('users_user_id_seq'::regclass),
    password character varying(255) COLLATE pg_catalog."default" NOT NULL,
    email character varying(100) COLLATE pg_catalog."default" NOT NULL,
    role character varying(20) COLLATE pg_catalog."default" NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    name character varying(256) COLLATE pg_catalog."default",
    CONSTRAINT users_pkey PRIMARY KEY (user_id),
    CONSTRAINT users_email_key UNIQUE (email)
)`,

		`CREATE TABLE IF NOT EXISTS user_profiles
(
    user_id integer NOT NULL,
    first_name character varying(50) COLLATE pg_catalog."default",
    last_name character varying(50) COLLATE pg_catalog."default",
    phone_number character varying(20) COLLATE pg_catalog."default",
    address text COLLATE pg_catalog."default",
    CONSTRAINT user_profiles_pkey PRIMARY KEY (user_id),
    CONSTRAINT user_profiles_user_id_fkey FOREIGN KEY (user_id)
        REFERENCES users (user_id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
)`,

		`CREATE TABLE IF NOT EXISTS vehicles
(
    vehicle_id integer NOT NULL DEFAULT nextval('lorries_lorry_id_seq'::regclass),
    make character varying(50) COLLATE pg_catalog."default" NOT NULL,
    model character varying(50) COLLATE pg_catalog."default" NOT NULL,
    year integer NOT NULL,
    registration_number character varying(20) COLLATE pg_catalog."default" NOT NULL,
    capacity integer NOT NULL,
    owner_id integer NOT NULL,
    CONSTRAINT lorries_pkey PRIMARY KEY (vehicle_id),
    CONSTRAINT lorries_registration_number_key UNIQUE (registration_number),
    CONSTRAINT fk_owner FOREIGN KEY (owner_id)
        REFERENCES users (user_id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
)
`,

		`CREATE TABLE IF NOT EXISTS subscriptions
(
    subscription_id integer NOT NULL DEFAULT nextval('subscriptions_subscription_id_seq'::regclass),
    user_id integer NOT NULL,
    plan_name character varying(50) COLLATE pg_catalog."default" NOT NULL,
    start_date date NOT NULL,
    end_date date NOT NULL,
    status character varying(20) COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT subscriptions_pkey PRIMARY KEY (subscription_id),
    CONSTRAINT subscriptions_user_id_fkey FOREIGN KEY (user_id)
        REFERENCES users (user_id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
)
`,
		`CREATE TABLE IF NOT EXISTS Products (
			ID SERIAL PRIMARY KEY,
			Name VARCHAR(100) NOT NULL,
			Description TEXT,
			Price DECIMAL(10, 2) NOT NULL,
			CreatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS expense_categories (
			category_id SERIAL PRIMARY KEY,
			category_name VARCHAR(50) NOT NULL
		)`,

		`CREATE TABLE IF NOT EXISTS expenses
		(
			expense_id integer NOT NULL DEFAULT nextval('expenses_expense_id_seq'::regclass),
			vehicle_id integer NOT NULL,
			category_id integer NOT NULL,
			amount numeric(10,2) NOT NULL,
			description text COLLATE pg_catalog."default",
			receipt text COLLATE pg_catalog."default",
			expense_date date NOT NULL,
			CONSTRAINT expenses_pkey PRIMARY KEY (expense_id),
			CONSTRAINT expenses_category_id_fkey FOREIGN KEY (category_id)
				REFERENCES expense_categories (category_id) MATCH SIMPLE
				ON UPDATE NO ACTION
				ON DELETE CASCADE,
			CONSTRAINT expenses_vehicle_id_fkey FOREIGN KEY (vehicle_id)
				REFERENCES vehicles (vehicle_id) MATCH SIMPLE
				ON UPDATE NO ACTION
				ON DELETE CASCADE
		)`,

		`CREATE TABLE IF NOT EXISTS incomes
(
    income_id integer NOT NULL DEFAULT nextval('incomes_income_id_seq'::regclass),
    vehicle_id integer NOT NULL,
    amount numeric(10,2) NOT NULL,
    payment_date date NOT NULL,
    status character varying(20) COLLATE pg_catalog."default" NOT NULL,
    description character varying(255) COLLATE pg_catalog."default",
    CONSTRAINT incomes_pkey PRIMARY KEY (income_id)
)`,

		// Add more table schemas as needed
	}

	for _, schema := range tableSchemas {
		_, err = db.Exec(schema)
		if err != nil {
			return nil, fmt.Errorf("error creating table: %v", err)
		}
	}

	return db, nil
}

// ---------------Sqllit beging ---------------
// import (
// 	"database/sql"
// 	"fmt"
// 	"log"

// 	"github.com/joho/godotenv"
// 	_ "github.com/mattn/go-sqlite3"
// )

// // Connect handles connection to SQLite, database creation, and table creation
// func Connect() (*sql.DB, error) {

// 	err := godotenv.Load()
// 	if err != nil {
// 		log.Fatalf("Error loading .env file: %v", err)
// 	}

// 	dbName := "vehicle.db" // SQLite uses file-based databases

// 	// Construct the connection string for SQLite
// 	connStr := fmt.Sprintf("file:%s?cache=shared&mode=rwc", dbName)

// 	db, err := sql.Open("sqlite3", connStr)
// 	if err != nil {
// 		return nil, err
// 	}

// 	tableSchemas := []string{
// 		`CREATE TABLE IF NOT EXISTS users (
// 			user_id INTEGER PRIMARY KEY AUTOINCREMENT,
// 			password TEXT NOT NULL,
// 			email TEXT NOT NULL UNIQUE,
// 			role TEXT NOT NULL,
// 			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
// 			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
// 			name TEXT
// 		)`,

// 		`CREATE TABLE IF NOT EXISTS user_profiles (
// 			user_id INTEGER PRIMARY KEY,
// 			first_name TEXT,
// 			last_name TEXT,
// 			phone_number TEXT,
// 			address TEXT,
// 			FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
// 		)`,

// 		`CREATE TABLE IF NOT EXISTS vehicles (
// 			vehicle_id INTEGER PRIMARY KEY AUTOINCREMENT,
// 			make TEXT NOT NULL,
// 			model TEXT NOT NULL,
// 			year INTEGER NOT NULL,
// 			registration_number TEXT NOT NULL UNIQUE,
// 			capacity INTEGER NOT NULL,
// 			owner_id INTEGER NOT NULL,
// 			FOREIGN KEY (owner_id) REFERENCES users(user_id)
// 		)`,

// 		`CREATE TABLE IF NOT EXISTS subscriptions (
// 			subscription_id INTEGER PRIMARY KEY AUTOINCREMENT,
// 			user_id INTEGER NOT NULL,
// 			plan_name TEXT NOT NULL,
// 			start_date DATE NOT NULL,
// 			end_date DATE NOT NULL,
// 			status TEXT NOT NULL,
// 			FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
// 		)`,

// 		`CREATE TABLE IF NOT EXISTS Products (
// 			ID INTEGER PRIMARY KEY AUTOINCREMENT,
// 			Name TEXT NOT NULL,
// 			Description TEXT,
// 			Price DECIMAL(10, 2) NOT NULL,
// 			CreatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP
// 		)`,

// 		`CREATE TABLE IF NOT EXISTS expense_categories (
// 			category_id INTEGER PRIMARY KEY AUTOINCREMENT,
// 			category_name TEXT NOT NULL
// 		)`,

// 		`CREATE TABLE IF NOT EXISTS expenses (
// 			expense_id INTEGER PRIMARY KEY AUTOINCREMENT,
// 			vehicle_id INTEGER NOT NULL,
// 			category_id INTEGER NOT NULL,
// 			amount DECIMAL(10, 2) NOT NULL,
// 			description TEXT,
// 			receipt TEXT,
// 			expense_date DATE NOT NULL,
// 			FOREIGN KEY (category_id) REFERENCES expense_categories(category_id) ON DELETE CASCADE,
// 			FOREIGN KEY (vehicle_id) REFERENCES vehicles(vehicle_id) ON DELETE CASCADE
// 		)`,

// 		`CREATE TABLE IF NOT EXISTS incomes (
// 			income_id INTEGER PRIMARY KEY AUTOINCREMENT,
// 			vehicle_id INTEGER NOT NULL,
// 			amount DECIMAL(10, 2) NOT NULL,
// 			payment_date DATE NOT NULL,
// 			status TEXT NOT NULL,
// 			description TEXT,
// 			FOREIGN KEY (vehicle_id) REFERENCES vehicles(vehicle_id)
// 		)`,
// 	}

// 	for _, schema := range tableSchemas {
// 		_, err = db.Exec(schema)
// 		if err != nil {
// 			return nil, fmt.Errorf("error creating table: %v", err)
// 		}
// 	}

// 	return db, nil
// }
