package kevlar_legacy

import (
	"github.com/boggydigital/testo"
	"strings"
	"testing"
)

func TestSha256(t *testing.T) {
	tests := []struct {
		in  string
		out string
	}{
		{"", "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"},
		{"1", "6b86b273ff34fce19d6b804eff5a3f5747ada4eaa22f1d49c01e52ddb7875b4b"},
		{"{key=value}", "8227a91a849cb343139ac0f4941cd77f284fc53f94acbfda6beb8679723ff59a"},
	}

	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			s, err := Sha256(strings.NewReader(tt.in))
			testo.EqualValues(t, s, tt.out)
			testo.Error(t, err, false)
		})
	}
}
