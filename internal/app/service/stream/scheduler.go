package stream

var Sch chan SchMsg

type SchMsg struct {
	*Session
	*Command
}

func Scheduler(msg SchMsg) {
	switch msg.Command.Sign {
	case PlayCommand:
		msg.Session.DoPlay()
	}
}
