# KeyLogger
A port of the C++ keylogger github.com/kernc/logkeys to Golang using the /dev/input/event* devices to read key presses outside of X11.  My next goal is
to combine with an X11 method of grabbing keystrokes in case we don't have root privileges.

# Example Usage

    LogKeys( os.Stdout )

That's it!  I tried to comment the code as best I could since it was a learning experience for me, hopefully it's useful to others as well.

# License

This is FREE SOFTWARE with absolutely NO WARRANTY of any kind.  Feel free to copy,
fork, or use the code how you like.
