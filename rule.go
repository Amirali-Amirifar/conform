package conform

import "github.com/amirali-amirifar/conform/predicate"

// Min requires input >= val.
func Min[T IntType](val T) predicate.Node[T] {
	return predicate.NewCmp(predicate.Ge, val)
}

// Max requires input <= val.
func Max[T IntType](val T) predicate.Node[T] {
	return predicate.NewCmp(predicate.Le, val)
}

// In restricts v to an allow-list of values.
func In[T IntType](allowed ...T) predicate.Node[T] {
	return predicate.NewIn(allowed...)
}
