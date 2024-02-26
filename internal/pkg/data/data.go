package data

// Data defines the interface for various *arr types
type Data interface {
	ID() int
	IMDBID() string
	Quality() string
	ReleaseDate() string
	ReleaseGroup() string
	Title() string
	Type() string
	URL() string
	Service() string
}
