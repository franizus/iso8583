package prefix

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/franizus/iso8583/encoding"
)

type ebcdicVarPrefixer struct {
	Digits int
}

var EBCDIC = Prefixers{
	Fixed: &ebcdicFixedPrefixer{},
	L:     &ebcdicVarPrefixer{1},
	LL:    &ebcdicVarPrefixer{2},
	LLL:   &ebcdicVarPrefixer{3},
	LLLL:  &ebcdicVarPrefixer{4},
}

func (p *ebcdicVarPrefixer) EncodeLength(maxLen, dataLen int) ([]byte, error) {
	if dataLen > maxLen {
		return nil, fmt.Errorf("field length: %d is larger than maximum: %d", dataLen, maxLen)
	}

	if len(strconv.Itoa(dataLen)) > p.Digits {
		return nil, fmt.Errorf("number of digits in length: %d exceeds: %d", dataLen, p.Digits)
	}

	strLen := fmt.Sprintf("%0*d", p.Digits, dataLen)
	res, err := encoding.EBCDIC.Encode([]byte(strLen))
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (p *ebcdicVarPrefixer) DecodeLength(maxLen int, data []byte) (int, error) {
	if len(data) < p.Length() {
		return 0, fmt.Errorf("length mismatch: want to read %d bytes, get only %d", p.Length(), len(data))
	}

	bDigits, _, err := encoding.EBCDIC.Decode(data[:p.Length()], p.Digits)
	if err != nil {
		return 0, err
	}

	digits := strings.TrimSpace(string(bDigits))

	dataLen, err := strconv.Atoi(digits)
	if err != nil {
		return 0, err
	}

	if dataLen > maxLen {
		return 0, fmt.Errorf("data length %d is larger than maximum %d", dataLen, maxLen)
	}

	return dataLen, nil
}

func (p *ebcdicVarPrefixer) Length() int {
	return p.Digits
}

func (p *ebcdicVarPrefixer) Inspect() string {
	return fmt.Sprintf("EBCDIC.%s", strings.Repeat("L", p.Digits))
}

type ebcdicFixedPrefixer struct {
}

func (p *ebcdicFixedPrefixer) EncodeLength(fixLen, dataLen int) ([]byte, error) {
	if dataLen > fixLen {
		return nil, fmt.Errorf("field length: %d should be fixed: %d", dataLen, fixLen)
	}

	return []byte{}, nil
}

// Returns number of characters that should be decoded
func (p *ebcdicFixedPrefixer) DecodeLength(fixLen int, data []byte) (int, error) {
	return fixLen, nil
}

func (p *ebcdicFixedPrefixer) Length() int {
	return 0
}

func (p *ebcdicFixedPrefixer) Inspect() string {
	return "EBCDIC.Fixed"
}
