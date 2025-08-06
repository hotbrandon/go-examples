package main

import (
	"log"
	"oracle-demo/handlers"
	"oracle-demo/models"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/sijms/go-ora/v2"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("cannot find .env file: %v", err)
	}

	oracle_dsn_lp := os.Getenv("ORACLE_DSN_LP")
	if oracle_dsn_lp == "" {
		log.Fatal("ORACLE_DSN_LP is not set")
	}
	models.OracleConfigs["LP"] = models.OracleConfig{
		DSN: oracle_dsn_lp,
	}

	oracle_dsn_nd := os.Getenv("ORACLE_DSN_ND")
	if oracle_dsn_nd == "" {
		log.Fatal("ORACLE_DSN_ND is not set")
	}
	models.OracleConfigs["ND"] = models.OracleConfig{
		DSN: oracle_dsn_nd,
	}

	// for key, dsn := range models.OracleConfigs {
	// 	log.Printf("Oracle DSN for %s: %s", key, dsn.DSN)
	// }

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	r.GET("/invoice/gen_c0401/:segment_no", handlers.GenerateC0401Handler)

	r.Run("0.0.0.0:8080")
}
