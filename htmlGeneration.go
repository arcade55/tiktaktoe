package main

import (
	"fmt"

	"github.com/arcade55/htma"
)

func homepage() htma.Element {
	doc := htma.HTML().
		LangAttr("en")

	head := htma.Head().
		AddChild(
			htma.Meta().
				CharsetAttr("UTF-8"),
		).
		AddChild(
			htma.Meta().
				NameAttr("viewport"),
		).
		AddChild(
			htma.Title("Tic-Tac-Toe"),
		).
		AddChild(
			htma.Link().
				RelAttr("stylesheet").
				HrefAttr("static/style.css"),
		).
		AddChild(
			htma.Script().
				TypeAttr("module").
				SrcAttr("https://cdn.jsdelivr.net/gh/starfederation/datastar@v1.0.0-RC.1/bundles/datastar.js"),
		)

	body := htma.Body().
		DataOnLoadAttr("@get('/sse')")

	mainElem := htma.Main()

	h1 := htma.H1().Text("Real-Time Tic-Tac-Toe")
	h2 := htma.H2().DataTextAttr("'You are player ' + $shape")
	//	pre := htma.Pre().DataJsonSignalsAttr("")

	resetButton := htma.Button().
		IDAttr("reset-button").
		DataOnClickAttr("@post('/reset')").
		Text("Reset Game")

	//mainElem = mainElem.AddChild(h1, h2, pre, board(), resetButton)
	mainElem = mainElem.AddChild(h1, h2, board(), resetButton)
	body = body.AddChild(mainElem)
	doc = doc.AddChild(head, body)

	return (doc)
}

func board() htma.Element {
	board := htma.Div().
		//DataRefAttr("foo").
		IDAttr("game-board").
		ClassAttr("tic-tac-toe-board")

	for i := 0; i < 9; i++ {
		cell := htma.Button().
			IDAttr(fmt.Sprintf("cell-%d", i)).
			ClassAttr("cell").
			DataOnClickAttr("$cell=el.id;@post('/cell')")

		board = board.AddChild(cell)
	}
	return board
}

func button(signals Signals) htma.Element {
	return htma.Button().
		IDAttr(signals.Cell).
		ClassAttr("cell").
		Text(signals.Shape)
}
