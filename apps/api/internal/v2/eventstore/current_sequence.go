package eventstore

type CurrentSequence func(current uint32) bool

func CheckSequence(current uint32, check CurrentSequence) bool {
	if check == nil {
		return true
	}
	return check(current)
}

// SequenceIgnore doesn't check the current sequence
func SequenceIgnore() CurrentSequence {
	return nil
}

// SequenceMatches exactly the provided sequence
func SequenceMatches(sequence uint32) CurrentSequence {
	return func(current uint32) bool {
		return current == sequence
	}
}

// SequenceAtLeast matches the given sequence <= the current sequence
func SequenceAtLeast(sequence uint32) CurrentSequence {
	return func(current uint32) bool {
		return current >= sequence
	}
}
