package utils

func Contains(slice []uint, element uint) bool {
	for _, e := range slice {
		if e == element {
			return true
		}
	}
	return false
}

// remove an int from a list of ints
// change the order of the element for performance reasons
func PopIdFromList(ids *[]uint, id uint) {
	for i, el := range *ids {
		if el == id {
			// remove the id from the list
			// order does not matter
			(*ids)[i] = (*ids)[len(*ids)-1]
			*ids = (*ids)[:len(*ids)-1]
			break
		}
	}
}
