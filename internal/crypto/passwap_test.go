package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zitadel/passwap/argon2"
	"github.com/zitadel/passwap/scrypt"
)

func TestPasswordHashConfig_BuildSwapper(t *testing.T) {
	type fields struct {
		Verifiers []HashName
		Hasher    HasherConfig
	}
	tests := []struct {
		name    string
		fields  fields
		want    bool
		wantErr bool
	}{
		{
			name: "invalid verifier",
			fields: fields{
				Verifiers: []HashName{
					HashNameArgon2,
					HashNameBcrypt,
					HashNameMd5,
					HashNameScrypt,
					"foobar",
				},
				Hasher: HasherConfig{
					Algorithm: HashNameBcrypt,
					Params: map[string]any{
						"cost": 5,
					},
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "invalid hasher",
			fields: fields{
				Verifiers: []HashName{
					HashNameArgon2,
					HashNameBcrypt,
					HashNameMd5,
					HashNameScrypt,
				},
				Hasher: HasherConfig{
					Algorithm: "foobar",
					Params: map[string]any{
						"cost": 5,
					},
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "argon2i, error",
			fields: fields{
				Hasher: HasherConfig{
					Algorithm: HashNameArgon2i,
					Params: map[string]any{
						"time":    3,
						"threads": 4,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "argon2i, ok",
			fields: fields{
				Hasher: HasherConfig{
					Algorithm: HashNameArgon2i,
					Params: map[string]any{
						"time":    3,
						"memory":  32768,
						"threads": 4,
					},
				},
			},
			want: true,
		},
		{
			name: "argon2id, error",
			fields: fields{
				Hasher: HasherConfig{
					Algorithm: HashNameArgon2id,
					Params: map[string]any{
						"time":    3,
						"threads": 4,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "argon2id, ok",
			fields: fields{
				Hasher: HasherConfig{
					Algorithm: HashNameArgon2id,
					Params: map[string]any{
						"time":    3,
						"memory":  32768,
						"threads": 4,
					},
				},
			},
			want: true,
		},
		{
			name: "bcrypt, error",
			fields: fields{
				Hasher: HasherConfig{
					Algorithm: HashNameBcrypt,
					Params: map[string]any{
						"foo": 3,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "bcrypt, ok",
			fields: fields{
				Hasher: HasherConfig{
					Algorithm: HashNameBcrypt,
					Params: map[string]any{
						"cost": 3,
					},
				},
			},
			want: true,
		},
		{
			name: "scrypt, error",
			fields: fields{
				Hasher: HasherConfig{
					Algorithm: HashNameScrypt,
					Params: map[string]any{
						"cost": "bar",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "scrypt, ok",
			fields: fields{
				Hasher: HasherConfig{
					Algorithm: HashNameScrypt,
					Params: map[string]any{
						"cost": 3,
					},
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &PasswordHashConfig{
				Verifiers: tt.fields.Verifiers,
				Hasher:    tt.fields.Hasher,
			}
			got, err := c.BuildSwapper()
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			if tt.want {
				require.NotNil(t, got)
				encoded, err := got.Hash("password")
				require.NoError(t, err)
				assert.NotEmpty(t, encoded)
			}
		})
	}
}

func TestHasherConfig_decodeParams(t *testing.T) {
	type dst struct {
		A int
		B uint32
	}
	tests := []struct {
		name    string
		params  map[string]any
		want    dst
		wantErr bool
	}{
		{
			name: "unused",
			params: map[string]any{
				"a": 1,
				"b": 2,
				"c": 3,
			},
			wantErr: true,
		},
		{
			name: "unset",
			params: map[string]any{
				"a": 1,
			},
			wantErr: true,
		},
		{
			name: "wrong type",
			params: map[string]any{
				"a": 1,
				"b": "2",
			},
			wantErr: true,
		},
		{
			name: "ok",
			params: map[string]any{
				"a": 1,
				"b": 2,
			},
			want: dst{
				A: 1,
				B: 2,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &HasherConfig{
				Params: tt.params,
			}
			var got dst
			err := c.decodeParams(&got)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestHasherConfig_argon2Params(t *testing.T) {
	type fields struct {
		Params map[string]any
	}
	type args struct {
		p argon2.Params
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    argon2.Params
		wantErr bool
	}{
		{
			name: "decode error",
			fields: fields{
				Params: map[string]any{
					"foo": "bar",
				},
			},
			args: args{
				p: argon2.RecommendedIDParams,
			},
			wantErr: true,
		},
		{
			name: "ok",
			fields: fields{
				Params: map[string]any{
					"time":    2,
					"memory":  256,
					"threads": 8,
				},
			},
			args: args{
				p: argon2.RecommendedIDParams,
			},
			want: argon2.Params{
				Time:    2,
				Memory:  256,
				Threads: 8,
				KeyLen:  32,
				SaltLen: 16,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &HasherConfig{
				Params: tt.fields.Params,
			}
			got, err := c.argon2Params(tt.args.p)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestHasherConfig_bcryptCost(t *testing.T) {
	type fields struct {
		Params map[string]any
	}
	tests := []struct {
		name    string
		fields  fields
		want    int
		wantErr bool
	}{
		{
			name: "decode error",
			fields: fields{
				Params: map[string]any{
					"foo": "bar",
				},
			},
			wantErr: true,
		},
		{
			name: "ok",
			fields: fields{
				Params: map[string]any{
					"cost": 12,
				},
			},
			want: 12,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &HasherConfig{
				Params: tt.fields.Params,
			}
			got, err := c.bcryptCost()
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestHasherConfig_scryptParams(t *testing.T) {
	type fields struct {
		Params map[string]any
	}
	tests := []struct {
		name    string
		fields  fields
		want    scrypt.Params
		wantErr bool
	}{
		{
			name: "decode error",
			fields: fields{
				Params: map[string]any{
					"foo": "bar",
				},
			},
			wantErr: true,
		},
		{
			name: "ok",
			fields: fields{
				Params: map[string]any{
					"cost": 2,
				},
			},
			want: scrypt.Params{
				N:       4,
				R:       8,
				P:       1,
				KeyLen:  32,
				SaltLen: 16,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &HasherConfig{
				Params: tt.fields.Params,
			}
			got, err := c.scryptParams()
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
