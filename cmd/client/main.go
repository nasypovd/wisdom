package main

import (
	"fmt"
	"log"
	"net"
	"time"
	"wisdom/pkg/client"
	"wisdom/pkg/pow"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var rootCmd = &cobra.Command{
	Use:   "client",
	Short: "This is the wisdom client",
	Long:  `This client sends requests to a server to get wisdom quotes after solving proof of work challenges.`,
	RunE:  runClient,
}

func main() {
	rootCmd.PersistentFlags().String("serverAddr", "localhost:8080", "Address of the server to connect to")
	rootCmd.PersistentFlags().Duration("timeout", 5*time.Second, "Timeout for the client to connect to the server")

	viper.BindPFlag("serverAddr", rootCmd.PersistentFlags().Lookup("serverAddr"))
	viper.BindPFlag("timeout", rootCmd.PersistentFlags().Lookup("timeout"))

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func runClient(cmd *cobra.Command, args []string) error {
	log, err := zap.NewDevelopment()
	if err != nil {
		return fmt.Errorf("failed to create logger: %w", err)
	}

	serverAddr := viper.GetString("serverAddr")
	timeout := viper.GetDuration("timeout")

	conn, err := net.DialTimeout("tcp", serverAddr, timeout)
	if err != nil {
		return fmt.Errorf("failed to connect to the server: %w", err)
	}

	defer conn.Close()

	c := client.NewClient(conn, pow.NewSolver(), log.With(zap.String("serverAddr", serverAddr)))

	quote, err := c.Run()
	if err != nil {
		return fmt.Errorf("failed to run client: %w", err)
	}

	fmt.Println("QUOTE:", quote)

	return nil
}
