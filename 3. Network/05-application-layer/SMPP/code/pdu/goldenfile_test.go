package pdu

import (
	"encoding/hex"
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// File-driven golden pattern (Eli Bendersky uslubi): kutilgan baytlar
// testdata/*.hex fayllarda, -update flag ular yangilaydi:
//
//	go test ./pdu -run TestGoldenFiles -update
//
// Hex-string (binary emas!) — diff o'qiladigan, git'da qulay. Kitobdagi
// asosiy golden'lar konstanta bo'lib qoldi (matn bilan yonma-yon turishi
// pedagogik qimmat); bu fayl PATTERN'ning o'zini namoyish qiladi — katta
// loyihada hammasi testdata/'ga ko'chadi.
var update = flag.Bool("update", false, "golden fayllarni qayta yozish")

func TestGoldenFiles(t *testing.T) {
	cases := []struct {
		name  string
		frame func() ([]byte, error)
	}{
		{"submit_sm", func() ([]byte, error) { return benchSubmit().Encode(42) }},
		{"enquire_link", func() ([]byte, error) { return EncodeEnquireLink(7), nil }},
		{"unbind_resp", func() ([]byte, error) { return EncodeUnbindResp(0, 9), nil }},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.frame()
			if err != nil {
				t.Fatal(err)
			}
			gotHex := strings.ToUpper(hex.EncodeToString(got))
			path := filepath.Join("testdata", tc.name+".hex")
			if *update {
				if err := os.WriteFile(path, []byte(gotHex+"\n"), 0o644); err != nil {
					t.Fatal(err)
				}
				return
			}
			want, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("%v (golden yo'qmi? -update bilan yarating)", err)
			}
			if gotHex != strings.TrimSpace(string(want)) {
				t.Errorf("%s golden mos emas:\n got  %s\n want %s", tc.name, gotHex, strings.TrimSpace(string(want)))
			}
		})
	}
}
