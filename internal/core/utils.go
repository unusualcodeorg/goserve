package core

func PathBuilder(segments ...string) string {
	path := ""
	for _, segment := range segments {
		path += segment
	}
	return path
}
