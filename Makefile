sub:
	git submodule add --force -b main git@github.com:hkloudou/mqx.git github.com/hkloudou/mqx
run:
	cd app/broker && make run