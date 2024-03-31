package utils

func AreSlicesEqual[T int | int64 | string](a, b []T) bool {
    if len(a) != len(b) {
        return false
    }

    for index, value := range a {
        if value != b[index] {
            return false
        }
    }

    return true
}
