// nolint
// autogenerated code using github.com/rigelrozanski/multitool
// aliases generated for the following subdirectories:
// ALIASGEN: github.com/terra-project/core/x/treasury/internal/types/
// ALIASGEN: github.com/terra-project/core/x/treasury/internal/keeper/
package treasury

import (
	"github.com/terra-project/core/x/treasury/internal/keeper"
	"github.com/terra-project/core/x/treasury/internal/types"
)

const (
	DefaultCodespace               = types.DefaultCodespace
	CodeInvalidEpoch               = types.CodeInvalidEpoch
	ModuleName                     = types.ModuleName
	StoreKey                       = types.StoreKey
	RouterKey                      = types.RouterKey
	QuerierRoute                   = types.QuerierRoute
	DefaultParamspace              = types.DefaultParamspace
	ProposalTypeTaxRateUpdate      = types.ProposalTypeTaxRateUpdate
	ProposalTypeRewardWeightUpdate = types.ProposalTypeRewardWeightUpdate
	QueryCurrentEpoch              = types.QueryCurrentEpoch
	QueryTaxRate                   = types.QueryTaxRate
	QueryTaxCap                    = types.QueryTaxCap
	QueryRewardWeight              = types.QueryRewardWeight
	QuerySeigniorageProceeds       = types.QuerySeigniorageProceeds
	QueryTaxProceeds               = types.QueryTaxProceeds
	QueryParameters                = types.QueryParameters
	QueryHistoricalIssuance        = types.QueryHistoricalIssuance
)

var (
	// functions aliases
	RegisterCodec                    = types.RegisterCodec
	ErrInvalidEpoch                  = types.ErrInvalidEpoch
	NewGenesisState                  = types.NewGenesisState
	DefaultGenesisState              = types.DefaultGenesisState
	ValidateGenesis                  = types.ValidateGenesis
	GetTaxRateKey                    = types.GetTaxRateKey
	GetRewardWeightKey               = types.GetRewardWeightKey
	GetTaxCapKey                     = types.GetTaxCapKey
	GetTaxProceedsKey                = types.GetTaxProceedsKey
	GetHistoricalIssuanceKey         = types.GetHistoricalIssuanceKey
	DefaultParams                    = types.DefaultParams
	NewTaxRateUpdateProposal         = types.NewTaxRateUpdateProposal
	NewRewardWeightUpdateProposal    = types.NewRewardWeightUpdateProposal
	NewQueryTaxCapParams             = types.NewQueryTaxCapParams
	NewQueryTaxRateParams            = types.NewQueryTaxRateParams
	NewQueryRewardWeightParams       = types.NewQueryRewardWeightParams
	NewQuerySeigniorageParams        = types.NewQuerySeigniorageParams
	NewQueryTaxProceedsParams        = types.NewQueryTaxProceedsParams
	NewQueryHistoricalIssuanceParams = types.NewQueryHistoricalIssuanceParams
	TaxRewardsForEpoch               = keeper.TaxRewardsForEpoch
	SeigniorageRewardsForEpoch       = keeper.SeigniorageRewardsForEpoch
	MiningRewardForEpoch             = keeper.MiningRewardForEpoch
	TRL                              = keeper.TRL
	SRL                              = keeper.SRL
	MRL                              = keeper.MRL
	UnitLunaIndicator                = keeper.UnitLunaIndicator
	SumIndicator                     = keeper.SumIndicator
	RollingAverageIndicator          = keeper.RollingAverageIndicator
	NewKeeper                        = keeper.NewKeeper
	ParamKeyTable                    = keeper.ParamKeyTable
	NewQuerier                       = keeper.NewQuerier

	// variable aliases
	ModuleCdc                            = types.ModuleCdc
	TaxRateKey                           = types.TaxRateKey
	RewardWeightKey                      = types.RewardWeightKey
	TaxCapKey                            = types.TaxCapKey
	TaxProceedsKey                       = types.TaxProceedsKey
	HistoricalIssuanceKey                = types.HistoricalIssuanceKey
	ParamStoreKeyTaxPolicy               = types.ParamStoreKeyTaxPolicy
	ParamStoreKeyRewardPolicy            = types.ParamStoreKeyRewardPolicy
	ParamStoreKeySeigniorageBurdenTarget = types.ParamStoreKeySeigniorageBurdenTarget
	ParamStoreKeyMiningIncrement         = types.ParamStoreKeyMiningIncrement
	ParamStoreKeyWindowShort             = types.ParamStoreKeyWindowShort
	ParamStoreKeyWindowLong              = types.ParamStoreKeyWindowLong
	ParamStoreKeyWindowProbation         = types.ParamStoreKeyWindowProbation
	DefaultTaxPolicy                     = types.DefaultTaxPolicy
	DefaultRewardPolicy                  = types.DefaultRewardPolicy
	DefaultSeigniorageBurdenTarget       = types.DefaultSeigniorageBurdenTarget
	DefaultMiningIncrement               = types.DefaultMiningIncrement
	DefaultWindowShort                   = types.DefaultWindowShort
	DefaultWindowLong                    = types.DefaultWindowLong
	DefaultWindowProbation               = types.DefaultWindowProbation
	DefaultTaxRate                       = types.DefaultTaxRate
	DefaultRewardWeight                  = types.DefaultRewardWeight
)

type (
	PolicyConstraints              = types.PolicyConstraints
	SupplyKeeper                   = types.SupplyKeeper
	MarketKeeper                   = types.MarketKeeper
	StakingKeeper                  = types.StakingKeeper
	DistributionKeeper             = types.DistributionKeeper
	GenesisState                   = types.GenesisState
	Params                         = types.Params
	TaxRateUpdateProposal          = types.TaxRateUpdateProposal
	RewardWeightUpdateProposal     = types.RewardWeightUpdateProposal
	QueryTaxCapParams              = types.QueryTaxCapParams
	QueryTaxRateParams             = types.QueryTaxRateParams
	QueryRewardWeightParams        = types.QueryRewardWeightParams
	QuerySeigniorageProceedsParams = types.QuerySeigniorageProceedsParams
	QueryTaxProceedsParams         = types.QueryTaxProceedsParams
	QueryHistoricalIssuanceParams  = types.QueryHistoricalIssuanceParams
	Keeper                         = keeper.Keeper
)
