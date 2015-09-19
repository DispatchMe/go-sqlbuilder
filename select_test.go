package sqlbuilder_test

// This is in a different package so we can also make sure everything works from outside of the sqlbuilder package

import (
	"fmt"
	. "github.com/jraede/go-sqlbuilder"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

type param struct {
	description string
	query       *Query
	expected    string
	vars        []interface{}
}

// Basic SELECT query to get the first and last name for people younger than 10 or older than 90, ordering by last name, limiting to 20 results
func ExampleSelect_basic() {
	sql, vars := Select("first_name", "last_name").From("people").Where(
		Expr{"gender", "=", "female"},
		Or(
			Expr{"age", ">", 90},
			Expr{"age", "<", 10},
		),
	).OrderBy("last_name", ASC).Limit(20).GetSQL()

	fmt.Println(sql, ",", vars)
	// Output: SELECT first_name, last_name FROM people WHERE (gender = $1 AND (age > $2 OR age < $3)) ORDER BY last_name ASC LIMIT 20 , [female 90 10]
}

// Complex SELECT query with multiple joins. This gets the total number of rushing plays that a player was involved in that either went for more than 10 yards or went for a score, if the player has more than 5 such plays. Orders by total matching plays, and show results #51-60
func ExampleSelect_complex() {
	sql, vars := Select("COUNT(plays) AS playcount", "players.name").From("players").Join(JOIN_INNER, "play_player",
		OnColumn("players.id", "=", "play_player.player_id"),
	).Join(JOIN_INNER, "plays",
		OnColumn("plays.id", "=", "play_player.play_id"),
		OnValue("plays.type", "=", "running"),
	).Where(
		Or(
			Expr{"plays.yards", ">", 10},
			Expr{"plays.scoring", "=", true},
		),
	).GroupBy("players.id").Having(Expr{"COUNT(plays)", ">", 5}).OrderBy("playcount", DESC).Limit(10).Skip(50).GetSQL()

	fmt.Println(sql, ",", vars)
	// Output: SELECT COUNT(plays) AS playcount, players.name FROM players INNER JOIN play_player ON players.id = play_player.player_id INNER JOIN plays ON plays.id = play_player.play_id AND ON plays.type = $1 WHERE (plays.yards > $2 OR plays.scoring = $3) GROUP BY players.id HAVING COUNT(plays) > $4 ORDER BY playcount DESC LIMIT 10 SKIP 50 , [running 10 true 5]
}

func TestSelect(t *testing.T) {

	params := []param{
		{"simple", Select("*").From("foos"), "SELECT * FROM foos", nil},
		{"simple where", Select("a", "b", "c").From("foos").Where(Expr{"a", "=", 10}), "SELECT a, b, c FROM foos WHERE a = $1", []interface{}{10}},
		{"complex where", Select("*").From("stats").Where(Expr{"rushing_attempts", "<", 10}, Or(Expr{"rushing_yards", ">", 100}, Expr{"rushing_tds", ">", 0})), "SELECT * FROM stats WHERE (rushing_attempts < $1 AND (rushing_yards > $2 OR rushing_tds > $3))", []interface{}{10, 100, 0}},
		{"ordering", Select("a").From("foos").OrderBy("a.timestamp", DESC), "SELECT a FROM foos ORDER BY a.timestamp DESC", nil},
		{"multiple ordering", Select("a").From("foos").OrderBy("a.category", DESC).OrderBy("a.timestamp", ASC), "SELECT a FROM foos ORDER BY a.category DESC, a.timestamp ASC", nil},
		{"group by", Select("SUM(a.price)").From("foos").GroupBy("a.category"), "SELECT SUM(a.price) FROM foos GROUP BY a.category", nil},
		{"single join", Select("*").From("foos").Join(JOIN_INNER, "bars", OnColumn("bars.foo_id", "=", "foos.id")), "SELECT * FROM foos INNER JOIN bars ON bars.foo_id = foos.id", nil},
		{"complex single join", Select("*").From("foos").Join(JOIN_INNER, "categories", OnColumn("foos.category_id", "=", "categories.id"), OnValue("categories.type", "=", "main")), "SELECT * FROM foos INNER JOIN categories ON foos.category_id = categories.id AND ON categories.type = $1", []interface{}{"main"}},
		{"multiple joins", Select("*").From("games").Join(JOIN_INNER, "drives", OnColumn("drives.game_id", "=", "games.id")).Join(JOIN_INNER, "plays", OnColumn("plays.drive_id", "=", "drives.id")), "SELECT * FROM games INNER JOIN drives ON drives.game_id = games.id INNER JOIN plays ON plays.drive_id = drives.id", nil},
		{"everything", Select("COUNT(plays)", "players.name").From("players").Join(JOIN_INNER, "play_player", OnColumn("players.id", "=", "play_player.player_id")).Join(JOIN_INNER, "plays", OnColumn("plays.id", "=", "play_player.play_id")).Where(Or(Expr{"plays.yards", ">", 10}, Expr{"plays.scoring", "=", true})).GroupBy("players.id").Having(Expr{"COUNT(plays)", ">", 5}).OrderBy("players.name", ASC).Limit(10).Skip(50), "SELECT COUNT(plays), players.name FROM players INNER JOIN play_player ON players.id = play_player.player_id INNER JOIN plays ON plays.id = play_player.play_id WHERE (plays.yards > $1 OR plays.scoring = $2) GROUP BY players.id HAVING COUNT(plays) > $3 ORDER BY players.name ASC LIMIT 10 SKIP 50", []interface{}{10, true, 5}},
	}

	Convey("Select queries", t, func() {
		for _, p := range params {
			Convey(p.description, func() {
				sql, vars := p.query.GetSQL()
				So(sql, ShouldEqual, p.expected)

				So(len(vars), ShouldEqual, len(p.vars))
				for i, v := range vars {
					So(p.vars[i], ShouldEqual, v)
				}

			})
		}
	})
}
