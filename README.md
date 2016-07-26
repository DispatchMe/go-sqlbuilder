This is a SQL query bulder inspired by my PHP days, specifically [FuelPHP's query builder](http://fuelphp.com/docs/classes/database/usage.html). It is meant for use with PostgreSQL but should also be fine with MySQL.

Please see the [godoc](http://godoc.org/github.com/jraede/go-sqlbuilder) for the complete documentation.

## Heres a teaser:

The below is a complex `SELECT` query with multiple joins.

```go
rows, err := Select("COUNT(plays) AS playcount", "players.name").From("players").InnerJoin("play_player",
	OnColumn("players.id", "play_player.player_id"),
).InnerJoin("plays",
	OnColumn("plays.id", "play_player.play_id"),
	OnExpression(Equal{"plays.type", "running"}),
).Where(
	Or(
		GreaterThan{"plays.yards", 10},
		Equal{"plays.scoring", true},
	),
).GroupBy("players.id").
Having(GreaterThan{"COUNT(plays)", 5}).
OrderBy("playcount", DESC).
Limit(10).Offset(50).ExecRead(myDBConnection)
```
