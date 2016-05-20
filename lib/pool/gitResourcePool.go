package pool

// Pool holds gitResources.
type gitResourcePool struct {
	pool chan *gitResource
}

// NewGitResourcePool creates a new pool of Clients.
func NewGitResourcePool(max int) *gitResourcePool {
	return &gitResourcePool{
		pool: make(chan *gitResource, max),
	}
}

// Borrow a GitResource from the pool.
func (p *gitResourcePool) Borrow() *gitResource {
	var c *gitResource
	select {
	case c = <-p.pool:
	default:
		c = &gitResource{}
	}
	return c
}

// Return returns a GitResource to the pool.
func (p *gitResourcePool) Return(c *gitResource) {
	select {
	case p.pool <- c:
	default:
		// let it go, let it go...
	}
}
