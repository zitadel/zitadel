package setup

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/database"
)

func Test_eventstoreAutovacuum_Execute(t *testing.T) {
	type fields struct {
		Enabled               bool
		VacuumInsertThreshold uint32
		AnalyzeThreshold      uint32
	}
	tests := []struct {
		name    string
		fields  fields
		expects func(sqlmock.Sqlmock)
		wantErr bool
	}{
		{
			name: "enabled, sets thresholds",
			fields: fields{
				Enabled:               true,
				VacuumInsertThreshold: 50000,
				AnalyzeThreshold:      25000,
			},
			expects: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(`ALTER TABLE eventstore.events2 SET (
	autovacuum_vacuum_scale_factor = 0.0,
	autovacuum_analyze_scale_factor = 0.0,
	autovacuum_vacuum_insert_scale_factor = 0.0,
	autovacuum_vacuum_insert_threshold = 50000,
	autovacuum_analyze_threshold = 25000,
	autovacuum_vacuum_threshold = 50000
)`)).WillReturnResult(sqlmock.NewResult(0, 0))
			},
		},
		{
			name: "enabled, db error",
			fields: fields{
				Enabled:               true,
				VacuumInsertThreshold: 50000,
				AnalyzeThreshold:      50000,
			},
			expects: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(".*").WillReturnError(assert.AnError)
			},
			wantErr: true,
		},
		{
			name: "enabled, vacuum threshold at minimum boundary fails",
			fields: fields{
				Enabled:               true,
				VacuumInsertThreshold: minAutovacuumThreshold,
				AnalyzeThreshold:      50000,
			},
			expects: func(sqlmock.Sqlmock) {},
			wantErr: true,
		},
		{
			name: "enabled, analyze threshold at minimum boundary fails",
			fields: fields{
				Enabled:               true,
				VacuumInsertThreshold: 50000,
				AnalyzeThreshold:      minAutovacuumThreshold,
			},
			expects: func(sqlmock.Sqlmock) {},
			wantErr: true,
		},
		{
			name: "enabled, thresholds just above minimum succeed",
			fields: fields{
				Enabled:               true,
				VacuumInsertThreshold: minAutovacuumThreshold + 1,
				AnalyzeThreshold:      minAutovacuumThreshold + 1,
			},
			expects: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(`ALTER TABLE eventstore.events2 SET (
	autovacuum_vacuum_scale_factor = 0.0,
	autovacuum_analyze_scale_factor = 0.0,
	autovacuum_vacuum_insert_scale_factor = 0.0,
	autovacuum_vacuum_insert_threshold = 10001,
	autovacuum_analyze_threshold = 10001,
	autovacuum_vacuum_threshold = 10001
)`)).WillReturnResult(sqlmock.NewResult(0, 0))
			},
		},
		{
			name: "disabled, resets to defaults",
			fields: fields{
				Enabled: false,
			},
			expects: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(eventstoreAutovacuumResetStmt)).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
		},
		{
			name: "disabled, db error",
			fields: fields{
				Enabled: false,
			},
			expects: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(".*").WillReturnError(assert.AnError)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer func() {
				require.NoError(t, mock.ExpectationsWereMet())
			}()
			defer db.Close()
			tt.expects(mock)

			mig := &eventstoreAutovacuum{
				dbClient:              &database.DB{DB: db},
				Enabled:               tt.fields.Enabled,
				VacuumInsertThreshold: tt.fields.VacuumInsertThreshold,
				AnalyzeThreshold:      tt.fields.AnalyzeThreshold,
			}
			err = mig.Execute(context.Background(), nil)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func Test_eventstoreAutovacuum_Check(t *testing.T) {
	type fields struct {
		Enabled               bool
		VacuumInsertThreshold uint32
		AnalyzeThreshold      uint32
	}
	tests := []struct {
		name    string
		fields  fields
		lastRun map[string]any
		want    bool
	}{
		{
			name: "no previous run",
			fields: fields{
				Enabled:               true,
				VacuumInsertThreshold: 50000,
				AnalyzeThreshold:      50000,
			},
			lastRun: nil,
			want:    true,
		},
		{
			name: "unchanged",
			fields: fields{
				Enabled:               true,
				VacuumInsertThreshold: 50000,
				AnalyzeThreshold:      25000,
			},
			lastRun: map[string]any{
				"enabled":               true,
				"vacuumInsertThreshold": float64(50000),
				"analyzeThreshold":      float64(25000),
			},
			want: false,
		},
		{
			name: "enabled changed",
			fields: fields{
				Enabled:               false,
				VacuumInsertThreshold: 50000,
				AnalyzeThreshold:      50000,
			},
			lastRun: map[string]any{
				"enabled":               true,
				"vacuumInsertThreshold": float64(50000),
				"analyzeThreshold":      float64(50000),
			},
			want: true,
		},
		{
			name: "vacuum threshold changed",
			fields: fields{
				Enabled:               true,
				VacuumInsertThreshold: 60000,
				AnalyzeThreshold:      50000,
			},
			lastRun: map[string]any{
				"enabled":               true,
				"vacuumInsertThreshold": float64(50000),
				"analyzeThreshold":      float64(50000),
			},
			want: true,
		},
		{
			name: "analyze threshold changed",
			fields: fields{
				Enabled:               true,
				VacuumInsertThreshold: 50000,
				AnalyzeThreshold:      60000,
			},
			lastRun: map[string]any{
				"enabled":               true,
				"vacuumInsertThreshold": float64(50000),
				"analyzeThreshold":      float64(50000),
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mig := &eventstoreAutovacuum{
				Enabled:               tt.fields.Enabled,
				VacuumInsertThreshold: tt.fields.VacuumInsertThreshold,
				AnalyzeThreshold:      tt.fields.AnalyzeThreshold,
			}
			assert.Equal(t, tt.want, mig.Check(tt.lastRun))
		})
	}
}
