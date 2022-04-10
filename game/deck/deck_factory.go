package deck

type DeckFactory interface {
	BuildRiverDeck() Deck
	BuildDeck() Deck
}
