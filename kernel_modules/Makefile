CONFIG_STACK_VALIDATION=

ccflags-y  :=  -I/usr/include -I../bpf/headers

obj-m += macvlan.o

all:
	make -C /lib/modules/$(shell uname -r)/build M=$(PWD) modules

clean:
	make -C /lib/modules/$(shell uname -r)/build M=$(PWD) clean
