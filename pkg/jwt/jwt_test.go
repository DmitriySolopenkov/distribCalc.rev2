package jwt

import (
	"testing"
)

func TestCreateUserToken(t *testing.T) {
	var userId = 1
	tokenJwt := New()
	_, err := tokenJwt.CreateUserToken(userId)

	if err != nil {
		t.Fatalf("CreateUsetToken(%d) return error %s", userId, err.Error())
	}
}

func TestTokenJwt_ParseToken(t *testing.T) {
	tokenJwt := New()
	token, _ := tokenJwt.CreateUserToken(1)

	tests := []struct {
		name    string
		token   string
		want    int
		wantErr bool
	}{
		{
			name:    "valid",
			token:   token.Token,
			want:    1,
			wantErr: false,
		},
		{
			name:    "invalid",
			token:   "",
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tokenJwt.ParseToken(tt.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got.ID != tt.want {
				t.Errorf("ParseToken() got = %v, want %v", got.ID, tt.want)
			}
		})
	}
}
