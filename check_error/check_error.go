package check_error

func Check(e error) {
	if e != nil {
			panic(e)
	}
}
