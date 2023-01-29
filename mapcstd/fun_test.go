package mapcstd

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func StringToInt(x string) int {
	return len(x)
}

func StringToTime(x string) (time.Time, error) {
	t, err := time.Parse(x, time.RFC3339)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

func StringToUUID(x string) (uuid.UUID, bool) {
	u, err := uuid.Parse(x)
	if err != nil {
		return [16]byte{}, false
	}
	return u, true
}

func TwoArguments(x, y string) int {
	return len(x) + len(y)
}

func ThreeReturns(x string) (int, bool, error) {
	_ = x
	return 0, true, nil
}

func TestNewFunFrom(t *testing.T) {
	tests := []struct {
		name    string
		arg     any
		want    *Fun
		wantErr bool
	}{
		{
			name: "normal",
			arg:  StringToInt,
			want: &Fun{
				srcTyp:  NewTypFrom("string"),
				destTyp: NewTypFrom(1),
				name:    "StringToInt",
				pkgPath: "github.com/abekoh/mapc/mapcstd",
				retType: OnlyValue,
			},
			wantErr: false,
		},
		{
			name: "with error",
			arg:  StringToTime,
			want: &Fun{
				srcTyp:  NewTypFrom("string"),
				destTyp: NewTypFrom(time.Time{}),
				name:    "StringToTime",
				pkgPath: "github.com/abekoh/mapc/mapcstd",
				retType: WithError,
			},
			wantErr: false,
		},
		{
			name: "with ok",
			arg:  StringToUUID,
			want: &Fun{
				srcTyp:  NewTypFrom("string"),
				destTyp: NewTypFrom(uuid.UUID{}),
				name:    "StringToUUID",
				pkgPath: "github.com/abekoh/mapc/mapcstd",
				retType: WithOk,
			},
			wantErr: false,
		},
		{
			name:    "failed: not a function",
			arg:     "hoge",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "failed: function has two arguments",
			arg:     TwoArguments,
			want:    nil,
			wantErr: true,
		},
		{
			name:    "failed: three returns",
			arg:     ThreeReturns,
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewFunFrom(tt.arg)
			if tt.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
			assert.Equalf(t, tt.want, got, "NewFunFrom(%v)", tt.arg)
		})
	}
}
