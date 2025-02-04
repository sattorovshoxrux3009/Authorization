package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"GitHub.com/sattorovshohruh3009/Authorization/config"
	"GitHub.com/sattorovshohruh3009/Authorization/storage"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	cfg := config.Load(".")
	mysqlUrl := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		cfg.Mysql.User,
		cfg.Mysql.Password,
		cfg.Mysql.Host,
		cfg.Mysql.Port,
		cfg.Mysql.Database,
	)
	db, err := sql.Open("mysql", mysqlUrl)
	if err != nil {
		log.Fatal("Bazaga ulanishda xatolik:", err)
	}
	defer db.Close()

	// Ulanishni tekshirish
	err = db.Ping()
	if err != nil {
		log.Fatal("Connection error:", err)
	}
	fmt.Println("Connection succs!")

	strg := storage.NewStorage(db)
	user_subject, err := strg.User().UpdatePassword(context.TODO(), "shohruh", "Shohruh2004")
	if err != nil {
		log.Fatal("error::::", err)
	}
	fmt.Println(user_subject)

}
