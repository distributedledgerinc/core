
#!/usr/bin/env python3

import operator
import argparse
import json
import sys
import dateutil.parser
import math

DECIMAL_UNIT = 3

def init_default_argument_parser(prog_desc, default_chain_id, default_start_time):
    parser = argparse.ArgumentParser(description=prog_desc)
    parser.add_argument(
        'exported_genesis',
        help='exported genesis.json file',
        type=argparse.FileType('r'), default=sys.stdin,
    )
    parser.add_argument('--chain-id', type=str, default=default_chain_id)
    parser.add_argument('--start-time', type=str, default=default_start_time)
    return parser


def main(argument_parser, process_genesis_func):
    args = argument_parser.parse_args()
    if args.chain_id.strip() == '':
        sys.exit('chain-id required')

    genesis = json.loads(args.exported_genesis.read())

    genesis = process_genesis_func(genesis=genesis, parsed_args=args,)
    print(json.dumps(genesis, indent=4, sort_keys=True))


def create_module_account(name, address, coins, permissions):
    return {
        'address': address,
        'coins': coins,
        'sequence_number': '0',
        'account_number': '0',
        'original_vesting': [],
        'delegated_free': [],
        'delegated_vesting': [],
        'start_time': '0',
        'end_time': '0',
        'module_name': name,
        'module_permissions': permissions,
        'vesting_schedules': []
    }

