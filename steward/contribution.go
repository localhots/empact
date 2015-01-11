package steward

type (
	Contribution struct {
		Author    string
		Repo      string
		Week      int64
		Commits   int
		Additions int
		Deletions int
	}
)
