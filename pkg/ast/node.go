// Copyright 2026 dywoq
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ast

import "github.com/dywoq/miniasm/pkg/token"

type Node interface {
	Node()
}

type TopLevel struct {
	Identifier string `json:"identifier"`
	Expression Node   `json:"expression"`
}

type Value struct {
	Literal string     `json:"literal"`
	Kind    token.Kind `json:"kind"`
}

func (d TopLevel) Node() {}
func (v Value) Node()    {}
