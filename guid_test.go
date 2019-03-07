package winguid

import (
	"strings"
	"testing"

	"golang.org/x/sys/windows"
)

var guidFixtures = []struct {
	Name      string
	Input     string
	Output    windows.GUID
	Roundtrip bool
}{
	{"NULL", "{00000000-0000-0000-0000-000000000000}", windows.GUID{0x00000000, 0x0000, 0x0000, [8]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}}, true},
	{"IUnknown", "{00000000-0000-0000-C000-000000000046}", windows.GUID{0x00000000, 0x0000, 0x0000, [8]byte{0xC0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x46}}, true},
	{"IDispatch", "{00020400-0000-0000-C000-000000000046}", windows.GUID{0x00020400, 0x0000, 0x0000, [8]byte{0xC0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x46}}, true},
	{"IEnumVariant", "{00020404-0000-0000-C000-000000000046}", windows.GUID{0x00020404, 0x0000, 0x0000, [8]byte{0xC0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x46}}, true},
	{"IConnectionPointContainer", "{B196B284-BAB4-101A-B69C-00AA00341D07}", windows.GUID{0xB196B284, 0xBAB4, 0x101A, [8]byte{0xB6, 0x9C, 0x00, 0xAA, 0x00, 0x34, 0x1D, 0x07}}, true},
	{"IConnectionPoint", "{B196B286-BAB4-101A-B69C-00AA00341D07}", windows.GUID{0xB196B286, 0xBAB4, 0x101A, [8]byte{0xB6, 0x9C, 0x00, 0xAA, 0x00, 0x34, 0x1D, 0x07}}, true},
	{"IInspectable", "{AF86E2E0-B12D-4C6A-9C5A-D7AA65101E90}", windows.GUID{0xaf86e2e0, 0xb12d, 0x4c6a, [8]byte{0x9c, 0x5a, 0xd7, 0xaa, 0x65, 0x10, 0x1e, 0x90}}, true},
	{"IProvideClassInfo", "{B196B283-BAB4-101A-B69C-00AA00341D07}", windows.GUID{0xb196b283, 0xbab4, 0x101a, [8]byte{0xB6, 0x9C, 0x00, 0xAA, 0x00, 0x34, 0x1D, 0x07}}, true},
	{"ICOMTestInt64", "{8D437CBC-B3ED-485C-BC32-C336432A1623}", windows.GUID{0x8d437cbc, 0xb3ed, 0x485c, [8]byte{0xbc, 0x32, 0xc3, 0x36, 0x43, 0x2a, 0x16, 0x23}}, true},
	{"Pattern1", "{10000000-1000-1000-1000-100000000000}", windows.GUID{0x10000000, 0x1000, 0x1000, [8]byte{0x10, 0x00, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00}}, true},
	{"Pattern2", "{01000000-0100-0100-0100-010000000000}", windows.GUID{0x01000000, 0x0100, 0x0100, [8]byte{0x01, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00}}, true},
	{"Pattern3", "{00100000-0010-0010-0010-001000000000}", windows.GUID{0x00100000, 0x0010, 0x0010, [8]byte{0x00, 0x10, 0x00, 0x10, 0x00, 0x00, 0x00, 0x00}}, true},
	{"Pattern4", "{00010000-0001-0001-0001-000100000000}", windows.GUID{0x00010000, 0x0001, 0x0001, [8]byte{0x00, 0x01, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00}}, true},
	{"Pattern5", "{a000a000-a000-a000-a000-a000a000a000}", windows.GUID{0xa000a000, 0xa000, 0xa000, [8]byte{0xa0, 0x00, 0xa0, 0x00, 0xa0, 0x00, 0xa0, 0x00}}, true},
	{"Pattern6", "{0aaa0aaa-0aaa-0aaa-0aaa-0aaa0aaa0aaa}", windows.GUID{0x0aaa0aaa, 0x0aaa, 0x0aaa, [8]byte{0x0a, 0xaa, 0x0a, 0xaa, 0x0a, 0xaa, 0x0a, 0xaa}}, true},
	{"Sequence1", "{12345678-1234-1234-1234-123456789abc}", windows.GUID{0x12345678, 0x1234, 0x1234, [8]byte{0x12, 0x34, 0x12, 0x34, 0x56, 0x78, 0x9a, 0xbc}}, true},
	{"Sequence2", "12345678-1234-1234-1234-123456789abc", windows.GUID{0x12345678, 0x1234, 0x1234, [8]byte{0x12, 0x34, 0x12, 0x34, 0x56, 0x78, 0x9a, 0xbc}}, false},
	{"Sequence3", "12345678123412341234123456789abc", windows.GUID{0x12345678, 0x1234, 0x1234, [8]byte{0x12, 0x34, 0x12, 0x34, 0x56, 0x78, 0x9a, 0xbc}}, false},
	{"CaseUpper1", "{ABCDEFAB-ABCD-ABCD-ABCD-ABCDEFABCDEF}", windows.GUID{0xabcdefab, 0xabcd, 0xabcd, [8]byte{0xab, 0xcd, 0xab, 0xcd, 0xef, 0xab, 0xcd, 0xef}}, true},
	{"CaseUpper2", "ABCDEFAB-ABCD-ABCD-ABCD-ABCDEFABCDEF", windows.GUID{0xabcdefab, 0xabcd, 0xabcd, [8]byte{0xab, 0xcd, 0xab, 0xcd, 0xef, 0xab, 0xcd, 0xef}}, false},
	{"CaseUpper3", "ABCDEFABABCDABCDABCDABCDEFABCDEF", windows.GUID{0xabcdefab, 0xabcd, 0xabcd, [8]byte{0xab, 0xcd, 0xab, 0xcd, 0xef, 0xab, 0xcd, 0xef}}, false},
	{"CaseLower1", "{abcdefab-abcd-abcd-abcd-abcdefabcdef}", windows.GUID{0xabcdefab, 0xabcd, 0xabcd, [8]byte{0xab, 0xcd, 0xab, 0xcd, 0xef, 0xab, 0xcd, 0xef}}, true},
	{"CaseLower2", "abcdefab-abcd-abcd-abcd-abcdefabcdef", windows.GUID{0xabcdefab, 0xabcd, 0xabcd, [8]byte{0xab, 0xcd, 0xab, 0xcd, 0xef, 0xab, 0xcd, 0xef}}, false},
	{"CaseLower3", "abcdefababcdabcdabcdabcdefabcdef", windows.GUID{0xabcdefab, 0xabcd, 0xabcd, [8]byte{0xab, 0xcd, 0xab, 0xcd, 0xef, 0xab, 0xcd, 0xef}}, false},
	{"CaseMixed1", "{AbCdEfAb-AbCd-AbCd-AbCd-AbCdEfAbCdEf}", windows.GUID{0xabcdefab, 0xabcd, 0xabcd, [8]byte{0xab, 0xcd, 0xab, 0xcd, 0xef, 0xab, 0xcd, 0xef}}, true},
	{"CaseMixed2", "{fEdCbAfE-fEdC-fEdC-fEdC-fEdCbAfEdCbA}", windows.GUID{0xfedcbafe, 0xfedc, 0xfedc, [8]byte{0xfe, 0xdc, 0xfe, 0xdc, 0xba, 0xfe, 0xdc, 0xba}}, true},
	{"Empty", "", windows.GUID{}, false},
	{"EmptyBrackets", "{}", windows.GUID{}, false},
	{"GarbageDash1", "----", windows.GUID{}, false},
	{"GarbageDash2", "------------------------------------", windows.GUID{}, false},
	{"GarbageDash3", "{------------------------------------}", windows.GUID{}, false},
	{"GarbagePadding1", " {abcdefab-abcd-abcd-abcd-abcdefabcdef}", windows.GUID{}, false},
	{"GarbagePadding2", "{abcdefab-abcd-abcd-abcd-abcdefabcdef} ", windows.GUID{}, false},
	{"GarbagePadding3", " abcdefab-abcd-abcd-abcd-abcdefabcdef", windows.GUID{}, false},
	{"GarbagePadding4", "abcdefab-abcd-abcd-abcd-abcdefabcdef ", windows.GUID{}, false},
	{"GarbagePadding5", " abcdefababcdabcdabcdabcdefabcdef", windows.GUID{}, false},
	{"GarbagePadding6", "abcdefababcdabcdabcdabcdefabcdef ", windows.GUID{}, false},
	{"Garbage1", "AFR*@)#$BNHRO*IABNFVaaa", windows.GUID{}, false},
	{"Garbage2", "#@*%@#&^%382765*@^#*&^%R*@&#%R7632", windows.GUID{}, false},
	{"Garbage3", "#@*%@#&^%382765*@^#*&^%R*@&#%R76377^2", windows.GUID{}, false},
	{"Garbage4", "{ABCDEFA*-ABCD-ABCD-ABCD-ABCDEFABCDEF}", windows.GUID{}, false},
	{"Garbage5", "{gggggggg-ABCD-ABCD-ABCD-ABCDEFABCDEF}", windows.GUID{}, false},
}

// TestGUID tests both NewGUID and GUID.String.
func TestGUID(t *testing.T) {
	for i := 0; i < len(guidFixtures); i++ {
		fixture := guidFixtures[i]
		result := New(fixture.Input)
		expected := fixture.Output
		if result != expected {
			t.Errorf("GUID test \"%v\" (%v of %v) failed. Expected %v from New. Received %v instead.", fixture.Name, i, len(guidFixtures), fixture.Output, result)
		}
		if fixture.Roundtrip && String(result) != strings.ToUpper(fixture.Input) {
			t.Errorf("GUID test \"%v\" (%v of %v) failed. Expected \"%v\" from GUID.String. Received \"%v\" instead.", fixture.Name, i, len(guidFixtures), strings.ToUpper(fixture.Input), result)
		}
	}
}
