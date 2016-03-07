NAME=deliveryagent
SOURCE=github.com/rickbau5/
EXE=bin/

all: deliveryagent

deliveryagent:
	go install $(SOURCE)$(NAME)

install: deliveryagent
	mv $(EXE)$(NAME) /etc/init.d/$(NAME)
