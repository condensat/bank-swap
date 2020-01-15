// Copyright 2020 Condensat Tech. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package liquid

import (
	"context"

	"github.com/condensat/bank-core/cache"
	"github.com/condensat/bank-core/logger"
	"github.com/condensat/bank-core/utils"

	"github.com/sirupsen/logrus"
)

type Swap int

func (p *Swap) Run(ctx context.Context) {
	log := logger.Logger(ctx).WithField("Method", "Swap.Run")

	p.registerHandlers(cache.RedisMutexContext(ctx))

	log.WithFields(logrus.Fields{
		"Hostname": utils.Hostname(),
	}).Info("Liquid Swap Service started")

	<-ctx.Done()
}

func (p *Swap) registerHandlers(ctx context.Context) {
	log := logger.Logger(ctx).WithField("Method", "Liquid.RegisterHandlers")

	log.Debug("Liquid Swap registered")
}
