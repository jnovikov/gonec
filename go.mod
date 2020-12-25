module github.com/jnovikov/gonec

go 1.12

require (
	github.com/boltdb/bolt v1.3.1
	github.com/covrom/decnum v0.0.0-20171025200506-7ae6a5b29284
	github.com/covrom/gonec v0.0.0-20190502071042-c145a972ff3a
	github.com/daviddengcn/go-colortext v0.0.0-20171126034257-17e75f6184bc
	github.com/dchest/siphash v1.2.0
	github.com/hashicorp/go-cleanhttp v0.5.0
	github.com/hashicorp/go-rootcerts v0.0.0-20160503143440-6bb64b370b90
	github.com/hashicorp/serf v0.8.1
	github.com/mattn/go-isatty v0.0.4
	github.com/mitchellh/go-homedir v1.0.0
	github.com/satori/go.uuid v1.2.0
	golang.org/x/sys v0.0.0-20181107165924-66b7b1311ac8
)

replace github.com/covrom/gonec => ./
