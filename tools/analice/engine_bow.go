package main

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"strings"
)

// codebeat:disable[TOO_MANY_IVARS]
type bowLicenseEngine struct {
	knownBags  map[license]bow
	curDep     dependency
	curBow     bow
	licsVsDeps map[license][]dependency
	rootLic    license
}

// codebeat:enable[TOO_MANY_IVARS]

func (b *bowLicenseEngine) start(p project, d dependency) error {
	b.curDep = d
	b.curBow = bow(make(map[string]int))
	return nil
}

func (b *bowLicenseEngine) stop(p project, d dependency) error {
	best, sim := b.bestMatch(b.curBow)
	if d.name != "" {
		b.licsVsDeps[best] = append(b.licsVsDeps[best], d)
	} else {
		b.rootLic = best
	}
	if sim < 0.5 {
		return fmt.Errorf("best match was '%s' but with poor similarity (%v)", best.name, sim)
	}
	return nil
}

func (b *bowLicenseEngine) analyze(line string) error {
	fields := strings.Fields(line)
	for _, field := range fields {
		b.curBow[field]++
	}
	return nil
}

func (b *bowLicenseEngine) initialize() {
	b.knownBags = make(map[license]bow)
	b.licsVsDeps = make(map[license][]dependency)
	for _, license := range licenses {
		ltm := b.bag(strings.NewReader(license.text))
		b.knownBags[license] = ltm
	}
}

func (b *bowLicenseEngine) bag(r io.Reader) bow {
	m := bow(make(map[string]int))
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

func (b *bowLicenseEngine) bestMatch(m bow) (license, float64) {
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

type bow map[string]int

func (x bow) dot(y bow) int {
	prod := 0
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
