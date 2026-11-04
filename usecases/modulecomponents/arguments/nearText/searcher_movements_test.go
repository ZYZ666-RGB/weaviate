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
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	enterrors "github.com/weaviate/weaviate/entities/errors"
)

func TestMoveVectorToAnother(t *testing.T) {
	t.Run("moving towards one another", func(t *testing.T) {
		type testCase struct {
			name           string
			source         []float32
			target         []float32
			weight         float32
			expectedResult []float32
			expectedError  error
		}

		tests := []testCase{
			{
				name:           "no force",
				source:         []float32{0, 1, 2, 3, 4},
				target:         []float32{0, 0, 0, 0, 0},
				weight:         0,
				expectedResult: []float32{0, 1, 2, 3, 4},
				expectedError:  nil,
			},
			{
				name:           "wrong vector sizes",
				source:         []float32{0, 1, 2, 3, 4},
				target:         []float32{0, 0, 0, 0},
				weight:         0,
				expectedResult: nil,
				expectedError:  fmt.Errorf("movement: vector lengths don't match: got 5 and 4"),
			},
			{
				name:           "force larger 1",
				source:         []float32{0, 1, 2, 3, 4},
				target:         []float32{0, 0, 0, 0, 0},
				weight:         1.5,
				expectedResult: nil,
				expectedError:  fmt.Errorf("movement: force must be between 0 and 1: got 1.500000"),
			},
			{
				name:           "force smaller 0",
				source:         []float32{0, 1, 2, 3, 4},
				target:         []float32{0, 0, 0, 0, 0},
				weight:         -0.2,
				expectedResult: nil,
				expectedError:  fmt.Errorf("movement: force must be between 0 and 1: got -0.200000"),
			},
			{
				name:           "force equals 1",
				source:         []float32{0, 1, 2, 3, 4},
				target:         []float32{1, 1, 1, 1, 1},
				weight:         1,
				expectedResult: []float32{0.5, 1, 1.5, 2, 2.5},
				expectedError:  nil,
			},
			{
				name:           "force equals 0.25",
				source:         []float32{0, 1, 2, 3, 4},
				target:         []float32{1, 1, 1, 1, 1},
				weight:         0.25,
				expectedResult: []float32{0.125, 1, 1.875, 2.75, 3.625},
				expectedError:  nil,
			},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				v := newMovements[[]float32]()
				res, err := v.MoveTo(test.source, test.target, test.weight)
				assert.Equal(t, test.expectedError, err)
				assert.Equal(t, test.expectedResult, res)
			})
		}
	})

	t.Run("moving away from one another", func(t *testing.T) {
		type testCase struct {
			name           string
			source         []float32
			target         []float32
			weight         float32
			expectedResult []float32
			expectedError  error
		}

		tests := []testCase{
			{
				name:           "no force",
				source:         []float32{1, 2, 3, 4},
				target:         []float32{0, 0, 0, 0},
				weight:         0,
				expectedResult: []float32{1, 2, 3, 4},
				expectedError:  nil,
			},
			{
				name:           "wrong vector sizes",
				source:         []float32{0, 1, 2, 3, 4},
				target:         []float32{0, 0, 0, 0},
				weight:         0,
				expectedResult: nil,
				expectedError:  fmt.Errorf("movement (moveAwayFrom): vector lengths don't match: got 5 and 4"),
			},
			{
				name:           "force smaller 0",
				source:         []float32{0, 1, 2, 3, 4},
				target:         []float32{0, 0, 0, 0, 0},
				weight:         -0.2,
				expectedResult: nil,
				expectedError:  fmt.Errorf("movement (moveAwayFrom): force must be 0 or positive: got -0.200000"),
			},
			{
				name:           "reproducing example from investigation period",
				source:         []float32{1.0, 1.0},
				target:         []float32{1.2, 0.8},
				weight:         1,
				expectedResult: []float32{0.9, 1.1},
				expectedError:  nil,
			},

			{
				name:           "force equals 0.25",
				source:         []float32{0, 1, 2, 3, 4},
				target:         []float32{1, 1, 1, 1, 1},
				weight:         0.25,
				expectedResult: []float32{-0.125, 1, 2.125, 3.25, 4.375},
				expectedError:  nil,
			},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				v := newMovements[[]float32]()
				res, err := v.MoveAwayFrom(test.source, test.target, test.weight)
				assert.Equal(t, test.expectedError, err)
				for i := range test.expectedResult {
					assert.InEpsilon(t, test.expectedResult[i], res[i], 0.01)
				}
			})
		}
	})

	t.Run("should error for multivectors", func(t *testing.T) {
		source := [][]float32{{0, 1}, {2, 3}}
		target := [][]float32{{0, 1}, {2, 3}}

		t.Run("multivectors moveTo", func(t *testing.T) {
			v := newMovements[[][]float32]()
			_, err := v.MoveTo(source, target, 0)
			assert.Equal(t, err, fmt.Errorf("move to operations are not applicable for multivector embeddings"))
		})

		t.Run("multivectors moveAwayFrom", func(t *testing.T) {
			v := newMovements[[][]float32]()
			_, err := v.MoveAwayFrom(source, target, 0)
			assert.Equal(t, err, fmt.Errorf("move away from operations are not applicable for multivector embeddings"))
		})
	})
}

// TestSearcherKeepsModulePrecondition: the REST handler maps the module's
// typed error to 422, so the searcher's wraps must keep it reachable.
func TestSearcherKeepsModulePrecondition(t *testing.T) {
	const msg = "no api key found neither in request header: X-Openai-Api-Key nor in environment variable under OPENAI_APIKEY"

	tests := []struct {
		name    string
		params  *NearTextParams
		okCalls int
		wantMsg string
	}{
		{
			name:    "keywords",
			params:  &NearTextParams{Values: []string{"space"}},
			wantMsg: "vectorize keywords: " + msg,
		},
		{
			name: "move to",
			params: &NearTextParams{
				Values: []string{"space"},
				MoveTo: ExploreMove{Values: []string{"stars"}, Force: 0.5},
			},
			okCalls: 1,
			wantMsg: "vectorize move to: vectorize move to: " + msg,
		},
		{
			name: "move away from",
			params: &NearTextParams{
				Values:       []string{"space"},
				MoveAwayFrom: ExploreMove{Values: []string{"stars"}, Force: 0.5},
			},
			okCalls: 1,
			wantMsg: "vectorize move away from: vectorize move to: " + msg,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vectorizer := &failingVectorizer{okCalls: tt.okCalls, err: enterrors.NewErrModulePrecondition(errors.New(msg))}
			searcher := NewSearcher[[]float32](vectorizer)

			_, err := searcher.VectorSearches()["nearText"](context.Background(), tt.params, "Movie", nil, tenantlessClassConfig{})
			require.Error(t, err)
			require.EqualError(t, err, tt.wantMsg)
			var precondition enterrors.ErrModulePrecondition
			require.ErrorAs(t, err, &precondition)
			require.EqualError(t, precondition, msg)
		})
	}
}
