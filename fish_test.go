package subcli

import (
	"fmt"
	//"testing"
)

func ExampleBuildSub() {
	data := testingData()
	for _, v := range data.subs {
		fmt.Println(buildSub("testing", v))
	}
	// Unordered output:# test1
	// complete -c testing -n "__fish_use_subcommand -a test1 -d 'test1 desc'
	//
	// # test2
	// complete -c testing -n "__fish_use_subcommand -a test2 -d 'test2 desc'
	//
	// # test3
	// complete -c testing -n "__fish_use_subcommand -a test3 -d 'test3 desc'
	//
	// # test4
	// complete -c testing -n "__fish_use_subcommand -a test4 -d 'test4 desc'
	//
	// # version
	// complete -c testing -n "__fish_use_subcommand -a version -d 'prints testing version'
}
