// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package term

import "syscall"

const (
	// Terminal attribute types.
	sshIflag = iota
	sshOflag
	sshCflag
	sshLflag
	sshCchar
	sshTspeed
	sshNOP

	// SSH terminal attributes.
	sshTTYOPEND    = 0
	sshVINTR       = 1
	sshVQUIT       = 2
	sshVERASE      = 3
	sshVKILL       = 4
	sshVEOF        = 5
	sshVEOL        = 6
	sshVEOL2       = 7
	sshVSTART      = 8
	sshVSTOP       = 9
	sshVSUSP       = 10
	sshVDSUSP      = 11
	sshVREPRINT    = 12
	sshVWERASE     = 13
	sshVLNEXT      = 14
	sshVFLUSH      = 15
	sshVSWTCH      = 16
	sshVSTATUS     = 17
	sshVDISCARD    = 18
	sshIGNPAR      = 30
	sshPARMRK      = 31
	sshINPCK       = 32
	sshISTRIP      = 33
	sshINLCR       = 34
	sshIGNCR       = 35
	sshICRNL       = 36
	sshIUCLC       = 37
	sshIXON        = 38
	sshIXANY       = 39
	sshIXOFF       = 40
	sshIMAXBEL     = 41
	sshISIG        = 50
	sshICANON      = 51
	sshXCASE       = 52
	sshECHO        = 53
	sshECHOE       = 54
	sshECHOK       = 55
	sshECHONL      = 56
	sshNOFLSH      = 57
	sshTOSTOP      = 58
	sshIEXTEN      = 59
	sshECHOCTL     = 60
	sshECHOKE      = 61
	sshPENDIN      = 62
	sshOPOST       = 70
	sshOLCUC       = 71
	sshONLCR       = 72
	sshOCRNL       = 73
	sshONOCR       = 74
	sshONLRET      = 75
	sshCS7         = 90
	sshCS8         = 91
	sshPARENB      = 92
	sshPARODD      = 93
	sshTTYOPISPEED = 128
	sshTTYOPOSPEED = 129
)

