package main

import (
    "fmt"
    "flag"
    "os"

    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
)

func main() {
    protectingTags := map[string]bool {"dont-kill": true, "production": true}
    killEnabled := flag.Bool("kill-enabled", false, "kill for real - default is dry run")
    probabilityString := flag.String("probability", "", "probability of a kill per ASG")

    flag.Parse()

    sess, err := session.NewSession()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error - problem creating session", err)
        os.Exit(1)
    }

    region := aws.Config{Region: aws.String("ap-southeast-2")}

    if *probabilityString == "" {
        fmt.Fprintf(os.Stderr, "Error - probability string not set\n")
        os.Exit(3)
    }

    probability, err := parseProbability(*probabilityString)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error - problem with probability string:", *probabilityString, "\n")
        os.Exit(2)
    }

    killProbability, err := assertProbability(probability)
    if err != nil {
        fmt.Fprintf(os.Stderr, "probability value assert fail\n")
        os.Exit(4)
    }

    fmt.Printf("kill enabled: %t, kill probability: %.2f%%\n", *killEnabled, killProbability*100)
    //fmt.Println("killProbability:", killProbability)

    initKiller()

    asgs := collateASGs(&region, sess)
    instances := identifyInstancesToKill(asgs, killProbability, &protectingTags)
    killInstances(&region, sess, instances, *killEnabled)
}
