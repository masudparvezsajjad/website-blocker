package reflectpause

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"time"

	"golang.org/x/term"
)

// ErrAborted means the user chose to cancel during the reflective pause.
var ErrAborted = errors.New("cancelled")

const defaultSeconds = 15

// habitPrompts mixes short reflective questions and habit-focused quotes.
var habitPrompts = []string{
	"Will this choice match the person you want to become?",
	"Is this impulse, or the decision your calmer self would make?",
	"What are you trading away for the next few minutes of comfort?",
	"Habits are the compound interest of self-improvement. — James Clear",
	"You do not rise to the level of your goals; you fall to the level of your systems.",
	"Small daily improvements are the hidden engine of real change.",
	"The chains of habit are too light to be felt until they are too heavy to be broken. — Burke",
	"Discipline is choosing what you want most over what you want now.",
	"Every time you pause instead of react, you vote for a stronger habit.",
	"Your future self is watching; make them proud.",
	"We are what we repeatedly do. Excellence, then, is not an act, but a habit.",
	"One clear boundary today saves a hundred regrets tomorrow.",
	"Comfort now often costs clarity later—take this second to notice which you want.",
}

// followUpQuestions are shown after the countdown; two or three are picked at random.
var followUpQuestions = []string{
	"Why do you want to take this step right now?",
	"Is this really necessary—or is it the easier path in the moment?",
	"In an hour, will you be glad you weakened this boundary?",
	"What need are you trying to meet, and is there a healthier way?",
	"Are you avoiding discomfort, or making a deliberate choice?",
	"Would you recommend this same decision to someone you care about?",
	"Does this match the rules you chose for yourself when you were thinking clearly?",
	"What story will you tell yourself about this choice tomorrow?",
	"Is this a one-time exception, or the start of a new habit?",
	"What would the version of you that installed this tool say right now?",
}

// Run runs two interactive phases: a countdown with a random quote, then random follow-up
// questions before final confirmation. Press x or X anytime to cancel; press Enter after phase 2
// to continue. If stdin is not a terminal, Run returns immediately without pausing.
func Run(seconds int, whatYouAreDoing string) error {
	if seconds <= 0 {
		seconds = defaultSeconds
	}

	fd := int(os.Stdin.Fd())
	if !term.IsTerminal(fd) {
		return nil
	}

	old, err := term.MakeRaw(fd)
	if err != nil {
		return fmt.Errorf("reflect pause: %w", err)
	}
	defer func() { _ = term.Restore(fd, old) }()

	inputCh := make(chan byte, 32)
	go readKeysToChannel(os.Stdin, inputCh)

	fmt.Fprintf(os.Stdout, "\nPausing before you %s\n\n", whatYouAreDoing)
	fmt.Fprintln(os.Stdout, habitPrompts[rand.Intn(len(habitPrompts))])
	fmt.Fprintln(os.Stdout)
	fmt.Fprintln(os.Stdout, "Press x at any time to cancel.")

	for left := seconds; left > 0; left-- {
		fmt.Fprintf(os.Stdout, "\r  %d seconds… ", left)
		_ = os.Stdout.Sync()
		if quit := waitSecondOrQuit(inputCh); quit {
			fmt.Fprintln(os.Stdout)
			fmt.Fprintln(os.Stdout, "Cancelled.")
			return ErrAborted
		}
	}
	fmt.Fprintln(os.Stdout)

	n := 2 + rand.Intn(2)
	perm := rand.Perm(len(followUpQuestions))
	fmt.Fprintln(os.Stdout, "Before you continue:")
	fmt.Fprintln(os.Stdout)
	for i := 0; i < n; i++ {
		fmt.Fprintf(os.Stdout, "  • %s\n", followUpQuestions[perm[i]])
	}
	fmt.Fprintln(os.Stdout)
	fmt.Fprintln(os.Stdout, "Press Enter to continue, or x to cancel.")

	for {
		b := <-inputCh
		switch b {
		case 'x', 'X':
			fmt.Fprintln(os.Stdout, "Cancelled.")
			return ErrAborted
		case '\r', '\n':
			return nil
		}
	}
}

func waitSecondOrQuit(inputCh <-chan byte) bool {
	t := time.NewTimer(time.Second)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			return false
		case b := <-inputCh:
			if b == 'x' || b == 'X' {
				return true
			}
		}
	}
}

func readKeysToChannel(r *os.File, ch chan<- byte) {
	buf := make([]byte, 1)
	for {
		n, err := r.Read(buf)
		if err != nil || n == 0 {
			return
		}
		select {
		case ch <- buf[0]:
		default:
		}
	}
}
