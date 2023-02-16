package modifiers

func SetStringEmpty() setStringDefaultModifier {
	var empty []string
	return setStringDefaultModifier{Default: &empty}
}
