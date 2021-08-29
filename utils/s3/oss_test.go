package s3

import (
	"fmt"
	"github.com/sabhiram/go-gitignore"
	"testing"
)

func TestIgnore(t *testing.T) {
	i, err := ignore.CompileIgnoreFile("/Users/fredliang/GolandProjects/let.cli/.gitignore")
	if err != nil {
		t.Error(err)
	}
	path := "/Users/fredliang/GolandProjects/let.cli/"
	fmt.Println("========================================================")
	fmt.Println(".idea/12313/13123", i.MatchesPath(path+".idea/12313/13123"))
	fmt.Println(".idea", i.MatchesPath(path+".idea"))
	fmt.Println(".idea", i.MatchesPath(path+".idea"))

	fmt.Println(".iea/12313/13123", i.MatchesPath(path+".iea/12313/13123"))
	fmt.Println("dist/12313/13123", i.MatchesPath(path+"dist/12313/13123"))
	fmt.Println("13123", i.MatchesPath(path+"13123"))

}
