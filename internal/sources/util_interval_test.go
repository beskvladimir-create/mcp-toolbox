// Copyright 2026 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sources

import (
	"encoding/json"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
)

func TestNormalizeValueIntervalAndTime(t *testing.T) {
	tcs := []struct {
		desc string
		val  any
		oid  uint32
		want any
	}{
		{"interval", pgtype.Interval{Microseconds: 7384500000, Days: 1, Valid: true}, 1186, "1 day 02:03:04.500000"},
		{"negative interval", pgtype.Interval{Microseconds: -1000000, Valid: true}, 1186, "-00:00:01"},
		{"interval with months", pgtype.Interval{Months: 14, Valid: true}, 1186, "14 mon 00:00:00"},
		{"time", pgtype.Time{Microseconds: 45296789000, Valid: true}, 1083, "12:34:56.789000"},
		{"NULL interval", nil, 1186, nil},
		{"interval value with other OID", pgtype.Interval{Microseconds: 1, Valid: true}, 1000, pgtype.Interval{Microseconds: 1, Valid: true}},
	}
	for _, tc := range tcs {
		t.Run(tc.desc, func(t *testing.T) {
			got := NormalizeValue(tc.val, tc.oid)
			if got != tc.want {
				t.Errorf("NormalizeValue(%#v, %d) = %#v, want %#v", tc.val, tc.oid, got, tc.want)
			}
		})
	}
}

func TestNormalizeValueIntervalJSON(t *testing.T) {
	// now() - xact_start in long_running_transactions, 6 minutes.
	got, err := json.Marshal(NormalizeValue(pgtype.Interval{Microseconds: 360000000, Valid: true}, 1186))
	if err != nil {
		t.Fatalf("json.Marshal returned error: %v", err)
	}
	if want := `"00:06:00"`; string(got) != want {
		t.Errorf("json.Marshal(NormalizeValue(interval)) = %s, want %s", got, want)
	}
}
