# AWS Tools

## random-instance-killer

Occasionally kill and AWS instance inside an AWS autoscale group. The probability of killing is a command line parameter. 

```random-instance-killer --help```

### enable killing

* default is off (ie dry run):

```random-instance-killer --kill-enabled```

### probability

Set probability of a kill per autoscale group. Some different ways of representing 50% probability.

```random-instance-killer --probability=50%```

```random-instance-killer --probability=50.1%```

```random-instance-killer --probability=0.5```

```random-instance-killer --probability=1/2```

### protect-tag

Any ASG tagged with a protect-tag will not have an instance deleted.

```random-instance-killer --probability=50% --protect-tag=dont-kill --protect-tag=production```


