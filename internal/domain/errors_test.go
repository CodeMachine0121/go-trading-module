package domain

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDomainErrors(t *testing.T) {
	tests := []struct {
		name    string
		err     error
		wantErr bool
		wantMsg string
	}{
		{
			name:    "ErrInvalidStrategy should be defined",
			err:     ErrInvalidStrategy,
			wantErr: true,
			wantMsg: "invalid strategy",
		},
		{
			name:    "ErrInvalidPrice should be defined",
			err:     ErrInvalidPrice,
			wantErr: true,
			wantMsg: "invalid price",
		},
		{
			name:    "ErrStrategyNotFound should be defined",
			err:     ErrStrategyNotFound,
			wantErr: true,
			wantMsg: "strategy not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr {
				assert.Error(t, tt.err, "expected error to be defined")
				assert.True(t, errors.Is(tt.err, tt.err), "error should be the same instance")
			}
		})
	}
}

func TestErrorsAsType(t *testing.T) {
	t.Run("should check error type using errors.Is", func(t *testing.T) {
		assert.True(t, errors.Is(ErrInvalidStrategy, ErrInvalidStrategy))
		assert.False(t, errors.Is(ErrInvalidStrategy, ErrInvalidPrice))
	})
}
