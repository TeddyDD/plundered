package cmdline

import (
	"io"
	"strings"
	"unicode"
)

// CmdLine lets people view the raw & parsed /proc/cmdline in one place
type CmdLine struct {
	Raw   string
	AsMap map[string]string
	Err   error
}

// NewCmdLine returns a populated CmdLine struct
func New(reader io.Reader) *CmdLine {
	return parse(reader)
}

func FromString(cmdline string) *CmdLine {
	return parse(strings.NewReader(cmdline))
}

// parse returns the current command line, trimmed
func parse(cmdlineReader io.Reader) *CmdLine {
	line := &CmdLine{}
	raw, err := io.ReadAll(cmdlineReader)
	line.Err = err
	// This works because string(nil) is ""
	line.Raw = strings.TrimRight(string(raw), "\n")
	line.AsMap = parseToMap(line.Raw)
	return line
}

func doParse(input string, handler func(flag, key, canonicalKey, value, trimmedValue string)) {
	lastQuote := rune(0)
	quotedFieldsCheck := func(c rune) bool {
		switch {
		case c == lastQuote:
			lastQuote = rune(0)
			return false
		case lastQuote != rune(0):
			return false
		case unicode.In(c, unicode.Quotation_Mark):
			lastQuote = c
			return false
		default:
			return unicode.IsSpace(c)
		}
	}

	for _, flag := range strings.FieldsFunc(string(input), quotedFieldsCheck) {
		// kernel variables must allow '-' and '_' to be equivalent in variable
		// names. We will replace dashes with underscores for processing.

		// Split the flag into a key and value, setting value="1" if none
		split := strings.Index(flag, "=")

		if len(flag) == 0 {
			continue
		}
		var key, value string
		if split == -1 {
			key = flag
			value = "1"
		} else {
			key = flag[:split]
			value = flag[split+1:]
		}
		canonicalKey := strings.Replace(key, "-", "_", -1)
		trimmedValue := strings.Trim(value, "\"'")

		// Call the user handler
		handler(flag, key, canonicalKey, value, trimmedValue)
	}
}

// parseToMap turns a space-separated kernel commandline into a map
func parseToMap(input string) map[string]string {
	flagMap := make(map[string]string)
	doParse(input, func(flag, key, canonicalKey, value, trimmedValue string) {
		// We store the value twice, once with dash, once with underscores
		// Just in case people check with the wrong method
		flagMap[canonicalKey] = trimmedValue
		flagMap[key] = trimmedValue
	})

	return flagMap
}

// ContainsFlag verifies that the kernel cmdline has a flag set
func (c *CmdLine) ContainsFlag(flag string) bool {
	_, present := c.Flag(flag)
	return present
}

// Flag returns the value of a flag, and whether it was set
func (c *CmdLine) Flag(flag string) (string, bool) {
	canonicalFlag := strings.Replace(flag, "-", "_", -1)
	value, present := c.AsMap[canonicalFlag]
	return value, present
}

func (c *CmdLine) AsBool(flag string) bool {
	canonicalFlag := strings.Replace(flag, "-", "_", -1)
	value, present := c.AsMap[canonicalFlag]
	if !present {
		return false
	}
	l := strings.ToLower(value)
	return l == "1" || l == "yes" || l == "true"
}
