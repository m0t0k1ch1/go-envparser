package envparser

import (
	"errors"
	"fmt"
	"math"
	"math/bits"
	"os"
	"strconv"
	"testing"

	"github.com/m0t0k1ch1/go-envparser/internal/testutils"
)

const (
	testEnvKey = "GO_ENVPARSER_TEST"
)

func TestParseFailedWithInvalidArgError(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		var iaerr *InvalidArgError
		err := Parse(testEnvKey, nil)
		testutils.Equal(t, true, errors.As(err, &iaerr))
		testutils.Contains(t, err.Error(), "v cannot be nil")
	})

	t.Run("non-pointer string", func(t *testing.T) {
		var iaerr *InvalidArgError
		err := Parse(testEnvKey, string(""))
		testutils.Equal(t, true, errors.As(err, &iaerr))
		testutils.Contains(t, err.Error(), "v cannot be non-pointer string")
	})

	t.Run("nil *string", func(t *testing.T) {
		var iaerr *InvalidArgError
		err := Parse(testEnvKey, (*string)(nil))
		testutils.Equal(t, true, errors.As(err, &iaerr))
		testutils.Contains(t, err.Error(), "v cannot be nil *string")
	})
}

func TestParseFailedWithUnsupportedTypeError(t *testing.T) {
	var b bool
	var uterr *UnsupportedTypeError
	err := Parse(testEnvKey, &b)
	testutils.Equal(t, true, errors.As(err, &uterr))
	testutils.Contains(t, err.Error(), "unsupported type: bool")
}

func TestParseFailedWithNotPresentError(t *testing.T) {
	var s string
	var nperr *NotPresentError
	err := Parse(testEnvKey, &s)
	testutils.Equal(t, true, errors.As(err, &nperr))
	testutils.Contains(t, err.Error(), fmt.Sprintf("%s is not present", testEnvKey))
}

func TestParseAsString(t *testing.T) {
	in := "string"
	out := "string"

	os.Setenv(testEnvKey, in)
	defer os.Clearenv()

	var s string
	if err := Parse(testEnvKey, &s); err != nil {
		t.Error(err)
	} else {
		testutils.Equal(t, out, s)
	}
}

func TestParseAsInt(t *testing.T) {
	cases := []struct {
		in  string
		out int
	}{{
		in:  minIntString(),
		out: minInt(),
	}, {
		in:  maxIntString(),
		out: maxInt(),
	}}

	for _, c := range cases {
		t.Run(c.in, func(t *testing.T) {
			os.Setenv(testEnvKey, c.in)
			defer os.Clearenv()

			var i int
			if err := Parse(testEnvKey, &i); err != nil {
				t.Error(err)
			} else {
				testutils.Equal(t, c.out, i)
			}
		})
	}
}

func TestParseAsIntFailedWithParseError(t *testing.T) {
	cases := []struct {
		in  string
		err string
	}{{
		in:  "zero",
		err: "invalid syntax",
	}, {
		in:  underIntString(),
		err: "value out of range",
	}, {
		in:  overIntString(),
		err: "value out of range",
	}}

	for _, c := range cases {
		t.Run(c.in, func(t *testing.T) {
			os.Setenv(testEnvKey, c.in)
			defer os.Clearenv()

			var i int
			var perr *ParseError
			var nerr *strconv.NumError
			err := Parse(testEnvKey, &i)
			testutils.Equal(t, true, errors.As(err, &perr))
			testutils.Equal(t, true, errors.As(err, &nerr))
			testutils.Contains(t, err.Error(), c.err)
		})
	}
}

func TestParseAsUint(t *testing.T) {
	cases := []struct {
		in  string
		out uint
	}{{
		in:  "0",
		out: 0,
	}, {
		in:  maxUintString(),
		out: maxUint(),
	}}

	for _, c := range cases {
		t.Run(c.in, func(t *testing.T) {
			os.Setenv(testEnvKey, c.in)
			defer os.Clearenv()

			var u uint
			if err := Parse(testEnvKey, &u); err != nil {
				t.Error(err)
			} else {
				testutils.Equal(t, c.out, u)
			}
		})
	}
}

func TestParseAsUintFailedWithParseError(t *testing.T) {
	cases := []struct {
		in  string
		err string
	}{{
		in:  "zero",
		err: "invalid syntax",
	}, {
		in:  overUintString(),
		err: "value out of range",
	}}

	for _, c := range cases {
		t.Run(c.in, func(t *testing.T) {
			os.Setenv(testEnvKey, c.in)
			defer os.Clearenv()

			var u uint
			var perr *ParseError
			var nerr *strconv.NumError
			err := Parse(testEnvKey, &u)
			testutils.Equal(t, true, errors.As(err, &perr))
			testutils.Equal(t, true, errors.As(err, &nerr))
			testutils.Contains(t, err.Error(), c.err)
		})
	}
}

func underIntString() string {
	if bits.UintSize == 32 {
		return "-2147483649"
	}

	return "-9223372036854775809"
}

func minInt() int {
	if bits.UintSize == 32 {
		return math.MinInt32
	}

	return math.MinInt64
}

func minIntString() string {
	return strconv.Itoa(minInt())
}

func maxInt() int {
	if bits.UintSize == 32 {
		return math.MaxInt32
	}

	return math.MaxInt64
}

func maxIntString() string {
	return strconv.Itoa(maxInt())
}

func overIntString() string {
	if bits.UintSize == 32 {
		return "2147483648"
	}

	return "9223372036854775808"
}

func maxUint() uint {
	if bits.UintSize == 32 {
		return math.MaxUint32
	}

	return math.MaxUint64
}

func maxUintString() string {
	return strconv.FormatUint(uint64(maxUint()), 10)
}

func overUintString() string {
	if bits.UintSize == 32 {
		return "4294967296"
	}

	return "18446744073709551616"
}
