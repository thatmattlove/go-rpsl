package serialize_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mdl.wtf/rpsl"
	"go.mdl.wtf/rpsl/internal/serialize"
)

type WillUnmarshalBinaryErr string

func (w WillUnmarshalBinaryErr) UnmarshalBinary([]byte) (WillUnmarshalBinaryErr, error) {
	return "", errors.New("an error")
}

type OtherUnmarshalBinary string

func (o OtherUnmarshalBinary) UnmarshalBinary([]byte) error {
	return errors.New("an error")
}

type PtrUnmarshalBinary string

func (p PtrUnmarshalBinary) UnmarshalBinary([]byte) (*PtrUnmarshalBinary, error) {
	return &p, nil
}

func Test_Decode(t *testing.T) {
	t.Run("base", func(t *testing.T) {
		t.Parallel()
		b := []byte(`
aut-num: 65000
as-name: AS-ACME-1
`)
		var autNum rpsl.AutNum
		err := serialize.Decode(b, &autNum)
		require.NoError(t, err)
		assert.Equal(t, rpsl.ASN(65000), autNum.AutNum)
		assert.Equal(t, "AS-ACME-1", autNum.ASName)
	})
	t.Run("with append single line", func(t *testing.T) {
		t.Parallel()
		b := []byte(`
aut-num: 65000
as-name: AS-ACME-1
descr: A Description`)
		var autNum rpsl.AutNum
		err := serialize.Decode(b, &autNum)
		require.NoError(t, err)
		assert.Equal(t, "A Description", string(autNum.Description))
	})
	t.Run("with append multiline", func(t *testing.T) {
		t.Parallel()
		b := []byte(`
aut-num: 65000
as-name: AS-ACME-1
descr: Line 1
descr: Line 2`)
		var autNum rpsl.AutNum
		err := serialize.Decode(b, &autNum)
		require.NoError(t, err)
		assert.Equal(t, "Line 1\nLine 2", string(autNum.Description))
	})
	t.Run("with as comma", func(t *testing.T) {
		t.Parallel()
		b := []byte(`route-set: RS-ACME
members: 192.0.2.0/24,RS-CORP`)
		var rs rpsl.RouteSet
		err := serialize.Decode(b, &rs)
		require.NoError(t, err)
		exp := []string{"192.0.2.0/24", "RS-CORP"}
		assert.Equal(t, exp, rs.Members)
	})
	t.Run("with as comma space", func(t *testing.T) {
		t.Parallel()
		b := []byte(`aut-num: AS65000
as-name: AS-65000
member-of: AS65001, AS65002, AS-ACME`)
		var autNum rpsl.AutNum
		err := serialize.Decode(b, &autNum)
		require.NoError(t, err)
		exp := []string{"AS65001", "AS65002", "AS-ACME"}
		assert.Equal(t, exp, autNum.MemberOf)
	})
	t.Run("with as multiline", func(t *testing.T) {
		t.Parallel()
		b := []byte(`as-set: AS-ACME
members: AS65000
members: AS-65001`)
		var asSet rpsl.ASSet
		err := serialize.Decode(b, &asSet)
		require.NoError(t, err)
		exp := []string{"AS65000", "AS-65001"}
		assert.Equal(t, exp, asSet.Members)
	})
	t.Run("with as comma whitespace", func(t *testing.T) {
		t.Parallel()
		b := []byte(`route-set: RS-ACME
members: 192.0.2.0/24, RS-CORP
members: 198.51.100.0/24`)
		var rs rpsl.RouteSet
		err := serialize.Decode(b, &rs)
		require.NoError(t, err)
		exp := []string{"192.0.2.0/24", "RS-CORP", "198.51.100.0/24"}
		assert.Equal(t, exp, rs.Members)
	})
	t.Run("with as comma space no space", func(t *testing.T) {
		t.Parallel()
		b := []byte(`aut-num: AS65000
as-name: AS-65000
member-of: AS-A,AS-B`)
		var autNum rpsl.AutNum
		err := serialize.Decode(b, &autNum)
		require.NoError(t, err)
		exp := []string{"AS-A", "AS-B"}
		assert.Equal(t, exp, autNum.MemberOf)
	})
	t.Run("with as multiline comma separated", func(t *testing.T) {
		t.Parallel()
		b := []byte(`as-set: AS-ACME
members: AS65000,AS65001
members:
members: AS-A, AS-B,`)
		var asSet rpsl.ASSet
		err := serialize.Decode(b, &asSet)
		require.NoError(t, err)
		exp := []string{"AS65000", "AS65001", "AS-A", "AS-B"}
		assert.Equal(t, exp, asSet.Members)
	})
	t.Run("with comments", func(t *testing.T) {
		t.Parallel()
		b := []byte(`route-set: RS-ACME
descr: Line 1 # kept verbatim
members: 192.0.2.0/24, RS-CORP # comment`)
		var rs rpsl.RouteSet
		err := serialize.Decode(b, &rs)
		require.NoError(t, err)
		assert.Equal(t, "Line 1 # kept verbatim", rs.Description)
		assert.Equal(t, []string{"192.0.2.0/24", "RS-CORP"}, rs.Members)
	})
	t.Run("with comments parsed", func(t *testing.T) {
		t.Parallel()
		type Struct struct {
			Origin rpsl.ASN `rpsl:"origin"`
			Uint32 uint32   `rpsl:"uint32"`
		}
		b := []byte(`origin: AS65000 # comment
uint32: 12345 # comment`)
		var s Struct
		err := serialize.Decode(b, &s)
		require.NoError(t, err)
		assert.Equal(t, rpsl.ASN(65000), s.Origin)
		assert.Equal(t, uint32(12345), s.Uint32)
	})
	t.Run("with comments plain strings", func(t *testing.T) {
		t.Parallel()
		b := []byte(`route: 192.0.2.0/24 # comment
origin: AS65000
descr: Line 1 # kept verbatim
remarks: Line 2 # kept verbatim
source: RIPE # Filtered`)
		var route rpsl.Route
		err := serialize.Decode(b, &route)
		require.NoError(t, err)
		assert.Equal(t, "192.0.2.0/24", route.Route)
		assert.Equal(t, "RIPE", route.Source)
		assert.Equal(t, "Line 1 # kept verbatim", route.Description)
		assert.Equal(t, map[string]string{"remarks": "Line 2 # kept verbatim"}, route.Extra)
	})
	t.Run("with comments continuation lines", func(t *testing.T) {
		t.Parallel()
		b := []byte("route-set: RS-ACME\n" +
			"descr: Line 1 # kept verbatim\n" +
			"  continued\n" +
			"members: 192.0.2.0/24, # customer A\n" +
			"  198.51.100.0/24, # customer B\n" +
			"  RS-CORP\n" +
			"mp-members: 2001:db8::/32, # c\n" +
			"+ 2001:db8:1::/48")
		var rs rpsl.RouteSet
		err := serialize.Decode(b, &rs)
		require.NoError(t, err)
		assert.Equal(t, "Line 1 # kept verbatim continued", rs.Description)
		assert.Equal(t, []string{"192.0.2.0/24", "198.51.100.0/24", "RS-CORP"}, rs.Members)
		assert.Equal(t, []string{"2001:db8::/32", "2001:db8:1::/48"}, rs.MPMembers)
		var route rpsl.Route
		err = serialize.Decode([]byte("route: 192.0.2.0/24\norigin: # comment\n  AS65000"), &route)
		require.NoError(t, err)
		assert.Equal(t, rpsl.ASN(65000), route.Origin)
	})
	t.Run("with extra", func(t *testing.T) {
		t.Parallel()
		b := []byte(`as-set: AS-ACME
members: AS65000
members: AS-65001
extra1: value1`)
		var asSet rpsl.ASSet
		err := serialize.Decode(b, &asSet)
		require.NoError(t, err)
		assert.Equal(t, map[string]string{"extra1": "value1"}, asSet.Extra)
	})
	t.Run("with repeated extra", func(t *testing.T) {
		t.Parallel()
		b := []byte(`as-set: AS-ACME
extra1: value1
extra1: value2`)
		var asSet rpsl.ASSet
		err := serialize.Decode(b, &asSet)
		require.NoError(t, err)
		assert.Equal(t, map[string]string{"extra1": "value1\nvalue2"}, asSet.Extra)
	})
	t.Run("with repeated attributes", func(t *testing.T) {
		t.Parallel()
		b := []byte(`aut-num: AS65000
as-name: AS-65000
admin-c: ADMIN-1
admin-c: ADMIN-2
tech-c: TECH-1
tech-c: TECH-2
mnt-by: MNT-1
mnt-by: MNT-2
import: from AS65001 accept ANY
import: from AS65002 accept ANY
export: to AS65001 announce AS65000
export: to AS65002 announce AS65000
mp-import: afi ipv6.unicast from AS65001 accept ANY
mp-import: afi ipv6.unicast from AS65002 accept ANY
mp-export: afi ipv6.unicast to AS65001 announce AS65000
mp-export: afi ipv6.unicast to AS65002 announce AS65000
default: to AS65001
default: to AS65002
mp-default: afi ipv6.unicast to AS65001
mp-default: afi ipv6.unicast to AS65002
source: ARIN`)
		var autNum rpsl.AutNum
		err := serialize.Decode(b, &autNum)
		require.NoError(t, err)
		assert.Equal(t, "ADMIN-1\nADMIN-2", autNum.AdminPOC)
		assert.Equal(t, "TECH-1\nTECH-2", autNum.TechPOC)
		assert.Equal(t, "MNT-1\nMNT-2", autNum.MntBy)
		assert.Equal(t, "from AS65001 accept ANY\nfrom AS65002 accept ANY", autNum.Import)
		assert.Equal(t, "to AS65001 announce AS65000\nto AS65002 announce AS65000", autNum.Export)
		assert.Equal(t, "to AS65001\nto AS65002", autNum.Default)
		assert.Equal(t, map[string]string{}, autNum.Extra)
		result, err := serialize.Encode(&autNum)
		require.NoError(t, err)
		assert.Equal(t, b, result)
	})
	t.Run("with repeated remarks", func(t *testing.T) {
		t.Parallel()
		b := []byte(`route-set: RS-ACME
remarks: Line 1
remarks: Line 2`)
		var rs rpsl.RouteSet
		err := serialize.Decode(b, &rs)
		require.NoError(t, err)
		assert.Equal(t, "Line 1\nLine 2", rs.Remarks)
	})
	t.Run("without extra", func(t *testing.T) {
		t.Parallel()
		b := []byte(`route: 192.0.2.0/24
origin: AS65000
mnt-by: MNT-TEST
source: ARIN`)
		var route rpsl.Route
		err := serialize.Decode(b, &route)
		require.NoError(t, err)
		assert.Equal(t, map[string]string{}, route.Extra)
		result, err := serialize.Encode(&route)
		require.NoError(t, err)
		assert.Equal(t, b, result)
	})
	t.Run("case insensitive keys", func(t *testing.T) {
		t.Parallel()
		b := []byte(`Route: 192.0.2.0/24
ORIGIN: AS65000
Source: ARIN`)
		var route rpsl.Route
		err := serialize.Decode(b, &route)
		require.NoError(t, err)
		assert.Equal(t, "192.0.2.0/24", route.Route)
		assert.Equal(t, rpsl.ASN(65000), route.Origin)
		assert.Equal(t, "ARIN", route.Source)
		assert.Equal(t, map[string]string{}, route.Extra)
	})
	t.Run("err non ptr", func(t *testing.T) {
		t.Parallel()
		err := serialize.Decode([]byte(""), struct{}{})
		assert.ErrorContains(t, err, "non-pointer")
	})
	t.Run("err non struct", func(t *testing.T) {
		t.Parallel()
		o := []string{"one"}
		err := serialize.Decode([]byte{0x1, 0x2}, &o)
		assert.ErrorContains(t, err, "ptr *[]string")
	})
	t.Run("skip malformed pair", func(t *testing.T) {
		t.Parallel()
		b := []byte(`
aut-num: 65000
as-name: AS-ACME-1
extra content`)
		var autNum rpsl.AutNum
		err := serialize.Decode(b, &autNum)
		require.NoError(t, err)
		assert.Equal(t, rpsl.ASN(65000), autNum.AutNum)
		assert.Equal(t, "AS-ACME-1", autNum.ASName)
	})
	t.Run("colon in value", func(t *testing.T) {
		t.Parallel()
		b := []byte(`route6: 2001:db8::/32
origin: AS65000
remarks: see https://example.com:8080/path`)
		var route6 rpsl.Route6
		err := serialize.Decode(b, &route6)
		require.NoError(t, err)
		assert.Equal(t, "2001:db8::/32", route6.Route6)
		assert.Equal(t, rpsl.ASN(65000), route6.Origin)
		assert.Equal(t, "see https://example.com:8080/path", route6.Extra["remarks"])
	})
	t.Run("colon in multiline value", func(t *testing.T) {
		t.Parallel()
		b := []byte(`as-set: AS65000:AS-FOO
members: AS65000:AS-BAR`)
		var asSet rpsl.ASSet
		err := serialize.Decode(b, &asSet)
		require.NoError(t, err)
		assert.Equal(t, "AS65000:AS-FOO", asSet.ASSet)
		assert.Equal(t, []string{"AS65000:AS-BAR"}, asSet.Members)
	})
	t.Run("continuation lines", func(t *testing.T) {
		t.Parallel()
		b := []byte("aut-num: AS65000\n" +
			"as-name: AS-ACME-1\n" +
			"import: from AS65001\n" +
			"        accept ANY\n" +
			"mp-import: afi ipv6.unicast\n" +
			"\tfrom AS65002 2001:db8::1\n" +
			"+\n" +
			"+  accept ANY\n" +
			"export:\n" +
			"   \n" +
			"  to AS65001 announce AS65000")
		orig := bytes.Clone(b)
		var autNum rpsl.AutNum
		err := serialize.Decode(b, &autNum)
		require.NoError(t, err)
		assert.Equal(t, "AS-ACME-1", autNum.ASName)
		assert.Equal(t, "from AS65001 accept ANY", autNum.Import)
		assert.Equal(t, "afi ipv6.unicast from AS65002 2001:db8::1 accept ANY", autNum.MPImport)
		assert.Equal(t, "to AS65001 announce AS65000", autNum.Export)
		assert.Equal(t, orig, b)
	})
	t.Run("continuation multiline", func(t *testing.T) {
		t.Parallel()
		b := []byte(`aut-num: AS65000
descr: Line 1
  continued
descr: Line 2`)
		var autNum rpsl.AutNum
		err := serialize.Decode(b, &autNum)
		require.NoError(t, err)
		assert.Equal(t, "Line 1 continued\nLine 2", autNum.Description)
	})
	t.Run("indented object", func(t *testing.T) {
		t.Parallel()
		b := []byte(`
		% comment
		aut-num: AS65000
		as-name: AS-ACME-1
		import: from AS65001
		  accept ANY
		mp-import: afi ipv6.unicast
		+ from AS65002 accept ANY
		descr: Line 1
		`)
		exp := []byte(`aut-num: AS65000
as-name: AS-ACME-1
import: from AS65001
  accept ANY
mp-import: afi ipv6.unicast
+ from AS65002 accept ANY
descr: Line 1`)
		var autNum, expAutNum rpsl.AutNum
		err := serialize.Decode(b, &autNum)
		require.NoError(t, err)
		err = serialize.Decode(exp, &expAutNum)
		require.NoError(t, err)
		assert.Equal(t, expAutNum, autNum)
		assert.Equal(t, rpsl.ASN(65000), autNum.AutNum)
		assert.Equal(t, "from AS65001 accept ANY", autNum.Import)
		assert.Equal(t, "afi ipv6.unicast from AS65002 accept ANY", autNum.MPImport)
	})
	t.Run("skip continuation without attribute", func(t *testing.T) {
		t.Parallel()
		b := []byte(`  orphan: value
aut-num: AS65000`)
		var autNum rpsl.AutNum
		err := serialize.Decode(b, &autNum)
		require.NoError(t, err)
		assert.Equal(t, rpsl.ASN(65000), autNum.AutNum)
		assert.NotContains(t, autNum.Extra, "orphan")
	})
	t.Run("crlf", func(t *testing.T) {
		t.Parallel()
		b := []byte("aut-num: AS65000\r\n" +
			"as-name: AS-ACME-1\r\n" +
			"import: from AS65001\r\n" +
			"  accept ANY\r\n" +
			"+\r\n" +
			"descr: Line 1\r\n" +
			"descr: Line 2\r\n")
		var autNum rpsl.AutNum
		err := serialize.Decode(b, &autNum)
		require.NoError(t, err)
		assert.Equal(t, rpsl.ASN(65000), autNum.AutNum)
		assert.Equal(t, "AS-ACME-1", autNum.ASName)
		assert.Equal(t, "from AS65001 accept ANY", autNum.Import)
		assert.Equal(t, "Line 1\nLine 2", autNum.Description)
	})
	t.Run("skip comment lines", func(t *testing.T) {
		t.Parallel()
		b := []byte(`% Note: whois banner
%
aut-num: AS65000
# comment: line
as-name: AS-ACME-1
import: from AS65001
# comment: between continuation lines
  accept ANY`)
		var autNum rpsl.AutNum
		err := serialize.Decode(b, &autNum)
		require.NoError(t, err)
		assert.Equal(t, rpsl.ASN(65000), autNum.AutNum)
		assert.Equal(t, "AS-ACME-1", autNum.ASName)
		assert.Equal(t, "from AS65001 accept ANY", autNum.Import)
		assert.NotContains(t, autNum.Extra, "% Note")
		assert.NotContains(t, autNum.Extra, "# comment")
	})
	t.Run("skip empty tag", func(t *testing.T) {
		t.Parallel()
		type Struct struct {
			Key  string `rpsl:"key"`
			Skip string
		}
		b := []byte(`key: value
skip: skip`)
		var s Struct
		err := serialize.Decode(b, &s)
		require.NoError(t, err)
		assert.Equal(t, "value", s.Key)
		assert.Zero(t, s.Skip)
	})
	t.Run("err nested unmarshal", func(t *testing.T) {
		t.Parallel()
		type Struct struct {
			String string                 `rpsl:"string"`
			Error  WillUnmarshalBinaryErr `rpsl:"error"`
		}
		b := []byte(`string: value
error: irrelevant`)
		var s Struct
		err := serialize.Decode(b, &s)
		assert.ErrorContains(t, err, "an error")
	})
	t.Run("skip other unmarshal signature", func(t *testing.T) {
		t.Parallel()
		type Struct struct {
			Other OtherUnmarshalBinary `rpsl:"other"`
			Ptr   PtrUnmarshalBinary   `rpsl:"ptr"`
		}
		b := []byte(`other: value1
ptr: value2`)
		var s Struct
		err := serialize.Decode(b, &s)
		require.NoError(t, err)
		assert.Equal(t, OtherUnmarshalBinary("value1"), s.Other)
		assert.Equal(t, PtrUnmarshalBinary("value2"), s.Ptr)
	})
	t.Run("uint32 field", func(t *testing.T) {
		t.Parallel()
		type Struct struct {
			String string `rpsl:"string"`
			Uint32 uint32 `rpsl:"uint32"`
		}
		b := []byte(`string: value
uint32: 12345`)
		var s Struct
		err := serialize.Decode(b, &s)
		require.NoError(t, err)
	})
	t.Run("uint32 invalid", func(t *testing.T) {
		t.Parallel()
		type Struct struct {
			String string `rpsl:"string"`
			Uint32 uint32 `rpsl:"uint32"`
		}
		b := []byte(`string: value
uint32: wrong`)
		var s Struct
		err := serialize.Decode(b, &s)
		assert.ErrorContains(t, err, "value could not be parsed as uint32")
	})
	t.Run("uint32 overflow", func(t *testing.T) {
		t.Parallel()
		type Struct struct {
			Uint32 uint32 `rpsl:"uint32"`
		}
		var s Struct
		err := serialize.Decode([]byte(`uint32: 4294967295`), &s)
		require.NoError(t, err)
		assert.Equal(t, uint32(4294967295), s.Uint32)
		err = serialize.Decode([]byte(`uint32: 4294967296`), &s)
		assert.ErrorContains(t, err, "value could not be parsed as uint32")
	})
	t.Run("with as comma non-slice", func(t *testing.T) {
		t.Parallel()
		type Struct struct {
			String  string `rpsl:"string"`
			String2 string `rpsl:"string2" as:"comma"`
		}
		b := []byte(`string: value
string2: value2,value3,value4 `)
		var s Struct
		err := serialize.Decode(b, &s)
		require.NoError(t, err)
		assert.Equal(t, "value", s.String)
		assert.Equal(t, "value2,value3,value4", s.String2)
	})
	t.Run("with as comma space non-slice", func(t *testing.T) {
		t.Parallel()
		type Struct struct {
			String  string `rpsl:"string"`
			String2 string `rpsl:"string2" as:"comma-space"`
		}
		b := []byte(`string: value
string2:  value2, value3, value4`)
		var s Struct
		err := serialize.Decode(b, &s)
		require.NoError(t, err)
		assert.Equal(t, "value", s.String)
		assert.Equal(t, "value2, value3, value4", s.String2)
	})
}
