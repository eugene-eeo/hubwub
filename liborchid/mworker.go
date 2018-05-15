package liborchid

const (
	PlaybackStart = iota
	PlaybackEnd
)

type PlaybackResult struct {
	State    int
	Song     *Song
	Stream   *Stream
	Complete bool
	Error    error
}

type MWorker struct {
	Results   chan *PlaybackResult
	SongQueue chan *Song
	stop      chan struct{}
}

func NewMWorker() *MWorker {
	return &MWorker{
		Results:   make(chan *PlaybackResult),
		SongQueue: make(chan *Song),
		stop:      make(chan struct{}),
	}
}

func (mw *MWorker) report(state int, song *Song, stream *Stream, complete bool, err error) {
	mw.Results <- &PlaybackResult{
		State:    state,
		Song:     song,
		Stream:   stream,
		Complete: complete,
		Error:    err,
	}
}

func (mw *MWorker) Stop() {
	mw.stop <- struct{}{}
}

func (mw *MWorker) Play() {
loop:
	for {
		select {
		case song := <-mw.SongQueue:
			stream, err := song.Stream()
			if err != nil {
				mw.report(PlaybackEnd, song, nil, false, err)
				break
			}
			stream.Play()
			mw.report(PlaybackStart, song, stream, false, nil)
			go func() {
				mw.report(PlaybackEnd, song, stream, <-stream.Complete(), nil)
			}()
		case <-mw.stop:
			mw.Results <- nil
			close(mw.Results)
			break loop
		}
	}
}