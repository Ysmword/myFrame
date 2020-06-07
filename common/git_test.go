package common

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
)


func TestWriteGitShellFile(t *testing.T){

	temps := []string{"","testing"}

	_,err := WriteGitShellFile(temps[0]) 
	if err == nil || !strings.Contains(err.Error(),"WriteGitShellFile temp is null"){
		t.Error("test error")
	}

	_,err = WriteGitShellFile(temps[1])
	if err==nil{
		t.Error("test error")
	}
}


