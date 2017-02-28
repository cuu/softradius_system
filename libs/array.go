package libs


func Slice_rm(sh string, arr []string) (ret []string){
	for i,v := range arr {
		if v == sh {
			ret = append( arr[:i],arr[i+1:]... )
		}
	}
	return
}
