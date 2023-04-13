package main

import (
	"flag"
	"fmt"
	"time"

	cronTask "github.com/erbaner/be/internal/cron_task"
)

func main() {
	var userID = flag.String("userID", "", "userID to clear msg and reset seq")
	var workingGroupID = flag.String("workingGroupID", "", "workingGroupID to clear msg and reset seq")
	flag.Parse()
	fmt.Println(time.Now(), "start cronTask", *userID, *workingGroupID)
	cronTask.StartCronTask(*userID, *workingGroupID)
}
