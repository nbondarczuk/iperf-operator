DOMAIN = nbt.pl
REPO = gitgub.com/nbondarczuk/iperf-operator
VERSION = v1alpha1

all: init create

init:
	operator-sdk init --domain $(DOMAIN) --repo $(REPO) 

create:
	operator-sdk create api --group operator --version $(VERSION) --kind IPerfOperator --resource --controller
