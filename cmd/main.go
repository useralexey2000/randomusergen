package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"randomusergen/apiclient"
	"randomusergen/repo/pg"
	"randomusergen/usergen"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	host       = getEnv("SERVER_HOST", "localhost")
	port       = getEnv("SERVER_PORT", "8080")
	dbhost     = getEnv("POSTGRES_HOST", "localhost")
	dbport     = getEnv("POSTGRES_PORT", "5432")
	dbuser     = getEnv("POSTGRES_USER", "postgres")
	dbpassword = getEnv("POSTGRES_PASSWORD", "postgres")
	dbname     = getEnv("POSTGRES_DB", "postgres")

	serverMode = getEnv("GIN_MODE", "release")
)

func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	return value
}

func main() {
	log.Printf("server started %v\n", net.JoinHostPort(host, port))
	defer func() {
		log.Println("server exited")
	}()

	db, err := pg.New(
		fmt.Sprintf("postgresql://%v:%v@%v:%v/%v",
			dbuser, dbpassword, dbhost, dbport, dbname),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		db.Close()
		log.Println("db stopped")
	}()

	apic := apiclient.New(http.DefaultClient)

	us := usergen.New(apic, db)

	gin.SetMode(serverMode)
	router := gin.New()
	v1 := router.Group("/api/v1")
	{
		v1.POST("/create", us.CreateUsers())
	}

	srv := http.Server{
		Addr:    net.JoinHostPort(host, port),
		Handler: router,
	}

	go func() {
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
}
