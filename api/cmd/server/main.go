package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/gorilla/mux"
	validator9 "gopkg.in/go-playground/validator.v9"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/gocql/gocql"
	translationsEn "gopkg.in/go-playground/validator.v9/translations/en"

	"github.com/ustrugany/projectx/api"
	"github.com/ustrugany/projectx/api/infrastructure/delivery/stdout"
	"github.com/ustrugany/projectx/api/infrastructure/persistence/cassandra"
	"github.com/ustrugany/projectx/api/infrastructure/validation"
	httpInterface "github.com/ustrugany/projectx/api/interfaces/http"
	"github.com/ustrugany/projectx/api/service"
)

func main() {
	// Config
	config, err := api.CreateConfig()
	if err != nil {
		panic(err)
	}

	// Logger
	baseLogger := CreateLogger()
	defer func() {
		_ = baseLogger.Sync()
	}()
	logger := *baseLogger.Sugar()

	logger.Debugw("@TODO remove me", "config", config)

	cluster := gocql.NewCluster(config.Cassandra.Host)
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: config.Cassandra.User,
		Password: config.Cassandra.Password,
	}
	cluster.Keyspace = config.Cassandra.Keyspace
	Session, err := cluster.CreateSession()
	if err != nil {
		panic(err)
	}

	defer Session.Close()

	// Initialize dependencies
	//messageRepository := inmemory.CreateMessageRepository(
	//	inmemory.CreateDb(),
	//)
	messageRepository := cassandra.CreateMessageRepository(
		Session,
	)

	// Validator & translator setup
	locale := "en"
	validator := validator9.New()
	translatorEn := en.New()
	uni := ut.New(translatorEn, translatorEn)
	translator, found := uni.GetTranslator(locale)
	if !found {
		logger.Fatalw("required translator not found", "locale", locale)
	}
	if err = translationsEn.RegisterDefaultTranslations(validator, translator); err != nil {
		logger.Fatal(err)
	}
	validator.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	// Create message endpoint
	messageValidator := validation.CreateMessageValidator(*validator, translator)
	createMessageService := service.CreateCreateMessage(messageRepository, messageValidator)
	createMessageHandler := httpInterface.CreateCreateMessageHandler(
		createMessageService,
		logger,
		config,
	)

	// List messages endpoint
	listMessagesService := service.CreateListMessages(
		messageRepository,
	)
	getMessageHandler := httpInterface.CreateListMessagesHandler(
		listMessagesService,
		logger,
		config,
	)

	// Send message endpoint
	delivery := stdout.CreateMessageDelivery(logger)
	sendMessageService := service.CreateSendMessage(messageRepository, delivery)
	updateMessageHandler := httpInterface.CreateSendMessageHandler(
		sendMessageService,
		logger,
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
	listMessages http.Handler,
	sendMessage http.Handler,
	createMessage http.Handler,
	logger zap.SugaredLogger,
) *cobra.Command {
	return &cobra.Command{
		Use:   "server",
		Short: "Runs the API server",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			r := mux.NewRouter()
			subRouter := r.PathPrefix("/api").
				Subrouter()
			subRouter.Handle("/message", createMessage).Methods(http.MethodPost)
			subRouter.Handle("/send", sendMessage).Methods(http.MethodPost)
			subRouter.Handle("/messages/{email}", listMessages).
				Methods(http.MethodGet)
			logger.Debugw("listening...", "config", config)
			if err := http.ListenAndServe(fmt.Sprintf(":%d", config.Port), r); err != nil {
				panic(err)
			}
		},
	}
}
