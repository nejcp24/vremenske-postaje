package vreme

import (
	"fmt"
	"math"
	"sort"
)

// Postaja predstavlja vremensko postajo z imenom, krajem in meritvami.
type Postaja struct {
	ime     string
	kraj    string
	meritve []float64
}

// DodajPostajo doda novo vremensko postajo v slovar postaj.
// Postaja je določena z oznako, imenom in krajem.
// Če oznaka že obstaja ali slovar ni inicializiran, vrne napako.
func DodajPostajo(postaje map[string]Postaja, oznaka, ime, kraj string) error {
	if postaje == nil {
		return fmt.Errorf("slovar postaj ni inicializiran")
	}

	if oznaka == "" {
		return fmt.Errorf("oznaka postaje ne sme biti prazna")
	}

	if _, obstaja := postaje[oznaka]; obstaja {
		return fmt.Errorf("postaja %s že obstaja", oznaka)
	}

	postaje[oznaka] = Postaja{
		ime:  ime,
		kraj: kraj,
	}

	return nil
}

// DodajMeritev doda temperaturno meritev obstoječi postaji.
// Meritev mora biti znotraj dovoljenega območja minTemp–maxTemp.
// Če postaja ne obstaja ali temperatura ni dovoljena, vrne napako.
func DodajMeritev(postaje map[string]Postaja, oznaka string, temperatura, minTemp, maxTemp float64) error {
	if minTemp > maxTemp || math.IsNaN(minTemp) || math.IsNaN(maxTemp) {
		return fmt.Errorf("neveljavno temperaturno območje")
	}

	if math.IsNaN(temperatura) || math.IsInf(temperatura, 0) {
		return fmt.Errorf("temperatura ni veljavna")
	}

	if temperatura < minTemp || temperatura > maxTemp {
		return fmt.Errorf("temperatura %.1f °C ni v dovoljenem območju [%.1f, %.1f] °C",
			temperatura, minTemp, maxTemp)
	}

	p, obstaja := postaje[oznaka]

	if !obstaja {
		return fmt.Errorf("postaja %s ne obstaja", oznaka)
	}

	p.meritve = append(p.meritve, temperatura)
	postaje[oznaka] = p

	return nil
}

// statistika izračuna najmanjšo, največjo in povprečno temperaturo
// za podano postajo. Če postaja ne obstaja ali nima dovolj meritev,
// vrne ok = false.
func statistika(postaje map[string]Postaja, oznaka string, stMeritev int) (najnizja, najvisja, povprecna float64, ok bool) {
	pos, obstaja := postaje[oznaka]

	if !obstaja || stMeritev < 1 || len(pos.meritve) < stMeritev {
		return 0, 0, 0, false
	}

	najnizja = pos.meritve[0]
	najvisja = pos.meritve[0]

	vsota := 0.0

	for _, vrednost := range pos.meritve {
		if vrednost < najnizja {
			najnizja = vrednost
		}

		if vrednost > najvisja {
			najvisja = vrednost
		}

		vsota += vrednost
	}

	povprecna = vsota / float64(len(pos.meritve))
	ok = true

	return
}

// IzpisMeritev izpiše vse vremenske postaje in njihove zbrane meritve.
func IzpisMeritev(postaje map[string]Postaja) {
	oznake := make([]string, 0, len(postaje))

	for oznaka := range postaje {
		oznake = append(oznake, oznaka)
	}

	sort.Strings(oznake)

	for _, oznaka := range oznake {
		vrednost := postaje[oznaka]

		fmt.Printf("%s - %s (%s): %v\n",
			oznaka, vrednost.ime, vrednost.kraj, vrednost.meritve)
	}
}

// IzpisPorocila izpiše temperaturno poročilo za vse postaje.
// Za izračun statistike mora imeti posamezna postaja vsaj stMeritev meritev.
func IzpisPorocila(postaje map[string]Postaja, stMeritev int) {
	oznake := make([]string, 0, len(postaje))

	for oznaka := range postaje {
		oznake = append(oznake, oznaka)
	}

	sort.Strings(oznake)

	for _, oznaka := range oznake {
		vrednost := postaje[oznaka]

		min, max, pov, ok := statistika(postaje, oznaka, stMeritev)

		switch {
		case !ok:
			fmt.Printf("%s: Premalo podatkov\n", vrednost.ime)

		case pov >= 25.0:
			fmt.Printf("%s: min %.1f °C, max %.1f °C, povprečje %.1f °C -> Vroče\n",
				vrednost.ime, min, max, pov)

		case pov >= 15.0:
			fmt.Printf("%s: min %.1f °C, max %.1f °C, povprečje %.1f °C -> Toplo\n",
				vrednost.ime, min, max, pov)

		case pov >= 5.0:
			fmt.Printf("%s: min %.1f °C, max %.1f °C, povprečje %.1f °C -> Hladno\n",
				vrednost.ime, min, max, pov)

		default:
			fmt.Printf("%s: min %.1f °C, max %.1f °C, povprečje %.1f °C -> Mrzlo\n",
				vrednost.ime, min, max, pov)
		}
	}
}
