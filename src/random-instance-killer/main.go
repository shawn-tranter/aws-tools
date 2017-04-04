package main

import (
    "fmt"
    "flag"
    "os"
    "reflect"

    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
)

func main() {
    var protectingTagsArray stringArrayFlags
    killEnabled := flag.Bool("kill-enabled", false, "kill for real - default is dry run")
    probabilityString := flag.String("probability", "", "probability of a kill per ASG")
    flag.Var(&protectingTagsArray, "protect-tag", "ASG tags")

    protectingTags := make(map[string]bool)

    flag.Parse()

    for _, tag := range protectingTagsArray {
        protectingTags[tag] = true
    }

    fmt.Println(protectingTagsArray)
    fmt.Println(protectingTags)

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

    fmt.Printf("kill enabled: %t, kill probability: %.2f%%, protecting tags: %v\n", *killEnabled, killProbability*100, reflect.ValueOf(protectingTags).MapKeys())

    initKiller()

    asgs := collateASGs(&region, sess)
    instances := identifyInstancesToKill(asgs, killProbability, &protectingTags)
    killInstances(&region, sess, instances, *killEnabled)
}

type stringArrayFlags []string

func (i *stringArrayFlags) String() string {
    return "my string representation"
}

func (i *stringArrayFlags) Set(value string) error {
    *i = append(*i, value)
    return nil
}