package utilities

// CheckErr checks if an error occurred
// If an error was found, the method panics
func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}
