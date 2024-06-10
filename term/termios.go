// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package term implements a subset of the C termios library to interface with Terminals.

This package allows the caller to get and set most Terminal capabilites
and sizes as well as create PTYs to enable writing things like script,
screen, tmux, and expect.

The Termios type is used for setting/getting Terminal capabilities while
the PTY type is used for handling virtual terminals.

Currently this part of this lib is Linux specific.

Also implements a simple version of readline in pure Go and some Stringers
for terminal colors and attributes.
*/

package term

import (
	"errors"
	"os"
	"strings"
	"syscall"
	"unsafe"

	"golang.org/x/sys/unix"
)

const (
	tNCCS = 32 // tNCCS    Termios CC size
)

// Termios merge of the C Terminal and Kernel termios structs.
type Termios struct {
	Iflag  uint32      // Iflag Handles the different Input modes
	Oflag  uint32      // Oflag For the different Output modes
	Cflag  uint32      // Cflag Control modes
	Lflag  uint32      // Lflag Local modes
	Line   byte        // Line
	Cc     [tNCCS]byte // Cc Control characters. How to handle special Characters eg. Backspace being ^H or ^? and so on
	Ispeed uint32      // Ispeed Hardly ever used speed of terminal
	Ospeed uint32      // Ospeed "
	Wz     Winsize     // Wz Terminal size information.
}

// Winsize handle the terminal window size.
type Winsize struct {
	WsRow    uint16 // WsRow 		Terminal number of rows
	WsCol    uint16 // WsCol 		Terminal number of columns
	WsXpixel uint16 // WsXpixel Terminal width in pixels
	WsYpixel uint16 // WsYpixel Terminal height in pixels
}

// PTY the PTY Master/Slave are always bundled together so makes sense to bundle here too.
//
// Slave  - implements the virtual terminal functionality and the place you connect client applications
// Master - Things written to the Master are forwarded to the Slave terminal and the other way around.
//
//	This gives reading from Master would give you nice line-by-line with no strange characters in
//	Cooked() Mode and every char in Raw() mode.
//
// Since Slave is a virtual terminal it depends on the terminal settings ( in this lib the Termios ) what
// and when data is forwarded through the terminal.
//
// See 'man pty' for further info
type PTY struct {
	Master *os.File // Master The Master part of the PTY
	Slave  *os.File // Slave The Slave part of the PTY
}

// Raw Sets terminal t to raw mode.
// This gives that the terminal will do the absolut minimal of processing, pretty much send everything through.
// This is normally what Shells and such want since they have their own readline and movement code.
func (t *Termios) Raw() {
	t.Iflag &^= syscall.IGNBRK | syscall.BRKINT | syscall.PARMRK | syscall.ISTRIP | syscall.INLCR | syscall.IGNCR | syscall.ICRNL | syscall.IXON
	// t.Iflag &^= BRKINT | ISTRIP | ICRNL | IXON // Stevens RAW
	t.Oflag &^= syscall.OPOST
	t.Lflag &^= syscall.ECHO | syscall.ECHONL | syscall.ICANON | syscall.ISIG | syscall.IEXTEN
	t.Cflag &^= syscall.CSIZE | syscall.PARENB
	t.Cflag |= syscall.CS8
	t.Cc[syscall.VMIN] = 1
	t.Cc[syscall.VTIME] = 0
}

// Cook Set the Terminal to Cooked mode.
// In this mode the Terminal process the information before sending it on to the application.
func (t *Termios) Cook() {
	t.Iflag |= syscall.BRKINT | syscall.IGNPAR | syscall.ISTRIP | syscall.ICRNL | syscall.IXON
	t.Oflag |= syscall.OPOST
	t.Lflag |= syscall.ISIG | syscall.ICANON
}

// Sane reset Term to sane values.
// Should be pretty much what the shell command "reset" does to the terminal.
func (t *Termios) Sane() {
	t.Iflag &^= syscall.IGNBRK | syscall.INLCR | syscall.IGNCR | syscall.IUTF8 | syscall.IXOFF | syscall.IUCLC | syscall.IXANY
	t.Iflag |= syscall.BRKINT | syscall.ICRNL | syscall.IMAXBEL
	t.Oflag |= syscall.OPOST | syscall.ONLCR
	t.Oflag &^= syscall.OLCUC | syscall.OCRNL | syscall.ONOCR | syscall.ONLRET
	t.Cflag |= syscall.CREAD
}

