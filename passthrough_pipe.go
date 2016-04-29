package pipeline

// PassThroughPipe passes the data read from input to output without modification
type PassThroughPipe struct {
}

// Process pipe input
func (p *PassThroughPipe) Process(in chan Data) chan Data {
	out := make(chan Data)
	go func() {
		defer close(out)
		for input := range in {
			out <- input
		}
	}()
	return out
}
