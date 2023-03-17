package chanh

// GetValAndKeep returns the given channel's value
// and keeps the value of the given channel in the given channel (persist value after retrived)
// By default if the value is retrieved from a channel
// the Channel value will be empty and it will block retriving value from a channel
// until the channel has a value
func GetValAndKeep[T any](ch chan T) T {
	val := <-ch // Get the value from the given channel
	ch <- val   // reassign value to keep the channel value
	return val  // return the channel value
}

// Make makes a new channel with value
func Make[T any](val T, size ...int) chan T {
	ch := make(chan T)
	if len(size) != 0 {
		ch = make(chan T, size[0])
	}
	ch <- val
	return ch
}

// ReplaceVal replaces the value at the given channel with the given value
func ReplaceVal[T any](ch chan T, vals ...T) {
	if cap(ch) == 1 { // if channel is only have one capacity
		switch len(ch) {
		case 1:
			<-ch
			ch <- vals[0]
		case 0:
			ch <- vals[0]
		}
		return
	} // otherwise
	for _, val := range vals {
		if len(ch) >= 1 {
			<-ch
			ch <- val
		} else {
			ch <- val
		}
	}
}
