package eventstore

type CurrentSequence func(current uint32) bool

func CheckSequence(current uint32, check CurrentSequence) bool {
	if check == nil {
		return true
	}
	return check(current)
}

func SequenceIgnore() CurrentSequence {
	return nil
}

func SequenceMatches(sequence uint32) CurrentSequence {
	return func(current uint32) bool {
		return current == sequence
	}
}

func SequenceAtLeast(sequence uint32) CurrentSequence {
	return func(current uint32) bool {
		return current >= sequence
	}
}
