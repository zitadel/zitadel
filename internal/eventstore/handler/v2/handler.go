package handler

import "github.com/caos/logging"

type Handler struct {
	preSteps  []PreStep
	postSteps []PostStep
}

func (h Handler) execPreSteps() error {
	for _, step := range h.preSteps {
		if err := step(); err != nil {
			return err
		}
	}
	return nil
}

func (h Handler) execPostSteps() error {
	for _, step := range h.postSteps {
		if err := step(); err != nil {
			logging.Log("V2-AtKUv").WithError(err).Warn("post step failed")
		}
	}
	return nil
}

type PreStep func() error

type PostStep func() error
