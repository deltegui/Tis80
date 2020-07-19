package tisasm

type Scanner interface {
	Scan() Token
	Advance(times int)
}
