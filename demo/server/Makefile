all:
	c++ -I. -I../demo/ice_interface/ -c ../demo/ice_interface/Printer.cpp Server.cpp
	c++ -o server Printer.o Server.o -lIce -lIceUtil -lpthread
	rm *.o
