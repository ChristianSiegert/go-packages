package sessions

import "encoding/json"

// Flashes manages flashes.
type Flashes interface {
	// Add adds flashes.
	Add(flashes ...Flash)

	// AddNew creates a new Flash and adds it. flashType is optional. Only the
	// first given flashType is used.
	AddNew(message string, flashType ...string) Flash

	// GetAll returns all flashes.
	GetAll() []Flash

	// Remove removes flashes.
	Remove(flashes ...Flash)

	// RemoveAll removes all flashes.
	RemoveAll()
}

type flashes []Flash

// NewFlashes returns a new instance of Flashes.
func NewFlashes() Flashes {
	return &flashes{}
}

// Add adds flashes.
func (f *flashes) Add(flashes ...Flash) {
	for _, flash := range flashes {
		*f = append(*f, flash)
	}
}

// AddNew creates a new Flash and adds it. flashType is optional. Only the
// first given flashType is used.
func (f *flashes) AddNew(message string, flashType ...string) Flash {
	flash := NewFlash(message, "")

	if len(flashType) > 0 {
		flash.SetType(flashType[0])
	}

	f.Add(flash)
	return flash
}

// GetAll returns all flashes.
func (f *flashes) GetAll() []Flash {
	return []Flash(*f)
}

// Remove removes flashes.
func (f *flashes) Remove(flashes ...Flash) {
	ff := f.GetAll()
	for _, flash := range flashes {
		for i := 0; i < len(ff); i++ {
			if ff[i] == flash {
				ff = append(ff[:i], ff[i+1:]...)
			}
		}
	}
	*f = ff
}

// RemoveAll removes all flashes.
func (f *flashes) RemoveAll() {
	*f = flashes{}
}

// FlashesFromJSON JSON decodes an array of Flash objects. The result is useful
// as input for Flashes.Add.
func FlashesFromJSON(data []byte) ([]Flash, error) {
	temp := []encodableFlash{}

	if err := json.Unmarshal(data, &temp); err != nil {
		return nil, err
	}

	ff := make([]Flash, 0, len(temp))
	for _, f := range temp {
		flash := NewFlash(f.Message, f.Type)
		ff = append(ff, flash)
	}
	return ff, nil
}
