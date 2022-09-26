//go:build debug
// +build debug

package test

import (
	"SMEI/lib/elevate"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"time"
)

var Cmd = &cobra.Command{
	Use:   "test",
	Short: "Testing the things",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("ayop")
		fmt.Println(os.Args[2])
		if elevate.IsElevated() {
			fmt.Println("sleeping")
			fmt.Println("\a")
			time.Sleep(time.Hour)
		}
		elevate.RerunElevatedFinal()
	},
}
