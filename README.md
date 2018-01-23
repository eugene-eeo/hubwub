# orchid

Very tiny music player for my needs. Mostly an excuse for me to learn how
to use termbox after being inspired by Brandon Rhode's [talk on terminal animations](https://www.youtube.com/watch?v=rrMnmLyYjU8).
Start it up in a directory where there are MP3 files, and `orchid` will
do the rest:

    $ go get -u github.com/eugene-eeo/orchid
    $ cd totally-legit-music
    $ orchid

Ideally ran in a terminal with size 50x10 and with the excellent [Iosevka Term](https://github.com/be5invis/Iosevka)
font, although Envy Code R works as well. Note on OSX: for some reason images
are only properly displayed on screen-256 (at least for me). YMMV.

| <img src='./screenshots/demo1.png' width='300px'>   | <img src='./screenshots/demo2.png' width='300px'>  |
|:---------------------------------------:|:--------------------------------------:|
| Default art for MP3s without album art. | Album art and in Shuffle mode.         |


## controls

- `r` toggle repeat song/playlist
- `<space>` pause/play
- `s` shuffle queue
- `n` next track
- `p` prev track
- `q` quit

## todo

- [ ] `go dep`
- [ ] write tests
- [ ] source MP3s from multiple directories
- [ ] extract + display more tags (?)
