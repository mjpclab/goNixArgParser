package goNixArgParser

func newArg(text string, argType argKind) *argToken {
	return &argToken{
		text: text,
		kind: argType,
	}
}
