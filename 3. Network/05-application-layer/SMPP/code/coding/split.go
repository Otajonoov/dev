package coding

import (
	"fmt"
	"sync"

	"smpp/tlv"
)

// Method — uzun xabarni bo'lish usuli (8-bob taqqoslash jadvali).
type Method uint8

const (
	MethodUDH8    Method = iota // UDH, 8-bit reference — DEFAULT tanlov
	MethodUDH16                 // UDH, 16-bit reference (segment sig'imi 1 belgiga kam)
	MethodSarTLV                // sar_* TLV uchligi — PDU metadata darajasida
	MethodPayload               // message_payload TLV — SMSC o'zi bo'ladi
)

// Segment — yuborishga tayyor bitta qism.
type Segment struct {
	// Data — short_message'ga qo'yiladigan baytlar (UDH usullarida UDH bilan
	// boshlanadi). MethodPayload'da esa message_payload TLV value'si.
	Data []byte
	// UDH — UDH usullarida to'ldiriladi (ma'lumot uchun; baytlari Data boshida).
	// Bu segmentni yuborishda esm_class'ga UDHI (WithUDHI) qo'yish SHART.
	UDH *UDH
	// Sar — MethodSarTLV'da to'ldiriladi; PDU'ga SarTLVs() bilan qo'shiladi.
	Sar *SarInfo
}

// SarInfo — sar_* TLV uchligi qiymatlari (§5.3.2.22–24).
type SarInfo struct {
	RefNum uint16
	Total  uint8
	Seq    uint8
}

// TLVs sar uchligini PDU'ga qo'shishga tayyor ko'rinishda qaytaradi.
// UCHALASI birga — spec talabi (§5.3.2.22: qolgan ikkisisiz kelsa ignore).
func (s SarInfo) TLVs() []tlv.TLV {
	return []tlv.TLV{
		tlv.U16(tlv.SarMsgRefNum, s.RefNum),
		tlv.U8(tlv.SarTotalSegments, s.Total),
		tlv.U8(tlv.SarSegmentSeqnum, s.Seq),
	}
}

// SarFromTLVs kelgan PDU'ning TLV'laridan sar uchligini o'qiydi.
// Spec qoidasi (§5.3.2.22–24): uchchala TLV birga bo'lmasa — IGNORE
// (found=false, xato EMAS); qiymatlari buzuq bo'lsa — xato.
func SarFromTLVs(tlvs []tlv.TLV) (SarInfo, bool, error) {
	ref, refOK := tlv.Find(tlvs, tlv.SarMsgRefNum)
	tot, totOK := tlv.Find(tlvs, tlv.SarTotalSegments)
	seq, seqOK := tlv.Find(tlvs, tlv.SarSegmentSeqnum)
	if !refOK || !totOK || !seqOK {
		return SarInfo{}, false, nil // chala uchlik — spec bo'yicha e'tiborsiz
	}
	var s SarInfo
	var err error
	if s.RefNum, err = ref.Uint16Value(); err != nil {
		return SarInfo{}, false, err
	}
	if s.Total, err = tot.Uint8Value(); err != nil {
		return SarInfo{}, false, err
	}
	if s.Seq, err = seq.Uint8Value(); err != nil {
		return SarInfo{}, false, err
	}
	return s, true, nil
}

// Segment sig'imlari: hammasi havo interfeysining 140 oktetidan (7-bob).
const (
	singleGSM7  = 160 // septet
	singleUCS2  = 70  // 16-bit unit
	udh8GSM7    = 153 // (140-6)*8/7 = 153 (+1 fill bit)
	udh8UCS2    = 67  // (140-6)/2
	udh16GSM7   = 152 // (140-7)*8/7 = 152
	udh16UCS2   = 66  // (140-7)/2 = 66.5 → 66
	maxSegments = 255 // Total field'i 1 oktet
)

