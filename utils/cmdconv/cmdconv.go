package cmdconv

func ToCmdLineStrings(cmd ...string) [][]byte {
	args := make([][]byte, len(cmd))
	for i, s := range cmd {
		args[i] = []byte(s)
	}
	return args
}

func ToCmdLineArgs(action string, args [][]byte) [][]byte {
	cmd := make([][]byte, 1)
	cmd[0] = []byte(action)
	return append(cmd, args...)
}
