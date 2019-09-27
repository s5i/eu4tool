package decode

import (
	"fmt"
	"reflect"
)

// DateFromBinaryMeta extracts the date from Clausewitz binary format "meta" file.
func DateFromBinaryMeta(b []byte) (d string, err error) {

	i := 0
	expectI := 0
	expect := []token{
		binCtl{ctlMagicNumber},
		binID{`date`},
		binCtl{ctlEquals},
	}
	for ; expectI < len(expect); expectI++ {
		t, n := getToken(b[i:])
		if !reflect.DeepEqual(t, expect[expectI]) {
			return "", fmt.Errorf("unexpected token (%T:%v) at i=%d", t, t, i)
		}
		i += n

		n, err := t.getN(b[i:])
		if err != nil {
			return "", fmt.Errorf("couldn't parse token %T at i=%d: %v", t, i, err)
		}
		i += n
	}

	t, n := getToken(b[i:])
	tInt, ok := t.(binInteger)
	if !ok {
		return "", fmt.Errorf("found token %T at i=%d, expected integer", t, i)
	}
	i += n

	date, n, err := tInt.getValue(b[i:])
	if err != nil {
		return "", fmt.Errorf("error parsing date integer at i=%d: %v", i, err)
	}
	i += n

	return formatDate(date), nil
}

func formatDate(x int) string {
	year := x/24/365 - 5000
	day := x/24%365 + 1
	month := 1

	for _, mLen := range []int{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31} {
		if day <= mLen {
			break

		}
		day -= mLen
		month++
	}
	return fmt.Sprintf("%04d_%02d_%02d", year, month, day)

}
