Flappy Disk
===========

A flappy bird clone written in Go for [Extra-Life](http://www.extra-life.org/participant/hurricanerix).

Building
--------

Currently I am doing my developemnt from Ubuntu Gnome 15.10. requrements to build are as follows:

```
$ go get -u github.com/jteeuwen/go-bindata/...
$ sudo apt-get install build-essential libegl1-mesa-dev libglfw3-dev libxrandr-dev libxcursor-dev libxinerama-dev libgl1-mesa-dev xorg-deUbuntu Install
$ go get
```

For OSX

```
# Install Xcode
$ brew install homebrew/versions/glfw3
$ go get
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

[*.pyxel assets created with Pyxel Edit](http://pyxeledit.com/)

[Original OpenGL code](https://github.com/go-gl/examples/tree/master/glfw31-gl41core-cube)
