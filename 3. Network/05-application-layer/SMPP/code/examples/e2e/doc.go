// Package e2e — kitobning YAKUNIY demo'si (16-bob, loyiha DoD stsenariysi):
// mock SMSC (quirk'lar yoqiq) ustida bind_transceiver → o'zbek lotin
// (U+02BB'li) va kirill matnlar → concatenation (UDH) → DLR'lar →
// hex/decimal message_id quirk'i bilan TO'G'RI korrelyatsiya.
//
// Yugurish:
//
//	go test ./examples/e2e -v
package e2e
