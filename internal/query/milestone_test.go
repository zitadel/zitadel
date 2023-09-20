package query

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"regexp"
	"testing"
)

var (
	expectedMilestoneQuery = regexp.QuoteMeta(`
		SELECT projections.milestones.instance_id,
		   projections.milestones.primary_domain,
		   projections.milestones.reached_date,
		   projections.milestones.last_pushed_date,
		   projections.milestones.type,
		   COUNT(*) OVER ()
		FROM projections.milestones AS OF SYSTEM TIME '-1 ms'
		`)

	milestoneCols = []string{
		"instance_id",
		"primary_domain",
		"reached_date",
		"last_pushed_date",
		"type",
		"ignore_client_ids",
	}
)

func Test_MilestonesPrepare(t *testing.T) {
	type want struct {
		sqlExpectations sqlExpectation
		err             checkErr
	}
	tests := []struct {
		name    string
		prepare interface{}
		want    want
		object  interface{}
	}{
		{
			name:    "prepareMilestonesQuery no result",
			prepare: prepareMilestonesQuery,
			want: want{
				sqlExpectations: mockQueries(
					expectedMilestoneQuery,
					nil,
					nil,
				),
			},
			object: &Milestones{Milestones: []*Milestone{}},
		},
		{
			name:    "prepareMilestonesQuery",
			prepare: prepareMilestonesQuery,
			want: want{
				sqlExpectations: mockQueries(
					expectedMilestoneQuery,
					milestoneCols,
					[][]driver.Value{
						{
							"instance-id",
							"primary.domain",
							testNow,
							testNow,
							1,
							1,
						},
					},
				),
			},
			object: &Milestones{
				SearchResponse: SearchResponse{
					Count: 1,
				},
				Milestones: []*Milestone{
					{
						InstanceID:    "instance-id",
						Type:          1,
						ReachedDate:   testNow,
						PushedDate:    testNow,
						PrimaryDomain: "primary.domain",
					},
				},
			},
		},
		{
			name:    "prepareMilestonesQuery multiple result",
			prepare: prepareMilestonesQuery,
			want: want{
				sqlExpectations: mockQueries(
					expectedMilestoneQuery,
					milestoneCols,
					[][]driver.Value{
						{
							"instance-id",
							"primary.domain",
							testNow,
							testNow,
							1,
							1,
						},
						{
							"instance-id",
							"primary.domain",
							testNow,
							testNow,
							2,
							2,
						},
						{
							"instance-id",
							"primary.domain",
							testNow,
							nil,
							3,
							3,
						},
						{
							"instance-id",
							"primary.domain",
							nil,
							nil,
							4,
							4,
						},
					},
				),
			},
			object: &Milestones{
				SearchResponse: SearchResponse{
					Count: 4,
				},
				Milestones: []*Milestone{
					{
						InstanceID:    "instance-id",
						Type:          1,
						ReachedDate:   testNow,
						PushedDate:    testNow,
						PrimaryDomain: "primary.domain",
					},
					{
						InstanceID:    "instance-id",
						Type:          2,
						ReachedDate:   testNow,
						PushedDate:    testNow,
						PrimaryDomain: "primary.domain",
					},
					{
						InstanceID:    "instance-id",
						Type:          3,
						ReachedDate:   testNow,
						PrimaryDomain: "primary.domain",
					},
					{
						InstanceID:    "instance-id",
						Type:          4,
						PrimaryDomain: "primary.domain",
					},
				},
			},
		},
		{
			name:    "prepareMilestonesQuery sql err",
			prepare: prepareMilestonesQuery,
			want: want{
				sqlExpectations: mockQueryErr(
					expectedMilestoneQuery,
					sql.ErrConnDone,
				),
				err: func(err error) (error, bool) {
					if !errors.Is(err, sql.ErrConnDone) {
						return fmt.Errorf("err should be sql.ErrConnDone got: %w", err), false
					}
					return nil, true
				},
			},
			object: (*Milestones)(nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertPrepare(t, tt.prepare, tt.object, tt.want.sqlExpectations, tt.want.err, defaultPrepareArgs...)
		})
	}
}