// Set Sets terminal t attributes on file.
func (t *Termios) Set(file *os.File) error {
	fd := file.Fd()
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), uintptr(syscall.TCSETS), uintptr(unsafe.Pointer(t)))
	if errno != 0 {
		return errno
	}
	return nil
}

// Attr Gets (terminal related) attributes from file.
func Attr(file *os.File) (Termios, error) {
	var t Termios
	fd := file.Fd()
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), uintptr(syscall.TCGETS), uintptr(unsafe.Pointer(&t)))
	if errno != 0 {
		return t, errno
	}
	t.Ispeed &= unix.CBAUD | unix.CBAUDEX
	t.Ospeed &= unix.CBAUD | unix.CBAUDEX
	return t, nil
}

// Isatty returns true if file is a tty.
func Isatty(file *os.File) bool {
	_, err := Attr(file)
	return err == nil
}

// GetPass reads password from a TTY with no echo.
func GetPass(prompt string, f *os.File, pbuf []byte) ([]byte, error) {
	t, err := Attr(f)
	if err != nil {
		return nil, err
	}
	defer t.Set(f)
	noecho := t
	noecho.Lflag = noecho.Lflag &^ syscall.ECHO
	if err := noecho.Set(f); err != nil {
		return nil, err
	}
	b := make([]byte, 1, 1)
	i := 0
	if _, err := f.Write([]byte(prompt)); err != nil {
		return nil, err
	}
	for ; i < len(pbuf); i++ {
		if _, err := f.Read(b); err != nil {
			b[0] = 0
			clearbuf(pbuf[:i+1])
		}
		if b[0] == '\n' || b[0] == '\r' {
			return pbuf[:i], nil
		}
		pbuf[i] = b[0]
		b[0] = 0
	}
	clearbuf(pbuf[:i+1])
	return nil, errors.New("ran out of bufferspace")
}

// clearbuf clears out the buffer incase we couldn't read the full password.
func clearbuf(b []byte) {
	for i := range b {
		b[i] = 0
	}
}

// GetChar reads a single byte.
func GetChar(f *os.File) (b byte, err error) {
	bs := make([]byte, 1, 1)
	if _, err = f.Read(bs); err != nil {
		return 0, err
	}
	return bs[0], err
}

// Winsz Fetches the current terminal windowsize.
// example handling changing window sizes with PTYs:
//
// import "os"
// import "os/signal"
//
// var sig = make(chan os.Signal,2) 		// Channel to listen for UNIX SIGNALS on
// signal.Notify(sig, syscall.SIGWINCH) // That'd be the window changing
//
//	for {
//		<-sig
//		term.Winsz(os.Stdin)			// We got signaled our terminal changed size so we read in the new value
//	 term.Setwinsz(pty.Slave) // Copy it to our virtual Terminal
//	}
func (t *Termios) Winsz(file *os.File) error {
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(file.Fd()), uintptr(syscall.TIOCGWINSZ), uintptr(unsafe.Pointer(&t.Wz)))
	if errno != 0 {
		return errno
	}
	return nil
}

// Setwinsz Sets the terminal window size.
func (t *Termios) Setwinsz(file *os.File) error {
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(file.Fd()), uintptr(syscall.TIOCSWINSZ), uintptr(unsafe.Pointer(&t.Wz)))
	if errno != 0 {
		return errno
	}
	return nil
}

// Close closes the PTYs that OpenPTY created.
func (p *PTY) Close() error {
	slaveErr := errors.New("Slave FD nil")
	if p.Slave != nil {
		slaveErr = p.Slave.Close()
	}
	masterErr := errors.New("Master FD nil")
	if p.Master != nil {
		masterErr = p.Master.Close()
	}
	if slaveErr != nil || masterErr != nil {
		var errs []string
		if slaveErr != nil {
			errs = append(errs, "Slave: "+slaveErr.Error())
		}
		if masterErr != nil {
			errs = append(errs, "Master: "+masterErr.Error())
		}
		return errors.New(strings.Join(errs, " "))
	}
	return nil
}

// ReadByte implements the io.ByteReader interface to read single char from the PTY.
func (p *PTY) ReadByte() (byte, error) {
	bs := make([]byte, 1, 1)
	_, err := p.Master.Read(bs)
	return bs[0], err
}

// GetChar fine old getchar() for a PTY.
func (p *PTY) GetChar() (byte, error) {
	return p.ReadByte()
}
