package tile

import "beeb/carcassonne/util/directions"

type EdgeSignature EdgeArray[FeatureType]

func (sig *EdgeSignature) Compatible(otherSig *EdgeSignature) bool {
	//compare each element in the array
	for j := range otherSig {

		//we dont care about matching 0's
		if otherSig[j] == 0 {
			continue
		}

		//if a non-blank mismatch occurs then this can't be the right orientation
		if otherSig[j] != sig[j] {
			return false
		}
	}

	return true
}

// IsRiverCurving
// probably not a good impl.
func (sig *EdgeSignature) IsRiverCurving() bool {
	numRivers := 0
	for i := 0; i < 4; i++ {
		if sig[i] == River {
			numRivers++
		}
	}

	if numRivers < 2 {
		return false
	}

	for i := 0; i < 4; i++ {
		if sig[i] == River {
			if sig[i] != sig[directions.Compliment[directions.Direction(i)]] {
				return true
			}
		}
	}

	return false
}

func (sig *EdgeSignature) Contains(t FeatureType) bool {
	ea := EdgeArray[FeatureType](*sig)
	return ea.Contains(t)
}
