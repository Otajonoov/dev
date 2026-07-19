#!/usr/bin/env bash
# ci.sh — loyihaning to'liq tekshiruv zanjiri (15-bob).
# Lokal ham, CI'da ham AYNAN shu skript yuguradi — "mashinamda ishlagan edi"
# muammosiga qarshi yagona haqiqat manbai.
set -euo pipefail
cd "$(dirname "$0")"

echo "== gofmt tekshiruvi"
unformatted=$(gofmt -l .)
if [ -n "$unformatted" ]; then
    echo "gofmt talab qiladi:"
    echo "$unformatted"
    exit 1
fi

echo "== go vet"
go vet ./...

echo "== go build (examples bilan)"
go build ./...

echo "== go test -race"
go test ./... -race -count=1

echo "== qisqa fuzz (regression rejimi)"
# Har target bir necha soniya: yangi bug QIDIRMAYDI, corpus'dagi va seed'lardagi
# case'lar hali ham o'tishini tasdiqlaydi. Chuqur qidiruv qo'lda:
#   go test ./pdu -run='^$' -fuzz=FuzzDecode -fuzztime=60s
go test ./pdu -run='^$' -fuzz=FuzzDecode -fuzztime=3s
go test ./pdu -run='^$' -fuzz=FuzzReadFrame -fuzztime=3s
go test ./dlr -run='^$' -fuzz=FuzzParse -fuzztime=3s
go test ./coding -run='^$' -fuzz=FuzzGSM7RoundTrip -fuzztime=3s
go test ./coding -run='^$' -fuzz=FuzzSplit -fuzztime=3s

echo "== HAMMASI TOZA"
