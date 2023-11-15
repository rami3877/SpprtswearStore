package file

func UncodeFromDataBase(data []byte) map[string]string {
	types := make(map[string]string)

	key := ""
	Type := ""

	for i := 0; i < len(data); i++ {
		if data[i] == ':' {
			i++
			for ; i < len(data); i++ {
				if data[i] == '\r' {
					types[key] = Type
					key = ""
					Type = ""
					break
				}
				Type += string(data[i])
			}
		}
		if data[i] != '\r' {
			key += string(data[i])
		}

	}
	return types
}
