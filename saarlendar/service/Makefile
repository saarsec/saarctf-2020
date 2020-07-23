
default:	build

clean:
	$(MAKE) -C ./src/nginx-1.17.4 clean
	$(MAKE) -C ./src/duktape clean
	rm ./Makefile
	rm ./nginx

build:
	$(MAKE) -C ./src/duktape
	$(MAKE) -C ./src/nginx-1.17.4
	cp ./src/nginx-1.17.4/objs/nginx ./nginx