var convertSSH = map[uint8]struct {
	tType  uint
	native uint32
}{
	sshTTYOPEND:    {tType: sshNOP},
	sshVINTR:       {tType: sshCchar, native: syscall.VINTR},
	sshVQUIT:       {tType: sshCchar, native: syscall.VQUIT},
	sshVERASE:      {tType: sshCchar, native: syscall.VERASE},
	sshVKILL:       {tType: sshCchar, native: syscall.VKILL},
	sshVEOF:        {tType: sshCchar, native: syscall.VEOF},
	sshVEOL:        {tType: sshCchar, native: syscall.VEOL},
	sshVEOL2:       {tType: sshCchar, native: syscall.VEOL2},
	sshVSTART:      {tType: sshCchar, native: syscall.VSTART},
	sshVSTOP:       {tType: sshCchar, native: syscall.VSTOP},
	sshVSUSP:       {tType: sshCchar, native: syscall.VSUSP},
	sshVDSUSP:      {tType: sshCchar, native: sshNOP},
	sshVREPRINT:    {tType: sshCchar, native: syscall.VREPRINT},
	sshVWERASE:     {tType: sshCchar, native: syscall.VWERASE},
	sshVLNEXT:      {tType: sshCchar, native: syscall.VLNEXT},
	sshVFLUSH:      {tType: sshNOP},
	sshVSWTCH:      {tType: sshCchar, native: syscall.VSWTC},
	sshVSTATUS:     {tType: sshNOP},
	sshVDISCARD:    {tType: sshCchar, native: syscall.VDISCARD},
	sshIGNPAR:      {tType: sshIflag, native: syscall.IGNPAR},
	sshPARMRK:      {tType: sshIflag, native: syscall.PARMRK},
	sshINPCK:       {tType: sshIflag, native: syscall.INPCK},
	sshISTRIP:      {tType: sshIflag, native: syscall.ISTRIP},
	sshINLCR:       {tType: sshIflag, native: syscall.INLCR},
	sshIGNCR:       {tType: sshIflag, native: syscall.IGNCR},
	sshICRNL:       {tType: sshIflag, native: syscall.ICRNL},
	sshIUCLC:       {tType: sshIflag, native: syscall.IUCLC},
	sshIXON:        {tType: sshIflag, native: syscall.IXON},
	sshIXANY:       {tType: sshIflag, native: syscall.IXANY},
	sshIXOFF:       {tType: sshIflag, native: syscall.IXOFF},
	sshIMAXBEL:     {tType: sshIflag, native: syscall.IMAXBEL},
	sshISIG:        {tType: sshLflag, native: syscall.ISIG},
	sshICANON:      {tType: sshLflag, native: syscall.ICANON},
	sshXCASE:       {tType: sshLflag, native: syscall.XCASE},
	sshECHO:        {tType: sshLflag, native: syscall.ECHO},
	sshECHOE:       {tType: sshLflag, native: syscall.ECHOE},
	sshECHOK:       {tType: sshLflag, native: syscall.ECHOK},
	sshECHONL:      {tType: sshLflag, native: syscall.ECHONL},
	sshNOFLSH:      {tType: sshLflag, native: syscall.NOFLSH},
	sshTOSTOP:      {tType: sshLflag, native: syscall.TOSTOP},
	sshIEXTEN:      {tType: sshLflag, native: syscall.IEXTEN},
	sshECHOCTL:     {tType: sshLflag, native: syscall.ECHOCTL},
	sshECHOKE:      {tType: sshLflag, native: syscall.ECHOKE},
	sshPENDIN:      {tType: sshNOP},
	sshOPOST:       {tType: sshOflag, native: syscall.OPOST},
	sshOLCUC:       {tType: sshOflag, native: syscall.OLCUC},
	sshONLCR:       {tType: sshOflag, native: syscall.ONLCR},
	sshOCRNL:       {tType: sshOflag, native: syscall.OCRNL},
	sshONOCR:       {tType: sshOflag, native: syscall.ONOCR},
	sshONLRET:      {tType: sshOflag, native: syscall.ONLRET},
	sshCS7:         {tType: sshCflag, native: syscall.CS7},
	sshCS8:         {tType: sshCflag, native: syscall.CS8},
	sshPARENB:      {tType: sshCflag, native: syscall.PARENB},
	sshPARODD:      {tType: sshCflag, native: syscall.PARODD},
	sshTTYOPISPEED: {tType: sshTspeed},
	sshTTYOPOSPEED: {tType: sshTspeed},
}

// ToSSH converts the Termios attributes to SSH attributes usable as syscall.TerminalModes.
func (t *Termios) ToSSH() map[uint8]uint32 {
	sshModes := make(map[uint8]uint32, len(convertSSH))
	var flags uint32
	for sshID, tios := range convertSSH {
		switch tios.tType {
		case sshIflag:
			flags = t.Iflag
		case sshOflag:
			flags = t.Oflag
		case sshLflag:
			flags = t.Lflag
		case sshCflag:
			flags = t.Cflag
		case sshCchar:
			sshModes[sshID] = uint32(t.Cc[tios.native])
			continue
		case sshTspeed:
			sshModes[sshTTYOPISPEED], sshModes[sshTTYOPOSPEED] = t.Ispeed, t.Ospeed
			continue
		default:
			continue
		}
		var onOff uint32
		if tios.native&flags > 0 {
			onOff = 1
		}
		sshModes[sshID] = onOff
	}
	return sshModes
}

// FromSSH converts SSH attributes to Termios attributes.
func (t *Termios) FromSSH(termModes map[uint8]uint32) {
	var flags *uint32
	for sshID, val := range termModes {
		switch convertSSH[sshID].tType {
		case sshIflag:
			flags = &t.Iflag
		case sshOflag:
			flags = &t.Oflag
		case sshLflag:
			flags = &t.Lflag
		case sshCflag:
			flags = &t.Cflag
		case sshCchar:
			t.Cc[convertSSH[sshID].native] = byte(val)
			continue
		case sshTspeed:
			if sshID == sshTTYOPISPEED {
				t.Ispeed = val
			} else {
				t.Ospeed = val
			}
			continue
		default:
			continue
		}
		if val > 0 {
			*flags |= convertSSH[sshID].native
		} else {
			*flags &^= convertSSH[sshID].native
		}
	}
}
