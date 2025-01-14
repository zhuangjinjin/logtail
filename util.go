package logtail

func isNumberChar(b byte) bool {
	return b >= '0' && b <= '9'
}

func isAlphabetChar(b byte) bool {
	return (b >= 'A' && b <= 'Z') || (b >= 'a' && b <= 'z')
}

func indexToNextLineStart(format *Format, message []byte) []byte {
	l := len(message)
	i := 0

	for i < l {
		indexLineEnd(message, &l, &i)
		ignoreLineEnd(message, &l, &i)

		if format == nil || format.PrefixMatch(message[i:]) {
			return message[i:]
		}
	}

	return nil
}

func indexToLineStart(format *Format, data []byte) []byte {
	if format == nil || format.PrefixMatch(data) {
		return data
	}

	return indexToNextLineStart(format, data)
}

func isFollowingLine(format *Format, bytes []byte) bool {
	if format == nil {
		format = defaultFormat
	}

	if format != nil {
		return !format.PrefixMatch(bytes)
	}

	return bytes[0] == ' ' || bytes[0] == '\t'
}

func isLineEnd(b byte) bool {
	return b == '\n' || b == '\r'
}

func indexLineEnd(bytes []byte, length, index *int) {
	for ; *index < *length && !isLineEnd(bytes[*index]); *index++ {
	}
}

func ignoreLineEnd(bytes []byte, length, index *int) {
	for ; *index < *length && isLineEnd(bytes[*index]); *index++ {
	}
}

const DoubleSize = 2

func EscapeLimitJSONBytes(b []byte, capacity int) []byte {
	size := len(b)
	num := capacity

	if size < num {
		num = size
	}

	t := make([]byte, num*DoubleSize)

	index := 0
	from := 0

	for i := 0; i < num; i++ {
		for ; i < num && b[i] != '\n' && b[i] != '\t' && b[i] != '"'; i++ {
		}

		copy(t[index:], b[from:i])
		index += i - from
		from = i + 1

		if i < num {
			fillJSONEscape(t, &index, b[i])
		}
	}

	// from <= size means not reach the end of the bytes
	if from <= size {
		// remove uncompleted utf8 bytes
		for i := index - 1; i >= 0 && t[i]&0xC0 == 0x80; i-- {
			index = i - 1
		}
	}

	return t[:index]
}

func fillJSONEscape(t []byte, index *int, b byte) {
	t[*index] = '\\'
	*index++

	switch b {
	case '\n':
		t[*index] = 'n'
	case '\t':
		t[*index] = 't'
	case '"':
		t[*index] = '"'
	}

	*index++
}
