Flappy Disk
===========

A flappy bird clone written in Go for Extra-Life 2016.

If you think this is a cool idea, [DONATE!](http://www.extra-life.org/participant/hurricanerix).  It is #FORTHEKIDS after all.

Extra-Life Game Day 2015 has come and gone, but I am not stopping.  I am going to continue fundraising until Dec 31st 2015.  During that time, I will continue working on the game.  I pledge that by Dec 31st, the game will be a fully functional and polished that anybody would be proud to play at Extra-Life 2016 (probably at the end, when things like Solitaire come out).

[First commit from Extra-Life 2016 event](https://github.com/hurricanerix/FlappyDisk/commit/8b0f5916fee5fc910fab048135529ee5e3573173)

[Last commit from Extra-Life 2016 event](https://github.com/hurricanerix/FlappyDisk/commit/0628e4290990e309a8af25286fab8849ee0f435f)

Additionally everybody who donates will get their name into the credits screen of the game.

![](https://github.com/hurricanerix/FlappyDisk/blob/master/screenshot.png)

Building
--------

Developemnt is being done on OS X El Capitan.

```
# Install Xcode
$ brew install homebrew/versions/glfw3
$ go get -u github.com/jteeuwen/go-bindata/...
$ make gen
$ go get
```

For Debian based Linux, this might work:

```
$ go get -u github.com/jteeuwen/go-bindata/...
$ sudo apt-get install build-essential libegl1-mesa-dev libglfw3-dev libxrandr-dev libxcursor-dev libxinerama-dev libgl1-mesa-dev xorg-deUbuntu Install
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
