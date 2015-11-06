Flappy Disk
===========

A flappy bird clone written in Go.

Building
--------

Currently I am doing my developemnt from Ubuntu Gnome 15.10. requrements to build are as follows:

```
$ go get -u github.com/jteeuwen/go-bindata/...
$ sudo apt-get install build-essential libegl1-mesa-dev libglfw3-dev libxrandr-dev libxcursor-dev libxinerama-dev libgl1-mesa-dev xorg-deUbuntu Install
```

To build the executable:

```
$ make
```

To make and run the executable:

```
$ make run
```

Credits
-------

[Original OpenGL code](https://github.com/go-gl/examples/tree/master/glfw31-gl41core-cube)

