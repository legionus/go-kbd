// Code generated by "stringer -type=Kind"; DO NOT EDIT

package lexer

import "fmt"

const _Kind_name = "UnknownSpaceNewLineCommentEqualsDashCommaPlusCharNumberUnicodeLiteralEscapedQuoteStrvalAltIsMetaStringsStringCharsetKeymapsKeycodePlainCapsShiftComposeControlCtrlLCtrlRAltGrAltShiftLShiftRShiftUsualForAsOnToInclude"

var _Kind_index = [...]uint8{0, 7, 12, 19, 26, 32, 36, 41, 45, 49, 55, 62, 69, 76, 81, 87, 96, 103, 109, 116, 123, 130, 135, 144, 151, 158, 163, 168, 173, 176, 182, 188, 193, 198, 201, 203, 205, 207, 214}

func (i Kind) String() string {
	if i < 0 || i >= Kind(len(_Kind_index)-1) {
		return fmt.Sprintf("Kind(%d)", i)
	}
	return _Kind_name[_Kind_index[i]:_Kind_index[i+1]]
}