def process_raw_genesis(genesis, parsed_args):

    bondedAmt = 0
    notBondedAmt = 0
    for val in genesis['app_state']['staking']['validators']:
        if val['status'] == 2:
            bondedAmt += int(val['tokens'])
        elif val['status'] == 0 or val['status'] == 1:
            notBondedAmt += int(val['tokens'])
        else:
            sys.exit('Invalid validator status')
        
    for ubd in genesis['app_state']['staking']['unbonding_delegations']:
        for entry in ubd['entries']:
            notBondedAmt += int(entry['balance'])

    newAccounts = []
    # Change old genesis accounts to new format genesis account
    for acc in genesis['app_state']['accounts']:
        newAcc = {
            'address': acc['address'],
            'coins': acc['coins'],
            'sequence_number': acc['sequence_number'],
            'account_number': '0', # ignored set by the account keeper during InitGenesis
            'original_vesting': acc['original_vesting'],
            'delegated_free': acc['delegated_free'],
            'delegated_vesting': acc['delegated_vesting'],
            'start_time': '0',
            'end_time': '0',
            'module_name': '',
            'module_permissions': [],
            'vesting_schedules': acc['lazy_vesting_schedules']
        }

        if acc['address'] == 'terra1fs7mmpducjf25j70sk3sz6k5phz2fllmyr5gwz':
            update_vesting_schedule(newAcc)

        newAccounts.append(newAcc)

    # Create fee collector account
    newAccounts.append(
        create_module_account(
            'fee_collector', 
            'terra17xpfvakm2amg962yls6f84z3kell8c5lkaeqfa', 
            genesis['app_state']['auth']['collected_fees'], 
            ['basic']
        )
    )

    # Create empty gov account with burner permissions
    newAccounts.append(
        create_module_account(
            'gov', 
            'terra10d07y265gmmuvt4z0w9aw880jnsr700juxf95n', 
            [], 
            ['burner']
        )
    )

    # Create empty distribution account
    # We assume genesis is exported for zero-height without any remaining rewards
    # only cares for community-pool
    communityPoolCoins = []
    for coin in genesis['app_state']['distr']['fee_pool']['community_pool']:
        communityPoolCoins.append({
            'denom': coin['denom'],
            'amount': str(int(float(coin['amount'])))
        })

    newAccounts.append(
        create_module_account(
            'distribution', 
            'terra1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8pm7utl', 
            communityPoolCoins, 
            ['basic']
        )
    )

    # Create bonded_tokens_pool account with burner & staking permissions
    newAccounts.append(
        create_module_account(
            'bonded_tokens_pool', 
            'terra1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3nln0mh', 
            [{'amount': str(bondedAmt), 'denom': 'uluna'}], 
            ['burner', 'staking']
        )
    )

    # Create not_bonded_tokens_pool account with burner & staking permissions
    newAccounts.append(
        create_module_account(
            'not_bonded_tokens_pool', 
            'terra1tygms3xhhs3yv487phx3dw4a95jn7t7l8l07dr', 
            [{'amount': str(notBondedAmt), 'denom': 'uluna'}], 
            ['burner', 'staking']
        )
    )

    # Create empty oracle account
    newAccounts.append(
        create_module_account(
            'oracle', 
            'terra1jgp27m8fykex4e4jtt0l7ze8q528ux2lh4zh0f', 
            [], 
            ['basic']
        )
    )

    # Create empty market account
    newAccounts.append(
        create_module_account(
            'market', 
            'terra1untf85jwv3kt0puyyc39myxjvplagr3wstgs5s', 
            [], 
            ['minter', 'burner']
        )
    )

    # Create empty treasury account
    newAccounts.append(
        create_module_account(
            'treasury', 
            'terra1vmafl8f3s6uuzwnxkqz0eza47v6ecn0t0yeca7', 
            [], 
            ['minter']
        )
    )

    # Migrate accounts 
    genesis['app_state']['accounts'] = newAccounts

    # Migrate auth module 
    genesis['app_state']['auth'] = {
        'params': genesis['app_state']['auth']['params']
    }

    # Add gov module genesis state
    genesis['app_state']['gov'] = {
        'deposits': [],
        'proposals': [],
        'votes': [],
        'starting_proposal_id': '1',
        'deposit_params': {
            'max_deposit_period': '1209600000000000',
            'min_deposit': [{
                'amount': '512000000',
                'denom': 'uluna'
            }]
        },
        'voting_params': {
            'voting_period': '1209600000000000'
        },
        'tally_params': {
            'quorum': '0.400000000000000000',
            'threshold': '0.500000000000000000',
            'vete': '0.334000000000000000'
        }
    }

    # Move genesis state key from distr to distribution
    genesis['app_state']['distribution'] = genesis['app_state']['distr']
    del genesis['app_state']['distr']
    
    # Migrate distribution module
    for vse in genesis['app_state']['distribution']['validator_slash_events']:
        vse['period'] = vse['validator_slash_event']['validator_period']

    # Migrate staking module
    for val in genesis['app_state']['staking']['validators']:
        val['commission'] = {
            'commission_rates': {
                'rate': val['commission']['rate'],
                'max_rate': val['commission']['max_rate'],
                'max_change_rate': val['commission']['max_change_rate']
            },
            'update_time': val['commission']['update_time']
        }

    # Add supply module genesis state
    genesis['app_state']['supply'] = {
        'supply': []
    }

    # Migrate market moudle
    genesis['app_state']['market'] = {
        'base_pool': '0.000000000000000000',
        'terra_pool': '0.000000000000000000',
        'last_update_height': '0',
        'params': {
            'pool_update_interval': '100800',
            'daily_terra_liquidity_ratio': '0.010000000000000000',
            'min_spread': '0.020000000000000000',
            'tobin_tax': '0.003000000000000000'
        }
    }

    # Migrate oracle module
    genesis['app_state']['oracle'] = {
        'voting_infos': {},
        'missed_votes': {},
        'feeder_delegations': genesis['app_state']['oracle']['feeder_delegations'],
        'prices': {},
        'price_prevotes': [],
        'price_votes': [],
        'params': {
            'vote_period': genesis['app_state']['oracle']['params']['vote_period'],
            'vote_threshold': genesis['app_state']['oracle']['params']['vote_threshold'],
            'reward_band': genesis['app_state']['oracle']['params']['oracle_reward_band'],
            'votes_window': '1000',
            'min_valid_votes_per_window': '0.050000000000000000', # 5%
            'slash_fraction': '0.000100000000000000', # 0.01%
            'reward_fraction': '0.010000000000000000' # 1%
        }
    }

    # Migrate treasury module
    # TODO - should we need to add window_probation?
    # del genesis['app_state']['treasury']['params']['window_probation']
    genesis['app_state']['treasury'] = {
        'params': genesis['app_state']['treasury']['params'],
        'tax_rate': genesis['app_state']['treasury']['tax_rate'],
        'reward_weight': genesis['app_state']['treasury']['reward_weight'],
        'tax_cap': {},
        'tax_proceeds': [],
        'historical_issuance': []
    }
    
    # Set new chain ID and genesis start time
    genesis['chain_id'] = parsed_args.chain_id.strip()
    genesis['genesis_time'] = parsed_args.start_time

    return genesis

