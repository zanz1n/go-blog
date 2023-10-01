//go:build !production
// +build !production

package cmd

import (
	"log"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalln(err)
	}
}
