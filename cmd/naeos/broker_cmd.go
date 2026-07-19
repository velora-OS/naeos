package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/NAEOS-foundation/naeos/internal/broker"
)

func newBrokerCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "broker",
		Short: "Message broker management",
		Long: `Manage message broker connections (Redis, RabbitMQ, Kafka).

Example:
  naeos broker connect --type redis --name myredis --host localhost --port 6379
  naeos broker list
  naeos broker publish --name myredis --channel events --message '{"event":"created"}'
  naeos broker subscribe --name myredis --channel events
  naeos broker disconnect --name myredis`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(newBrokerConnectCommand())
	cmd.AddCommand(newBrokerListCommand())
	cmd.AddCommand(newBrokerPublishCommand())
	cmd.AddCommand(newBrokerDisconnectCommand())

	return cmd
}

func newBrokerConnectCommand() *cobra.Command {
	var brokerType, name, host, password string
	var port, db int

	cmd := &cobra.Command{
		Use:   "connect",
		Short: "Connect to a message broker",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			mgr := broker.NewManager()

			var b broker.Broker
			switch brokerType {
			case "redis":
				b = broker.NewRedis()
			case "rabbitmq":
				b = broker.NewRabbitMQ()
			case "kafka":
				b = broker.NewKafka()
			default:
				return fmt.Errorf("unsupported broker type: %s", brokerType)
			}

			mgr.Register(name, b)

			cfg := &broker.Config{
				Host:     host,
				Port:     port,
				Password: password,
				DB:       db,
			}

			if err := b.Connect(cfg); err != nil {
				return fmt.Errorf("failed to connect: %w", err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Connected to %s broker '%s'\n", brokerType, name)
			return nil
		},
	}

	cmd.Flags().StringVar(&brokerType, "type", "redis", "broker type (redis, rabbitmq, kafka)")
	cmd.Flags().StringVar(&name, "name", "", "connection name (required)")
	cmd.Flags().StringVar(&host, "host", "localhost", "broker host")
	cmd.Flags().IntVar(&port, "port", 6379, "broker port")
	cmd.Flags().StringVar(&password, "password", "", "broker password")
	cmd.Flags().IntVar(&db, "db", 0, "Redis database number")
	_ = cmd.MarkFlagRequired("name")
	return cmd
}

func newBrokerListCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all broker connections",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			mgr := broker.NewManager()

			names := mgr.List()
			if len(names) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "No broker connections.")
				return nil
			}

			out := cmd.OutOrStdout()
			fmt.Fprintf(out, "%-15s\n", "NAME")
			fmt.Fprintf(out, "%-15s\n", "----")
			for _, name := range names {
				fmt.Fprintf(out, "%-15s\n", name)
			}
			return nil
		},
	}
}

func newBrokerPublishCommand() *cobra.Command {
	var name, channel, message string

	cmd := &cobra.Command{
		Use:   "publish",
		Short: "Publish a message to a channel",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			mgr := broker.NewManager()

			b, ok := mgr.Get(name)
			if !ok {
				return fmt.Errorf("broker '%s' not found", name)
			}

			msg := &broker.Message{
				Payload: []byte(message),
			}

			if err := b.Publish(channel, msg); err != nil {
				return fmt.Errorf("failed to publish: %w", err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Published to '%s' on '%s'\n", channel, name)
			return nil
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "broker name (required)")
	cmd.Flags().StringVar(&channel, "channel", "", "channel name (required)")
	cmd.Flags().StringVar(&message, "message", "", "message payload (required)")
	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("channel")
	_ = cmd.MarkFlagRequired("message")
	return cmd
}

func newBrokerDisconnectCommand() *cobra.Command {
	var name string

	cmd := &cobra.Command{
		Use:   "disconnect",
		Short: "Disconnect from a broker",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			mgr := broker.NewManager()

			b, ok := mgr.Get(name)
			if !ok {
				return fmt.Errorf("broker '%s' not found", name)
			}

			if err := b.Close(); err != nil {
				return fmt.Errorf("failed to disconnect: %w", err)
			}

			mgr.Remove(name)
			fmt.Fprintf(cmd.OutOrStdout(), "Disconnected from '%s'.\n", name)
			return nil
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "broker name (required)")
	_ = cmd.MarkFlagRequired("name")
	return cmd
}
