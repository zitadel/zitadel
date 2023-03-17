package runtime

import "os"

type Phase struct {
	latestSig       os.Signal
	originalSigChan chan os.Signal
}

func StartTracking(sigChan chan os.Signal) *Phase {
	p := &Phase{
		originalSigChan: sigChan,
	}
	go func() {
		p.latestSig = <-p.ForkShutdown()
	}()
	return p
}

func (p *Phase) ShuttingDown() bool {
	return p.latestSig != nil
}

func (p *Phase) ForkShutdown() chan os.Signal {
	fork := make(chan os.Signal)
	go func() {
		sig := <-p.originalSigChan
		go func() {
			p.originalSigChan <- sig
		}()
		fork <- sig
	}()
	return fork
}
