package applicationstate

func (m Model) View() string {
	switch m.programState {
	case editTimerState:
		return m.inputView()

	case sessionRunningState:
		return m.runningSessionView()

	case sessionCompleteState:
		return m.completedView()

	case sessionEndedEarlyState:
		return m.completedView()

	case breakState:
		return m.breakView()

	default:
		return m.completedView()
	}
}
