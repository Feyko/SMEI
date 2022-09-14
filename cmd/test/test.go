//go:build debug
// +build debug

package test

import (
	"SMEI/lib/crypt"
	"fmt"
	"github.com/spf13/cobra"
	"log"
)

var Cmd = &cobra.Command{
	Use:   "test",
	Short: "Testing the things",
	Run: func(cmd *cobra.Command, args []string) {
		password := "ayoo"
		secure, err := crypt.Encrypt(password, "heyadgffgddfhhdfhfdhfdhdfhdfdfhdfhdfhhdfdhfhdfhdfhdfhfdhdf")
		if err != nil {
			log.Fatalf("Encryption error: %v", err)
		}
		fmt.Println(secure)
		unsecure, err := crypt.Decrypt(password, secure)
		if err != nil {
			log.Fatalf("Decryption error: %v", err)
		}
		fmt.Println(unsecure)
	},
}