// Split matnni segmentlarga bo'ladi. dc — Choose'dan kelgan qiymat (DCDefault
// yoki DCUCS2); ref — xabar reference raqami (RefCounter.Next bilan oling;
// 8-bit usulda quyi bayti ishlatiladi). Matn bitta segmentga sig'sa UDH'siz,
// yagona segment qaytadi (method nima bo'lishidan qat'i nazar).
func Split(text string, dc DataCoding, method Method, ref uint16) ([]Segment, error) {
	switch dc {
	case DCDefault, DCUCS2:
	default:
		return nil, fmt.Errorf("coding: Split faqat GSM7/UCS2 uchun (dc=0x%02X emas)", uint8(dc))
	}

	// MethodPayload: bo'lish YO'Q — butun matn bitta value, SMSC o'zi segmentlaydi.
	if method == MethodPayload {
		data, err := encodeText(text, dc)
		if err != nil {
			return nil, err
		}
		return []Segment{{Data: data}}, nil
	}

	// Bitta segmentga sig'adimi?
	fits, err := fitsSingle(text, dc)
	if err != nil {
		return nil, err
	}
	if fits {
		data, err := encodeText(text, dc)
		if err != nil {
			return nil, err
		}
		return []Segment{{Data: data}}, nil
	}

	chunks, err := splitChunks(text, dc, method)
	if err != nil {
		return nil, err
	}
	if len(chunks) > maxSegments {
		return nil, fmt.Errorf("coding: %d segment — Total field'iga (max %d) sig'maydi", len(chunks), maxSegments)
	}

	total := uint8(len(chunks))
	segs := make([]Segment, 0, len(chunks))
	for i, chunk := range chunks {
		data, err := encodeText(chunk, dc)
		if err != nil {
			return nil, err
		}
		seq := uint8(i + 1)
		switch method {
		case MethodUDH8, MethodUDH16:
			u := &UDH{RefNum: ref, Total: total, Seq: seq, Is16bit: method == MethodUDH16}
			segs = append(segs, Segment{Data: append(u.Encode(), data...), UDH: u})
		case MethodSarTLV:
			segs = append(segs, Segment{Data: data, Sar: &SarInfo{RefNum: ref, Total: total, Seq: seq}})
		}
	}
	return segs, nil
}

// CountSegments matn uchun (Normalize + auto-detect'dan keyin) data_coding
// va default usulda (UDH 8-bit) ketadigan segment sonini aytadi — narx
// kalkulyatori. Natijadagi dc Choose bilan mos.
func CountSegments(text string) (DataCoding, int, error) {
	text = Normalize(text)
	dc := DCUCS2
	if _, err := EncodeGSM7(text); err == nil {
		dc = DCDefault
	}
	segs, err := Split(text, dc, MethodUDH8, 0)
	if err != nil {
		return dc, 0, err
	}
	return dc, len(segs), nil
}

func encodeText(text string, dc DataCoding) ([]byte, error) {
	if dc == DCDefault {
		return EncodeGSM7(text)
	}
	return EncodeUCS2(text), nil
}

func fitsSingle(text string, dc DataCoding) (bool, error) {
	if dc == DCDefault {
		n, err := SeptetLen(text)
		return n <= singleGSM7, err
	}
	return len(EncodeUCS2(text)) <= singleUCS2*2, nil
}

// splitChunks matnni belgi chegaralarida bo'ladi: GSM7'da septet byudjeti
// (extension belgi 2 septet — ESC juftligi HECH QACHON bo'linmaydi), UCS2'da
// unit byudjeti (surrogate pair 2 unit — u ham bo'linmaydi: ish rune
// darajasida ketadi).
func splitChunks(text string, dc DataCoding, method Method) ([]string, error) {
	budget := udh8GSM7
	if dc == DCUCS2 {
		budget = udh8UCS2
	}
	if method == MethodUDH16 {
		if dc == DCUCS2 {
			budget = udh16UCS2
		} else {
			budget = udh16GSM7
		}
	}

	var chunks []string
	var cur []rune
	used := 0
	for _, r := range text {
		cost, err := runeCost(r, dc)
		if err != nil {
			return nil, err
		}
		if used+cost > budget {
			chunks = append(chunks, string(cur))
			cur = cur[:0]
			used = 0
		}
		cur = append(cur, r)
		used += cost
	}
	if len(cur) > 0 {
		chunks = append(chunks, string(cur))
	}
	return chunks, nil
}

// runeCost — bitta belgining segment byudjetidagi narxi: GSM7'da septet
// (extension=2), UCS2'da 16-bit unit (surrogate pair=2).
func runeCost(r rune, dc DataCoding) (int, error) {
	if dc == DCDefault {
		return SeptetLen(string(r))
	}
	if r > 0xFFFF {
		return 2, nil // BMP tashqarisi — surrogate pair
	}
	return 1, nil
}

// RefCounter — har destination uchun alohida reference hisoblagich:
// bitta qabul qiluvchiga parallel ketgan uzun xabarlarning ref'lari
// to'qnashmasligi uchun (8-bit ref — atigi 256 qiymat; telefon segmentlarni
// ref + yuboruvchi manzil bo'yicha yig'adi).
type RefCounter struct {
	mu sync.Mutex
	m  map[string]uint16
}

// NewRefCounter bo'sh hisoblagich yaratadi.
func NewRefCounter() *RefCounter {
	return &RefCounter{m: make(map[string]uint16)}
}

// Next dest uchun navbatdagi reference'ni qaytaradi (uint16 aylanib ketishi
// normal — muhimi qo'shni xabarlar farqli bo'lsin).
func (rc *RefCounter) Next(dest string) uint16 {
	rc.mu.Lock()
	defer rc.mu.Unlock()
	rc.m[dest]++
	return rc.m[dest]
}
