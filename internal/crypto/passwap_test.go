package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zitadel/passwap/argon2"
	"github.com/zitadel/passwap/bcrypt"
	"github.com/zitadel/passwap/md5"
	"github.com/zitadel/passwap/scrypt"
)

func TestPasswordHasher_EncodingSupported(t *testing.T) {
	tests := []struct {
		name        string
		encodedHash string
		want        bool
	}{
		{
			name:        "empty string, false",
			encodedHash: "",
			want:        false,
		},
		{
			name:        "scrypt, false",
			encodedHash: "$scrypt$ln=16,r=8,p=1$cmFuZG9tc2FsdGlzaGFyZA$Rh+NnJNo1I6nRwaNqbDm6kmADswD1+7FTKZ7Ln9D8nQ",
			want:        false,
		},
		{
			name:        "bcrypt, true",
			encodedHash: "$2y$12$hXUrnqdq1RIIYZ2HPytIIe5lXdIvbhqrTvdPsSF7o.jFh817Z6lwm",
			want:        true,
		},
		{
			name:        "argo2i, true",
			encodedHash: "$argon2i$v=19$m=4096,t=3,p=1$cmFuZG9tc2FsdGlzaGFyZA$YMvo8AUoNtnKYGqeODruCjHdiEbl1pKL2MsYy9VgU/E",
			want:        true,
		},
		{
			name:        "argo2id, true",
			encodedHash: "$argon2d$v=19$m=4096,t=3,p=1$cmFuZG9tc2FsdGlzaGFyZA$CB0Du96aj3fQVcVSqb0LIA6Z6fpStjzjVkaC3RlpK9A",
			want:        true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &PasswordHasher{
				Prefixes: []string{bcrypt.Prefix, argon2.Prefix},
			}
			got := h.EncodingSupported(tt.encodedHash)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestPasswordHashConfig_PasswordHasher(t *testing.T) {
	type fields struct {
		Verifiers []HashName
		Hasher    HasherConfig
	}
	tests := []struct {
		name         string
		fields       fields
		wantPrefixes []string
		wantErr      bool
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
			wantErr: true,
		},
		{
			name: "missing algorithm",
			fields: fields{
				Hasher: HasherConfig{},
			},
			wantErr: true,
		},
		{
			name: "invalid md5",
			fields: fields{
				Hasher: HasherConfig{
					Algorithm: HashNameMd5,
				},
			},
			wantErr: true,
		},
		{
			name: "invalid argon2",
			fields: fields{
				Hasher: HasherConfig{
					Algorithm: HashNameArgon2,
					Params: map[string]any{
						"time":    3,
						"memory":  32768,
						"threads": 4,
					},
				},
			},
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
				Verifiers: []HashName{HashNameBcrypt, HashNameMd5, HashNameScrypt},
			},
			wantPrefixes: []string{argon2.Prefix, bcrypt.Prefix, md5.Prefix, scrypt.Prefix, scrypt.Prefix_Linux},
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
				Verifiers: []HashName{HashNameBcrypt, HashNameMd5, HashNameScrypt},
			},
			wantPrefixes: []string{argon2.Prefix, bcrypt.Prefix, md5.Prefix, scrypt.Prefix, scrypt.Prefix_Linux},
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
				Verifiers: []HashName{HashNameArgon2, HashNameMd5, HashNameScrypt},
			},
			wantPrefixes: []string{bcrypt.Prefix, argon2.Prefix, md5.Prefix, scrypt.Prefix, scrypt.Prefix_Linux},
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
				Verifiers: []HashName{HashNameArgon2, HashNameBcrypt, HashNameMd5},
			},
			wantPrefixes: []string{scrypt.Prefix, scrypt.Prefix_Linux, argon2.Prefix, bcrypt.Prefix, md5.Prefix},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &PasswordHashConfig{
				Verifiers: tt.fields.Verifiers,
				Hasher:    tt.fields.Hasher,
			}
			got, err := c.PasswordHasher()
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			if tt.wantPrefixes != nil {
				require.NotNil(t, got)
				assert.Equal(t, tt.wantPrefixes, got.Prefixes)
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
