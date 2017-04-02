Occasionally kill and AWS instance inside an AWS autoscale group. The probability of killing is a command line parameter. 

killa --help

enable a kill; default is off (ie dry run):

killa --kill-enabled

Set probability of a kill per autoscale group. Some different ways of representing 50% probability.

killa --probability=50%

killa --probability=50.1%

killa --probability=0.5

killa --probability=1/2


