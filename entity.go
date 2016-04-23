package dirt

type Entity struct {
	*Location
}

type Location struct {
	Entities []Entity
}
