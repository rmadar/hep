package fads

import (
	"math/rand"
	"reflect"

	"github.com/go-hep/fwk"
	"github.com/go-hep/random"
)

type efficiency struct {
	fwk.TaskBase

	input  string
	output string

	eff  func(pt, eta float64) float64
	seed int64
	dist random.Dist
}

func (tsk *efficiency) Configure(ctx fwk.Context) error {
	var err error
	err = tsk.DeclInPort(tsk.input, reflect.TypeOf([]Candidate{}))
	if err != nil {
		return err
	}

	err = tsk.DeclOutPort(tsk.output, reflect.TypeOf([]Candidate{}))
	if err != nil {
		return err
	}

	return err
}

func (tsk *efficiency) StartTask(ctx fwk.Context) error {
	var err error
	src := rand.NewSource(tsk.seed)
	tsk.dist = random.Flat(0, 1, &src)
	return err
}

func (tsk *efficiency) StopTask(ctx fwk.Context) error {
	var err error

	return err
}

func (tsk *efficiency) Process(ctx fwk.Context) error {
	var err error
	store := ctx.Store()
	msg := ctx.Msg()

	v, err := store.Get(tsk.input)
	if err != nil {
		return err
	}

	input := v.([]Candidate)
	msg.Debugf(">>> input: %v\n", len(input))

	output := make([]Candidate, 0, len(input))
	defer func() {
		err = store.Put(tsk.output, output)
	}()

	for i := range input {
		cand := &input[i]
		eta := cand.Pos.Eta()
		pt := cand.Mom.Pt()

		// apply efficiency
		eff := tsk.dist()
		max := tsk.eff(pt, eta)
		if eff > max {
			continue
		}

		output = append(output, *cand)
	}

	msg.Debugf(">>> output: %v\n", len(output))

	return err
}

func init() {
	fwk.Register(reflect.TypeOf(efficiency{}),
		func(typ, name string, mgr fwk.App) (fwk.Component, error) {
			var err error
			tsk := &efficiency{
				TaskBase: fwk.NewTask(typ, name, mgr),
				input:    "InputParticles",
				output:   "OutputParticles",
				eff:      func(x, y float64) float64 { return 1 },
				seed:     1234,
			}
			err = tsk.DeclProp("Input", &tsk.input)
			if err != nil {
				return nil, err
			}

			err = tsk.DeclProp("Output", &tsk.output)
			if err != nil {
				return nil, err
			}

			err = tsk.DeclProp("Eff", &tsk.eff)
			if err != nil {
				return nil, err
			}

			err = tsk.DeclProp("Seed", &tsk.seed)
			if err != nil {
				return nil, err
			}

			return tsk, err
		},
	)
}
