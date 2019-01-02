package main

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"strings"
)

type bow map[string]uint64

func (x bow) dot(y bow) uint64 {
	prod := uint64(0)
	for k, v := range x {
		prod += v * y[k]
	}
	return prod
}

func (x bow) cos(y bow) float64 {
	xy := float64(x.dot(y))
	xx := float64(x.dot(x))
	yy := float64(y.dot(y))
	return xy / (math.Sqrt(xx * yy))
}

type bowEngine struct {
	knownBags map[license]bow
	curBow    bow
	bestGuess license
}

func (b *bowEngine) start() {
	b.curBow = bow(make(map[string]uint64))
}

func (b *bowEngine) stop() error {
	best, sim := b.bestMatch(b.curBow)
	b.bestGuess = best
	if sim < 0.5 {
		return fmt.Errorf("best match was '%s' but with poor similarity (%v)", best.name, sim)
	}
	return nil
}

func (b *bowEngine) analyze(line string) error {
	fields := strings.Fields(line)
	for _, field := range fields {
		b.curBow[field]++
	}
	return nil
}

func (b *bowEngine) initialize() {
	b.knownBags = make(map[license]bow)
	for _, license := range licenses {
		ltm := b.bag(strings.NewReader(license.text))
		b.knownBags[license] = ltm
	}
}

func (b *bowEngine) bag(r io.Reader) bow {
	m := bow(make(map[string]uint64))
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		for _, field := range fields {
			m[field]++
		}
	}
	return m
}

func (b *bowEngine) bestMatch(m bow) (license, float64) {
	var bestLicense license
	bestScore := 0.0
	for l, bag := range b.knownBags {
		cos := m.cos(bag)
		if math.Abs(1.0-cos) < math.Abs(1.0-bestScore) {
			bestScore, bestLicense = cos, l
		}
	}
	return bestLicense, bestScore
}
