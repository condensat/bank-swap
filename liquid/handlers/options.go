// Copyright 2020 Condensat Tech. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package handlers

import (
	"fmt"
	"io"

	"github.com/condensat/bank-core/utils"
	"github.com/condensat/bank-swap/liquid/common"

	"github.com/condensat/bank-core/utils/shellexec"
)

type SwapCommand string

const (
	LiquidSwapCli = "liquidswap-cli"

	SwapCommandInfo     = SwapCommand("info")
	SwapCommandPropose  = SwapCommand("propose")
	SwapCommandFinalize = SwapCommand("finalize")
	SwapCommandAccept   = SwapCommand("accept")

	FeeRatePrecision       = 9 // BTC/Kb = 1000 / 100000000 sat/B
	FeeRatePrecisionFormat = "%.9f"
)

var (
	elementsConfFile = "/etc/liquidswap/elements.conf"
)

func SetElementsConf(elementsConf string) {
	elementsConfFile = elementsConf
}

func liquidSwapOptions(args ...interface{}) shellexec.Options {
	defaultEnv := []string{
		"LC_ALL=C.UTF-8",
		"LANG=C.UTF-8",
	}

	var payload io.Reader
	var finalArgs []string
	finalArgs = append(finalArgs, "--conf-file", elementsConfFile)

	for _, a := range args {
		switch arg := a.(type) {

		case string:
			finalArgs = append(finalArgs, arg)

		case common.ProposalInfo:
			finalArgs = append(finalArgs, arg.Args()...)

		case common.Payload:
			finalArgs = append(finalArgs, "-")
			payload = arg.Stdin()

		default:
			finalArgs = append(finalArgs, fmt.Sprintf("%v", arg))
		}
	}

	return shellexec.DefaultOptions().
		WithEnv(defaultEnv...).
		WithPath("/usr/local/bin").
		WithProgram(LiquidSwapCli).
		WithArgs(finalArgs...).
		WithStdin(payload)
}

func LiquidSwapPropose(address common.ConfidentialAddress, proposal common.ProposalInfo, feeRate float64) shellexec.Options {
	if feeRate < common.MinumumFeeRate {
		feeRate = common.MinumumFeeRate
	}
	feeRate = utils.ToFixed(feeRate, FeeRatePrecision)

	return liquidSwapOptions(
		"--with-address", address,
		SwapCommandPropose,
		"--fee-rate", fmt.Sprintf(FeeRatePrecisionFormat, feeRate),
		proposal)
}

func LiquidSwapInfo(payload common.Payload) shellexec.Options {
	return liquidSwapOptions(SwapCommandInfo, payload)
}

func LiquidSwapFinalize(payload common.Payload) shellexec.Options {
	return liquidSwapOptions(
		SwapCommandFinalize,
		"--send",
		payload,
	)
}

func LiquidSwapAccept(address common.ConfidentialAddress, payload common.Payload, feeRate float64) shellexec.Options {
	if feeRate < common.MinumumFeeRate {
		feeRate = common.MinumumFeeRate
	}
	feeRate = utils.ToFixed(feeRate, FeeRatePrecision)

	return liquidSwapOptions(
		"--with-address", address,
		SwapCommandAccept,
		"--fee-rate", fmt.Sprintf(FeeRatePrecisionFormat, feeRate),
		payload,
	)
}
