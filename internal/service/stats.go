package service

import (
	"context"
	"github.com/lotostudio/financial-api/internal/domain"
	"github.com/lotostudio/financial-api/internal/repo"
	"golang.org/x/sync/errgroup"
)

type StatsService struct {
	accRepo   repo.Accounts
	balRepo   repo.Balances
	transRepo repo.Transactions
}

func newStatsService(accRepo repo.Accounts, balRepo repo.Balances, transRepo repo.Transactions) *StatsService {
	return &StatsService{
		accRepo:   accRepo,
		balRepo:   balRepo,
		transRepo: transRepo,
	}
}

func (s *StatsService) Statement(ctx context.Context, filter domain.TransactionsFilter) (domain.Statement, error) {
	var acc domain.Account
	var balIn, balOut domain.Balance
	var txs []domain.Transaction

	errs, ctx := errgroup.WithContext(ctx)

	errs.Go(func() error {
		var err error
		acc, err = s.accRepo.Get(ctx, *filter.AccountId)

		if err != nil {
			return err
		}

		return nil
	})

	errs.Go(func() error {
		var err error
		balIn, err = s.balRepo.Get(ctx, *filter.AccountId, *filter.CreatedFrom)

		if err != nil {
			return err
		}

		return nil
	})

	errs.Go(func() error {
		var err error
		balOut, err = s.balRepo.Get(ctx, *filter.AccountId, *filter.CreatedTo)

		if err != nil {
			return err
		}

		return nil
	})

	errs.Go(func() error {
		var err error
		txs, err = s.transRepo.List(ctx, filter)

		if err != nil {
			return err
		}

		return nil
	})

	if err := errs.Wait(); err != nil {
		return domain.Statement{}, err
	}

	balIn.Date = *filter.CreatedFrom
	balOut.Date = *filter.CreatedTo

	return domain.Statement{
		Account:      acc,
		BalanceIn:    balIn,
		BalanceOut:   balOut,
		Transactions: txs,
	}, nil
}
