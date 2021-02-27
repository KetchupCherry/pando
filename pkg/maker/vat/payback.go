package vat

import (
	"context"

	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/pkg/maker"
	"github.com/fox-one/pando/pkg/uuid"
	"github.com/shopspring/decimal"
)

func HandlePayback(
	collaterals core.CollateralStore,
	vaults core.VaultStore,
	transactions core.TransactionStore,
	wallets core.WalletStore,
) maker.HandlerFunc {
	frob := HandleFrob(
		collaterals,
		vaults,
		transactions,
		wallets,
	)

	return func(ctx context.Context, r *maker.Request) error {
		var (
			user   uuid.UUID
			follow uuid.UUID
			id     uuid.UUID
		)

		if err := require(r.Scan(&user, &follow, &id) == nil, "bad-data"); err != nil {
			return err
		}

		_, debt := r.Payment()
		r = r.WithBody(user, follow, id, decimal.Zero, debt.Neg())
		return frob(ctx, r)
	}
}
