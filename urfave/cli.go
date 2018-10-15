package urfave

import (
	"gopkg.in/urfave/cli.v2"
)

func ConcatFlags(flags ...[]cli.Flag) []cli.Flag {
	var result []cli.Flag

	for _, f := range flags {
		result = append(result, f...)
	}

	return result
}

func ChainBeforeFunc(before ...cli.BeforeFunc) cli.BeforeFunc {
	return func(ctx *cli.Context) error {
		for _, f := range before {
			if err := f(ctx); err != nil {
				return err
			}
		}
		return nil
	}
}
