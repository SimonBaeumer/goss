package goss_testing

//ConvertStringSliceToInterfaceSlice is a helper function to match
// system interfaces
func ConvertStringSliceToInterfaceSlice(strings []string) []interface{} {
    var iStrings = make([]interface{}, len(strings))
    for i, char := range strings {
        iStrings[i] = char
    }
    return iStrings
}