package utils

import "sort"

func CheckIfExistIntArray[V int | float64 | string](s V, in []V) bool {
	for _, v := range in {
		if v == s {
			return true
		}
	}
	return false
}

func CheckIfHasCombination[V int | float64 | string](s []V, in []V) bool {
	for _, v := range s {
		if CheckIfExistIntArray[V](v, in) {
			return true
		}
	}
	return false
}

func CheckMissingValue[V int | float64 | string](new, existing []V) []V {
	missing := make([]V, 0)

	mapNew := make(map[V]struct{})
	for _, v := range new {
		mapNew[v] = struct{}{}
	}

	for _, v := range existing {
		if _, found := mapNew[v]; !found {
			missing = append(missing, v)
		}
	}

	return missing
}

func CombinationsValue[V int | float64 | string](new, existing []V) []V {
	combinations := make([]V, 0)

	mapNew := make(map[V]struct{})
	for _, v := range new {
		mapNew[v] = struct{}{}
	}

	for _, v := range existing {
		if _, found := mapNew[v]; found {
			combinations = append(combinations, v)
		}
	}

	return combinations
}

func PermutationValue[V int | float64 | string](new, existing []V) []V {
	permutations := make([]V, 0)

	newValue := append(new, existing...)

	mapNew := make(map[V]struct{})
	for _, v := range newValue {
		mapNew[v] = struct{}{}
	}

	for k, _ := range mapNew {
		permutations = append(permutations, k)
	}
	sort.SliceStable(permutations, func(i, j int) bool {
		return permutations[i] < permutations[j]
	})
	return permutations
}
