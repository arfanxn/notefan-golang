package synch

// GetChanValAndKeep returns the given channel's value
// and keeps the value of the given channel in the given channel (persist value after retrived)
// By default if the value is retrieved from a channel
// the Channel value will be empty and it will block retriving value from a channel
// until the channel has a value
func GetChanValAndKeep[T any](ch chan T) T {
	val := <-ch // Get the value from the given channel
	ch <- val   // reassign value to keep the channel value
	return val  // return the channel value
}
