package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/ustrugany/projectx/api"
	"github.com/ustrugany/projectx/api/infrastructure/persistence/inmemory"
	httpInterface "github.com/ustrugany/projectx/api/interfaces/http"
	"github.com/ustrugany/projectx/api/service"
)

func main() {
	// Config
	config, err := api.NewConfig()
	if err != nil {
		panic(err)
	}

	// Logger
	baseLogger := CreateLogger()
	defer func() {
		_ = baseLogger.Sync()
	}()
	logger := baseLogger.Sugar()

	logger.Infow("@TODO remove", "config", config)

	// Initialize dependencies
	messageRepository := inmemory.CreateMessageRepository(
		inmemory.CreateDb(),
	)

	createMessageUseCase := service.CreateCreateMessageUseCase(messageRepository)
	createMessageHandler := httpInterface.CreateCreateMessageHandler(
		createMessageUseCase,
		*logger,
		config,
	)

	listMessagesUseCase := service.CreateListMessagesUseCase(messageRepository)
	getMessageHandler := httpInterface.CreateListMessagesHandler(
		listMessagesUseCase,
		*logger,
		config,
	)

	sendMessageUseCase := service.CreateSendMessageUseCase(messageRepository)
	updateMessageHandler := httpInterface.CreateSendMessageHandler(
		sendMessageUseCase,
		*logger,
		config,
	)

	// CLI initialization
	serverCommand := CreateServerCommand(
		&config,
		getMessageHandler,
		updateMessageHandler,
		createMessageHandler,
		logger,
	)
	serverCommand.Flags().IntVarP(&config.Port, "port", "p", 0, "port (required)")
	if err = serverCommand.MarkFlagRequired("port"); err != nil {
		panic(err)
	}

	// Execute root command
	rootCmd := &cobra.Command{
		Use: "projectx",
	}
	rootCmd.AddCommand(serverCommand)
	if err = rootCmd.Execute(); err != nil {
		panic(err)
	}
}

func CreateLogger() *zap.Logger {
	rawJSON := []byte(`{
	  "level": "debug",
	  "encoding": "json",
	  "outputPaths": ["stdout", "/tmp/logs"],
	  "errorOutputPaths": ["stderr"]
	}`)
	var cfg zap.Config
	if err := json.Unmarshal(rawJSON, &cfg); err != nil {
		panic(err)
	}
	cfg.EncoderConfig = zap.NewProductionEncoderConfig()
	baseLogger, err := cfg.Build()
	if err != nil {
		panic(err)
	}

	return baseLogger
}

func CreateServerCommand(
	config *api.Config,
	gmh httpInterface.ListMessagesHandler,
	umh httpInterface.SendMessageHandler,
	cmh httpInterface.CreateMessageHandler,
	logger *zap.SugaredLogger,
) *cobra.Command {
	return &cobra.Command{
		Use:   "server",
		Short: "Runs the API server",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			r := mux.NewRouter()
			subRouter := r.PathPrefix("/api").Subrouter()
			subRouter.Handle("/message", cmh).Methods(http.MethodPost)
			subRouter.Handle("/message/send", umh).Methods(http.MethodPost)
			subRouter.Handle("/messages/{email}", gmh).Methods(http.MethodGet)
			logger.Infow("Listening...",
				"config", config,
			)
			if err := http.ListenAndServe(fmt.Sprintf(":%d", config.Port), r); err != nil {
				panic(err)
			}
		},
	}
}