def update_vesting_schedule(account):
    initial_genesis_time = '2019-04-24T06:00:00.000000Z'
    
    # parse genesis date
    genesis_date = dateutil.parser.parse(initial_genesis_time)
    vesting_schedules = []

    # Luna Schedule Update
    luna_vesting_schedule = {
        'denom': 'uluna',
        'schedules': [
            {
                'start_time': str(get_time_after_n_month(genesis_date, 4)),
                'end_time': str(get_time_after_n_month(genesis_date, 5)),
                'ratio': '0.100000000000000000',
            },
            {
                'start_time': str(get_time_after_n_month(genesis_date, 5)),
                'end_time': str(get_time_after_n_month(genesis_date, 6)),
                'ratio': '0.100000000000000000',
            },
            {
                'start_time': str(get_time_after_n_month(genesis_date, 6)),
                'end_time': str(get_time_after_n_month(genesis_date, 7)),
                'ratio': '0.100000000000000000',
            },
            {
                'start_time': str(get_time_after_n_month(genesis_date, 14)),
                'end_time': str(get_time_after_n_month(genesis_date, 15)),
                'ratio': '0.700000000000000000',
            }
        ]
    }

    # Terra Schedule Update
    terra_schedules = []
    cumulated_ratio = 0
    for i in range(17):
        ratio = correct_decimal(1/18)
        cumulated_ratio += ratio
        terra_schedules.append({
            'start_time': str(get_time_after_n_month(genesis_date, 4 + i)),
            'end_time': str(get_time_after_n_month(genesis_date, 5 + i)),
            'ratio': str(ratio),
        })
    

    terra_schedules.append({
        'start_time': str(get_time_after_n_month(genesis_date, 4 + 17)),
        'end_time': str(get_time_after_n_month(genesis_date, 5 + 17)),
        'ratio': str(1 - cumulated_ratio),
    })

    terra_vesting_schedule = {
        'denom': 'usdr',
        'schedules': terra_schedules
    }

    vesting_schedules.append(luna_vesting_schedule)
    vesting_schedules.append(terra_vesting_schedule)
    account['vesting_schedules'] = vesting_schedules


def get_time_after_n_month(start_date, n):
    year = start_date.year
    month = start_date.month+n
    while month > 12:
        year += 1
        month -= 12

    return int(start_date.replace(year=year, month=month).timestamp())

def correct_decimal(float_num):
    return int(float_num * (10**DECIMAL_UNIT)) / (10**DECIMAL_UNIT)

if __name__ == '__main__':
    parser = init_default_argument_parser(
        prog_desc='Convert genesis.json for columbus-3-drill',
        default_chain_id='columbus-3-drill',
        default_start_time='2019-10-02T19:00:00Z',
    )
    main(parser, process_raw_genesis)


# TODO - terra1n2kzv00yjanjpjplqtwucug45lurr8tzgrvj2p -> terra1760w9ckqt5xh4urs235vm27xxkdwgeunay06v3