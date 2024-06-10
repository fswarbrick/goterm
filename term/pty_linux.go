package term

import (
	"os"
	"strconv"
	"syscall"
	"unsafe"
)

const (
	TIOCGPTN   = 0x80045430
	TIOCSPTLCK = 0x40045431
)

// PTSName return the name of the pty.
func (p *PTY) PTSName() (string, error) {
	n, err := p.PTSNumber()
	if err != nil {
		return "", err
	}
	return "/dev/pts/" + strconv.Itoa(int(n)), nil
}

// PTSNumber return the pty number.
func (p *PTY) PTSNumber() (uint, error) {
	var ptyno uint
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(p.Master.Fd()), uintptr(TIOCGPTN), uintptr(unsafe.Pointer(&ptyno)))
	if errno != 0 {
		return 0, errno
	}
	return ptyno, nil
}

func (p *PTY) PTSUnlock() error {
	// unlock pty slave
	var unlock int // 0 => Unlock
	if _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(p.Master.Fd()), uintptr(TIOCSPTLCK), uintptr(unsafe.Pointer(&unlock))); errno != 0 {
		p.Master.Close()
		return errno
	}
	return nil
}

// OpenPTY Creates a new Master/Slave PTY pair.
func OpenPTY() (*PTY, error) {
	// Opening ptmx gives you the FD of a brand new PTY
	master, err := os.OpenFile("/dev/ptmx", os.O_RDWR, 0)
	if err != nil {
		return nil, err
	}
	pty := &PTY{Master: master}

	err = pty.PTSUnlock()
	if err != nil {
		master.Close()
		return nil, err
	}

	// get path of pts slave
	slaveStr, err := pty.PTSName()
	if err != nil {
		master.Close()
		return nil, err
	}

	// open pty slave
	pty.Slave, err = os.OpenFile(slaveStr, os.O_RDWR|syscall.O_NOCTTY, 0)
	if err != nil {
		master.Close()
		return nil, err
	}

	return pty, nil
}
