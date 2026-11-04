//                           _       _
// __      _____  __ ___   ___  __ _| |_ ___
// \ \ /\ / / _ \/ _` \ \ / / |/ _` | __/ _ \
//  \ V  V /  __/ (_| |\ V /| | (_| | ||  __/
//   \_/\_/ \___|\__,_| \_/ |_|\__,_|\__\___|
//
//  Copyright © 2016 - 2026 Weaviate B.V. All rights reserved.
//
//  CONTACT: hello@weaviate.io
//

package nearText

import (
	"context"

	"github.com/weaviate/weaviate/entities/moduletools"
)

type fakeTransformer struct{}

func (t *fakeTransformer) Transform(in []string) ([]string, error) {
	result := make([]string, len(in))
	for i, txt := range in {
		if txt == "transform this" {
			result[i] = "transformed text"
		} else {
			result[i] = txt
		}
	}
	return result, nil
}

// failingVectorizer succeeds on the first okCalls calls, then returns err.
type failingVectorizer struct {
	okCalls int
	err     error
}

func (v *failingVectorizer) Texts(context.Context, []string, moduletools.ClassConfig) ([]float32, error) {
	if v.okCalls > 0 {
		v.okCalls--
		return []float32{1, 0}, nil
	}
	return nil, v.err
}

type tenantlessClassConfig struct {
	moduletools.ClassConfig
}

func (tenantlessClassConfig) Tenant() string { return "" }
