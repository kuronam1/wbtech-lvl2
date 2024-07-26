package main

import "fmt"

type DVDPlayer struct{}

func (d *DVDPlayer) On() {
	fmt.Println("DVD Player is on")
}

func (d *DVDPlayer) Off() {
	fmt.Println("DVD Player is off")
}

func (d *DVDPlayer) Play(movie string) {
	fmt.Printf("Playing movie: %s\n", movie)
}

type Projector struct{}

func (p *Projector) On() {
	fmt.Println("Projector is on")
}

func (p *Projector) Off() {
	fmt.Println("Projector is off")
}

func (p *Projector) WideScreenMode() {
	fmt.Println("Projector in widescreen mode")
}

type SoundSystem struct{}

func (s *SoundSystem) On() {
	fmt.Println("Sound System is on")
}

func (s *SoundSystem) Off() {
	fmt.Println("Sound System is off")
}

func (s *SoundSystem) SetVolume(volume int) {
	fmt.Printf("Setting volume to %d\n", volume)
}

type HomeTheaterFacade struct {
	dvdPlayer   *DVDPlayer
	projector   *Projector
	soundSystem *SoundSystem
}

func NewHomeTheaterFacade(dvdPlayer *DVDPlayer, projector *Projector, soundSystem *SoundSystem) *HomeTheaterFacade {
	return &HomeTheaterFacade{
		dvdPlayer:   dvdPlayer,
		projector:   projector,
		soundSystem: soundSystem,
	}
}

func (h *HomeTheaterFacade) WatchMovie(movie string) {
	fmt.Println("Get ready to watch a movie...")
	h.projector.On()
	h.projector.WideScreenMode()
	h.soundSystem.On()
	h.soundSystem.SetVolume(10)
	h.dvdPlayer.On()
	h.dvdPlayer.Play(movie)
}

func (h *HomeTheaterFacade) StopMovie() {
	fmt.Println("Shutting movie theater down...")
	h.dvdPlayer.Off()
	h.projector.Off()
	h.soundSystem.Off()
}
