package main

import (
    "fmt"
    "math/rand"
    "time"

    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/autoscaling"
    "github.com/aws/aws-sdk-go/service/ec2"
)

func initKiller () {
    rand.Seed(time.Now().UnixNano())
}

func collateASGs (region *aws.Config, sess *session.Session) []*autoscaling.Group {

    svc := autoscaling.New(sess, region)
    params := &autoscaling.DescribeAutoScalingGroupsInput{
        AutoScalingGroupNames: []*string{
            aws.String("shawn-tester-asg"),
        },
        MaxRecords: aws.Int64(100),
    }

    resp, err := svc.DescribeAutoScalingGroups(params)

    if err != nil {
        fmt.Println(err.Error())
        return nil
    }

    return resp.AutoScalingGroups
}

func identifyInstancesToKill (groups []*autoscaling.Group, probability float64, protectingTags *map[string]bool) []*autoscaling.Instance {
    var instances []*autoscaling.Instance

    for _, group := range groups {
        protected := false

        for _, tag := range group.Tags {
            if (*protectingTags)[*tag.Key] {
                protected = true
                fmt.Println("tag:", *tag.Key, "protecting:", *group.AutoScalingGroupName)
                break
            }
        }

        if !protected {
            if doKill (probability) {
                instance := pickVictimInstance(group)
                if instance != nil {
                    fmt.Printf("Kill an instance in \"%s\": %s\n", *group.AutoScalingGroupName, *instance.InstanceId)
                    instances = append(instances, instance)
                } else {
                    fmt.Println("No instance to kill - likely 0 instances in ASG", *group.AutoScalingGroupName)
                }
            } else {
                fmt.Printf("Do not kill an instance in \"%s\"\n", *group.AutoScalingGroupName)
            }
        } else {
            fmt.Println(*group.AutoScalingGroupName, "protected")
        }
    }

    return instances
}

func killInstances (region *aws.Config, sess *session.Session, instances []*autoscaling.Instance, killEnabled bool) int {

    if len(instances) == 0 {
        return 0
    }

    svc := ec2.New(sess, region)

    ids := make([]*string, len(instances))
    var i int = 0

    for _, instance := range instances {
        ids[i] = instance.InstanceId
        i++
    }

    params := &ec2.TerminateInstancesInput{
        InstanceIds: ids,
        DryRun: aws.Bool(!killEnabled),
    }

    _, err := svc.TerminateInstances(params)

    if err != nil {
        fmt.Println(err.Error())
        return -1
    }

    return len(instances)
}

func pickVictimInstance (asg *autoscaling.Group) *autoscaling.Instance {
    if len(asg.Instances) > 0 {
        index := rand.Intn(len(asg.Instances))
        return asg.Instances[index]
    }
    return nil
}

func doKill (probability float64) bool {
    n := int(1000/probability)
    x := rand.Intn(n)
    //fmt.Printf("n: %d, random number: %d\n", n, x)
    return x < 1000
}

