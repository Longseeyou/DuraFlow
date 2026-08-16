package mapper

func StringValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
