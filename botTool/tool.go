package botTool

func Contains(slice map[string]struct{}, item string) bool {
    _, ok := slice[item] 
    return ok
}