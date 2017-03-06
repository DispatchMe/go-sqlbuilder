package sqlbuilder_test

import (
	. "github.com/jraede/go-sqlbuilder"
	"testing"
)

//
// Select benchmarks
//

func BenchmarkSelectBasic(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Select("first_name", "last_name").From("people").Where(
			Equal{"gender", "female"},
			Or(
				GreaterThan{"age", 90},
				LessThan{"age", 10},
			),
		).OrderBy("last_name", ASC).Limit(20).GetFullSQL()
	}
}

func BenchmarkSelectComplex(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Select("COUNT(plays) AS playcount", "players.name").From("players").InnerJoin("play_player",
			OnColumn("players.id", "play_player.player_id"),
		).InnerJoin("plays",
			OnColumn("plays.id", "play_player.play_id"),
			OnExpression(Equal{"plays.type", "running"}),
		).Where(
			Or(
				GreaterThan{"plays.yards", 10},
				Equal{"plays.scoring", true},
			),
		).GroupBy("players.id").Having(GreaterThan{"COUNT(plays)", 5}).OrderBy("playcount", DESC).Limit(10).Offset(50).GetFullSQL()
	}
}

func BenchmarkSelectWithSubqueries(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Select("foo").From(Alias(Union(Select("foo").From("table1"), Select("foo").From("table2")), "u"))
	}
}
