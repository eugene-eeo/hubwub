package main

import "strings"
import "sort"
import "github.com/eugene-eeo/orchid/player"
import "github.com/nsf/termbox-go"

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

type item struct {
	str string
	idx int
}

func match(a, b string) bool {
	B := []rune(b)
	n := len(B)
	j := 0
	for _, c := range a {
		found := false
		for j < n && !found {
			found = c == B[j]
			j++
		}
		if !found {
			return false
		}
	}
	return true
}

func rankMatch(a, b string) float64 {
	return float64(len(a)) / float64(len(b))
}

func matchAll(query string, haystack []*item) []*item {
	matching := []*item{}
	for _, x := range haystack {
		if match(query, x.str) {
			matching = append(matching, x)
		}
	}
	sort.Slice(matching, func(i, j int) bool {
		return rankMatch(query, matching[i].str) < rankMatch(query, matching[j].str)
	})
	return matching
}

type Finder struct {
	items []*item
	songs []player.Song
}

func finderFromPlayer(p *player.Player) *Finder {
	all := p.Songs()
	songs := make([]player.Song, len(all))
	items := make([]*item, len(all))
	for i, song := range all {
		songs[i] = song
		items[i] = &item{
			str: strings.ToLower(song.Name()),
			idx: i,
		}
	}
	return &Finder{
		items: items,
		songs: songs,
	}
}

func (f *Finder) Find(q string) []*item {
	return matchAll(strings.ToLower(q), f.items)
}

func (f *Finder) Get(i *item) player.Song {
	return f.songs[i.idx]
}

type FinderUI struct {
	finder   *Finder
	requests chan func(*FinderUI)
	choice   chan *player.Song
	input    *Input
}

func newFinderUIFromPlayer(p *player.Player) *FinderUI {
	finder := finderFromPlayer(p)
	return &FinderUI{
		finder:   finder,
		requests: make(chan func(*FinderUI)),
		choice:   make(chan *player.Song),
		input:    newInput(),
	}
}

/*

 > ________________ 48x1 => 46x1 query
 1                  48x7
 2
 ...

*/

func (f *FinderUI) RenderQuery(query string) {
	termbox.SetCell(1, 0, '>', termbox.ColorRed, termbox.ColorDefault)
	m := -1
	unicodeCells(query, 46, false, func(x int, r rune) {
		termbox.SetCell(3+x, 0, r, termbox.AttrBold, termbox.ColorDefault)
		m = x
	})
	termbox.SetCell(3+m+1, 0, '_', 0xf0, termbox.ColorDefault)
}

func (f *FinderUI) RenderResults(cursor int, results []*item) {
	start := 0
	end := min(7, len(results))
	if cursor > 6 {
		start = cursor - 6
		end = cursor + 1
	}
	j := 0
	for i := start; i < end; i++ {
		song := f.finder.Get(results[i])
		color := termbox.ColorDefault
		if i == cursor {
			color = termbox.AttrReverse
		}
		unicodeCells(song.Name(), 48, true, func(x int, r rune) {
			termbox.SetCell(1+x, 1+j, r, color, termbox.ColorDefault)
		})
		j++
	}
}

func (f *FinderUI) HandleKeyStrokes() {
	results := f.finder.items
	cursor := 0
	exit := false
	for !exit {
		must(termbox.Clear(termbox.ColorDefault, termbox.ColorDefault))
		f.RenderQuery(f.input.String())
		f.RenderResults(cursor, results)
		must(termbox.Sync())
		ev := termbox.PollEvent()
		switch ev.Key {
		case termbox.KeyArrowUp:
			if cursor > 0 {
				cursor--
			}
		case termbox.KeyArrowDown:
			if cursor < len(results)-1 {
				cursor++
			}
		case termbox.KeyEsc:
			cursor = -1
			exit = true
		case termbox.KeyEnter:
			exit = true
		default:
			f.input.Feed(ev)
			results = f.finder.Find(f.input.String())
			cursor = 0
		}
	}
	if cursor < len(results) && cursor >= 0 {
		song := f.finder.Get(results[cursor])
		f.choice <- &song
		return
	}
	f.choice <- nil
}
