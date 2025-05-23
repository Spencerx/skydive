//go:build !linux

/*
 * Copyright (C) 2019 Red Hat, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
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

package targets

import (
	"github.com/google/gopacket"

	"github.com/skydive-project/skydive/api/types"
	"github.com/skydive-project/skydive/flow"
	"github.com/skydive-project/skydive/graffiti/graph"
	"github.com/skydive-project/skydive/probe"
)

// ERSpanTarget defines a ERSpanTarget target
type ERSpanTarget struct {
}

// SendPacket implements the Target interface
func (ers *ERSpanTarget) SendPacket(packet gopacket.Packet, bpf *flow.BPF) {
}

// Start start the target
func (ers *ERSpanTarget) Start() {
}

// Stop stops the target
func (ers *ERSpanTarget) Stop() {
}

// NewERSpanTarget returns a new ERSpan target
func NewERSpanTarget(g *graph.Graph, n *graph.Node, capture *types.Capture) (*ERSpanTarget, error) {
	return nil, probe.ErrNotImplemented
}
