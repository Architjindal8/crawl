package utils

func IsMaster() bool {
	return true // hardcoded for local dev
}

func GetNodeType() string {
	if IsMaster() {
		return "master"
	} else {
		return "worker"
	}
}
