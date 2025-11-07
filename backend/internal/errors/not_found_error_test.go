package errors

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNotFoundError(t *testing.T) {
	tests := []struct {
		name    string
		msg     string
		wantErr string
	}{
		{
			name:    "empty message",
			msg:     "",
			wantErr: "Data not found",
		},
		{
			name:    "with message",
			msg:     ": paymentId: 123",
			wantErr: "Data not found: paymentId: 123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewNotFoundError(tt.msg)
			require.Equal(t, tt.wantErr, err.Error())
		})
	}
}
