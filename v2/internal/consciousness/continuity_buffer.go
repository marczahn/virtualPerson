package consciousness

func NewContinuityBuffer(capacity int) *ContinuityBuffer {
	if capacity < 0 {
		capacity = 0
	}
	return &ContinuityBuffer{
		capacity: capacity,
		thoughts: make([]Thought, 0, capacity),
	}
}

func (b *ContinuityBuffer) Add(thought Thought) {
	if b == nil || b.capacity == 0 {
		return
	}
	b.thoughts = append(b.thoughts, thought)
	if len(b.thoughts) > b.capacity {
		b.thoughts = b.thoughts[len(b.thoughts)-b.capacity:]
	}
}

func (b *ContinuityBuffer) Items() []Thought {
	if b == nil || len(b.thoughts) == 0 {
		return nil
	}
	out := make([]Thought, len(b.thoughts))
	copy(out, b.thoughts)
	return out
}
