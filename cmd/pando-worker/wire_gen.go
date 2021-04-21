// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//+build !wireinject

package main

import (
	"github.com/fox-one/pando/cmd/pando-worker/config"
	"github.com/fox-one/pando/handler/node"
	"github.com/fox-one/pando/parliament"
	asset2 "github.com/fox-one/pando/service/asset"
	message2 "github.com/fox-one/pando/service/message"
	oracle2 "github.com/fox-one/pando/service/oracle"
	"github.com/fox-one/pando/service/user"
	wallet2 "github.com/fox-one/pando/service/wallet"
	"github.com/fox-one/pando/store/asset"
	"github.com/fox-one/pando/store/collateral"
	"github.com/fox-one/pando/store/flip"
	"github.com/fox-one/pando/store/message"
	"github.com/fox-one/pando/store/oracle"
	"github.com/fox-one/pando/store/proposal"
	"github.com/fox-one/pando/store/transaction"
	user2 "github.com/fox-one/pando/store/user"
	"github.com/fox-one/pando/store/vault"
	"github.com/fox-one/pando/store/wallet"
	"github.com/fox-one/pando/worker/cashier"
	"github.com/fox-one/pando/worker/events"
	"github.com/fox-one/pando/worker/keeper"
	"github.com/fox-one/pando/worker/messenger"
	"github.com/fox-one/pando/worker/payee"
	"github.com/fox-one/pando/worker/pricesync"
	"github.com/fox-one/pando/worker/spentsync"
	"github.com/fox-one/pando/worker/syncer"
	"github.com/fox-one/pando/worker/txsender"
	"github.com/fox-one/pkg/store/property"
)

// Injectors from wire.go:

func buildApp(cfg *config.Config) (app, error) {
	db, err := provideDatabase(cfg)
	if err != nil {
		return app{}, err
	}
	walletStore := wallet.New(db)
	client, err := provideMixinClient(cfg)
	if err != nil {
		return app{}, err
	}
	walletConfig := provideWalletServiceConfig(cfg)
	walletService := wallet2.New(client, walletConfig)
	system := provideSystem(cfg)
	cashierCashier := cashier.New(walletStore, walletService, system)
	messageStore := message.New(db)
	messageService := message2.New(client)
	messengerMessenger := messenger.New(messageStore, messageService)
	assetStore := asset.New(db)
	assetService := asset2.New(client)
	transactionStore := transaction.New(db)
	proposalStore := proposal.New(db)
	collateralStore := collateral.New(db)
	vaultStore := vault.New(db)
	flipStore := flip.New(db)
	store := propertystore.New(db)
	userConfig := user.Config{}
	userService := user.New(client, userConfig)
	coreParliament := parliament.New(messageStore, userService, assetService, walletService, collateralStore, system)
	oracleStore := oracle.New(db)
	oracleService := oracle2.New(oracleStore)
	payeePayee := payee.New(assetStore, assetService, walletStore, walletService, transactionStore, proposalStore, collateralStore, vaultStore, flipStore, store, coreParliament, oracleStore, oracleService, system)
	sync := pricesync.New(assetStore, assetService)
	userStore := user2.New(db)
	localizer, err := provideLocalizer(cfg)
	if err != nil {
		return app{}, err
	}
	notifier := provideNotifier(system, assetService, messageStore, vaultStore, collateralStore, userStore, localizer)
	spentSync := spentsync.New(walletStore, notifier)
	sender := txsender.New(walletStore)
	syncerSyncer := syncer.New(walletStore, walletService, store)
	eventsEvents := events.New(transactionStore, notifier, store)
	keeperKeeper := keeper.New(collateralStore, oracleStore, vaultStore, walletService, notifier, system)
	v := provideWorkers(cashierCashier, messengerMessenger, payeePayee, sync, spentSync, sender, syncerSyncer, eventsEvents, keeperKeeper)
	server := node.New(system, store, oracleStore)
	mux := provideRoute(server)
	serverServer := provideServer(mux)
	mainApp := app{
		workers: v,
		server:  serverServer,
	}
	return mainApp, nil
}
