package cat

import (
	"github.com/fox-one/pando/cmd/pando-cli/cmds/actions"
	"github.com/fox-one/pando/cmd/pando-cli/cmds/pay"
	"github.com/fox-one/pando/cmd/pando-cli/internal/call"
	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/handler/rpc/api"
	"github.com/fox-one/pando/pkg/mtg/types"
	"github.com/fox-one/pando/pkg/number"
	"github.com/spf13/cobra"
)

func NewSupplyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "supply <collateral id> <amount>",
		Short: "execute Supply action on a collateral",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]
			amount := number.Decimal(args[1])

			cat, err := call.RPC().FindCollateral(cmd.Context(), &api.Req_FindCollateral{Id: id})
			if err != nil {
				return err
			}

			memo, err := actions.Build(cmd, core.ActionCatSupply, types.UUID(cat.Id))
			if err != nil {
				return err
			}

			return pay.Request(cmd.Context(), cat.Dai, amount, memo)
		},
	}

	return cmd
}
