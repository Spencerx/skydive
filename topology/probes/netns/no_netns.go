//go:build !linux

/*
 * Copyright (C) 2018 Red Hat, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy ofthe License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specificlanguage governing permissions and
 * limitations under the License.
 *
 */

package netns

import (
	"github.com/skydive-project/skydive/probe"
	tp "github.com/skydive-project/skydive/topology/probes"
)

// NewProbe returns a new network namespace probe
func NewProbe(ctx tp.Context, bundle *probe.Bundle) (probe.Handler, error) {
	return nil, probe.ErrNotImplemented
}
