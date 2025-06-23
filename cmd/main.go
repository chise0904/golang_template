package main

import (
	"github.com/chise0904/golang_template/cmd/server"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := cobra.Command{}

	rootCmd.AddCommand(server.ServerCmd())

	rootCmd.Execute()
}
