package calculon

var precedense = map[string]struct {
	prec       int
	rightAssoc bool
}{
	"+": {2, false},
	"-": {2, false},
	"*": {3, false},
	"/": {3, false},
	"%": {3, false},
	"^": {4, true},
}
