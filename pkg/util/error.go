package util

func HandleError(err error) {
	if err != nil {
		PrintRed(err.Error() + "/n")
	}
}