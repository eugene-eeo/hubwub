package main

import "time"
import "math/rand"
import "github.com/nsf/termbox-go"
import "github.com/eugene-eeo/orchid/liborchid"

type request func(*hub)

type hub struct {
	Player     *liborchid.Player
	Stream     *liborchid.Stream
	Song       *liborchid.Song
	done       chan struct{}
	Requests   chan request
	playerView *playerView
}

func (h *hub) Paused() bool {
	return h.Stream == nil || h.Stream.Paused()
}

func (h *hub) Toggle() {
	if h.Stream != nil {
		h.Stream.Toggle()
	}
}

func newHub(p *liborchid.Player) *hub {
	h := &hub{
		Player:     p,
		Stream:     nil,
		Requests:   make(chan request),
		done:       make(chan struct{}),
		playerView: newPlayerView(),
	}
	return h
}

func (h *hub) Render() {
	h.playerView.Update(
		h.Player,
		h.Paused(),
		h.Player.Shuffle,
		h.Player.Repeat,
	)
}

func (h *hub) Play() {
	h.Song = h.Player.Song()
	if h.Song == nil {
		return
	}
	stream, err := h.Song.Stream()
	if err != nil {
		h.Stream = nil
		h.Player.Remove()
		go func() {
			h.Requests <- func(c *hub) { c.Play() }
		}()
		return
	}
	h.Stream = stream
	stream.Play()
	go func() {
		complete := <-stream.Complete()
		if complete {
			h.Requests <- func(c *hub) {
				c.Player.Next(1, false)
				c.Play()
			}
		}
	}()
}

func (h *hub) handle(evt termbox.Event) {
	if evt.Type != termbox.EventKey {
		return
	}
	switch evt.Ch {
	case 'q':
		h.done <- struct{}{}
	case 'n':
		h.Player.Next(1, true)
		h.Play()
	case 'p':
		h.Player.Next(-1, true)
		h.Play()
	case 's':
		h.Player.ToggleShuffle()
	case 'r':
		h.Player.ToggleRepeat()
	case 'f':
		f := newFinderUIFromPlayer(h.Player)
		go f.Loop()
		song := <-f.choice
		if song != nil {
			h.Player.SetCurrent(song)
			h.Play()
		}
	}
	if evt.Key == termbox.KeySpace {
		h.Toggle()
	}
}

func (h *hub) Loop(events <-chan termbox.Event) {
	for {
		select {
		case evt := <-events:
			h.handle(evt)
			h.Render()
		case req := <-h.Requests:
			req(h)
			h.Render()
		}
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	songs := liborchid.FindSongs(".")

	must(termbox.Init())
	termbox.SetOutputMode(termbox.Output256)
	defer termbox.Close()

	h := newHub(liborchid.NewPlayer(songs))
	events := make(chan termbox.Event)
	go func() {
		for {
			events <- termbox.PollEvent()
		}
	}()
	go h.Loop(events)

	h.Requests <- func(h *hub) {
		h.Play()
	}
	<-h.done
}
