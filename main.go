package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/nejcp24/vremenske-postaje/vreme"
	"github.com/urfave/cli/v3"
)

func main() {

	app := &cli.Command{
		Name:  "vremenske-postaje",
		Usage: "Program za upravljanje vremenskih postaj in njihovih meritev",

		Flags: []cli.Flag{
			&cli.Float64Flag{
				Name:  "minTemp",
				Value: -50,
				Usage: "Najnižja dovoljena temperatura meritve",
			},

			&cli.Float64Flag{
				Name:  "maxTemp",
				Value: 50,
				Usage: "Najvišja dovoljena temperatura meritve",
			},

			&cli.IntFlag{
				Name:  "stMeritev",
				Value: 3,
				Usage: "Najmanjše število meritev za izračun statistike",
			},
		},

		Action: func(_ context.Context, cmd *cli.Command) error {

			minTemp := cmd.Float64("minTemp")
			maxTemp := cmd.Float64("maxTemp")
			stMeritev := cmd.Int("stMeritev")

			if minTemp > maxTemp {
				return fmt.Errorf("minTemp ne sme biti večji od maxTemp")
			}

			if stMeritev < 1 {
				return fmt.Errorf("stMeritev mora biti vsaj 1")
			}

			// Ustvarimo slovar vremenskih postaj.
			postaje := make(map[string]vreme.Postaja)

			// Dodamo začetne vremenske postaje.
			zacetnePostaje := []struct {
				oznaka string
				ime    string
				kraj   string
			}{
				{"KP01", "Koper Luka", "Koper"},
				{"KR01", "Kredarica", "Triglav"},
				{"LJ01", "Ljubljana Bežigrad", "Ljubljana"},
				{"MB01", "Maribor Tabor", "Maribor"},
			}

			for _, p := range zacetnePostaje {
				if err := vreme.DodajPostajo(postaje, p.oznaka, p.ime, p.kraj); err != nil {
					return err
				}
			}

			// Dodamo začetne temperaturne meritve.
			zacetneMeritve := []struct {
				oznaka      string
				temperatura float64
			}{
				{"KP01", 24.5},
				{"KP01", 27},
				{"KP01", 26.1},
				{"KP01", 28.4},

				{"KR01", 2},
				{"KR01", -1.5},
				{"KR01", 3.5},

				{"LJ01", 18.2},
				{"LJ01", 21.5},
				{"LJ01", 19.8},
				{"LJ01", 16.5},
				{"LJ01", 17.4},

				{"MB01", 12.4},
			}

			for _, m := range zacetneMeritve {
				if err := vreme.DodajMeritev(
					postaje,
					m.oznaka,
					m.temperatura,
					minTemp,
					maxTemp,
				); err != nil {
					fmt.Fprintln(os.Stderr, "Opozorilo:", err)
				}
			}

			// Izpišemo vse zbrane meritve.
			fmt.Println("Meritve vremenskih postaj:")
			vreme.IzpisMeritev(postaje)

			// Izpišemo temperaturno poročilo.
			fmt.Println("\nTemperaturno poročilo:")
			vreme.IzpisPorocila(postaje, stMeritev)

			return nil
		},
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
