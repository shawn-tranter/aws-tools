package main

import (
    "fmt"
    "flag"

    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
)

func main() {
    protectingTags := map[string]bool {"dont-kill": true, "production": true}
    killEnabledPtr := flag.Bool("kill-enabled", false, "kill for real - default is dry run")
    probabilityPtr := flag.String("probability", "0", "probability of a kill per ASG")

    flag.Parse()

    sess, err := session.NewSession()
    if err != nil {
        panic(err)
    }

    region := aws.Config{Region: aws.String("ap-southeast-2")}

    killEnabled := *killEnabledPtr

    probability, err := parseProbability(*probabilityPtr)
    if err != nil {
        fmt.Println("probability string:", *probabilityPtr)
        panic(err)
    }

    killProbability, err := assertProbability(probability)
    if err != nil {
        fmt.Println("probability value assert fail")
        panic(err)
    }

    fmt.Printf("kill enabled: %t, kill probability: %2f\n", killEnabled, killProbability)
    //fmt.Println("killProbability:", killProbability)

    initKiller()

    asgs := collateASGs(&region, sess)
    instances := identifyInstancesToKill(asgs, killProbability, &protectingTags)
    killInstances(&region, sess, instances, killEnabled)
}
