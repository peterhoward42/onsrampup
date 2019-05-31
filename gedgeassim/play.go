package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gedge/graphson"

	"github.com/gedge/gremgo-neptune"
	_ "github.com/schwartzmx/gremgo-neptune"
)

func main() {
	errs := make(chan error)
	go func(chan error) {
		err := <-errs
		log.Fatal("Lost connection to the database: " + err.Error())
	}(errs)

	// Schwartzmz tolerates / encourages just the ip:port url and automatically
	// adds on gremlin under the hood. Gedge requires it to be fully
	// formed in the first place.
	dialer := gremgo.NewDialer("ws://127.0.0.1:8182/gremlin")
	g, err := gremgo.Dial(dialer, errs)
	if err != nil {
		fmt.Println(err)
		return
	}

	// This is the signature for schwartzmx/gremgo-neptune
	//res, err := g.Execute(myGremlinScript)

	// And this for gedge/gremgo-neptune.
	// Neptune doesn't support bindings<>placholders in Gremlin scripts.
	// And so schwartzmx/gremgo removed them from this signature.
	// Gedge has reinstated them in the signature, but they won't / can't do
	// anything when pointing at Neptune. Hence nil, nil.
	res, err := g.Execute(myGremlinScript, nil, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	// This is the native json way to inspect results
	j, err := json.Marshal(res[0].Result.Data)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("XXXXXX: Native JSON form: %s\n", j)

	// And this the graphson-parsed way
	raw := res[0].Result.Data
	graphsonRes, err := graphson.DeserializeStringListFromBytes(raw)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("XXXXXX: Parsed grapson JSON form: %s\n", graphsonRes)
}

// here
const myGremlinScript = `

g.V().hasLabel('person').has('name', within('peter', 'chris', 'pras')).
    drop().iterate()

g.addV('person').property('name', 'peter').property('age', 60).iterate()
g.addV('person').property('name', 'chris').property('age', 30).iterate()
g.addV('person').property('name', 'pras').property('age', 25).iterate()

g.addV('sw').property('name', 'network').property('lang', 'python').iterate()
g.addV('sw').property('name', 'asc').property('lang', 'react').iterate()

g.addE('knows').
    from(V().has('person', 'name', 'peter')).
    to(V().has('person', 'name', 'chris')).iterate()

g.addE('knows').
    from(V().has('person', 'name', 'peter')).
    to(V().has('person', 'name', 'pras')).iterate()

g.addE('knows').
    from(V().has('person', 'name', 'chris')).
    to(V().has('person', 'name', 'pras')).iterate()

g.addE('knows').
    from(V().has('person', 'name', 'chris')).
    to(V().has('person', 'name', 'peter')).iterate()

g.addE('knows').
    from(V().has('person', 'name', 'pras')).
    to(V().has('person', 'name', 'peter')).iterate()

g.addE('knows').
    from(V().has('person', 'name', 'pras')).
    to(V().has('person', 'name', 'chris')).iterate()

g.addE('created').
    from(V().has('person', 'name', 'peter')).
    to(V().has('sw', 'name', 'network')).iterate()

g.addE('created').
    from(V().has('person', 'name', 'peter')).
    to(V().has('sw', 'name', 'asc')).iterate()

g.addE('created').
    from(V().has('person', 'name', 'chris')).
    to(V().has('sw', 'name', 'network')).iterate()

g.addE('created').
    from(V().has('person', 'name', 'pras')).
    to(V().has('sw', 'name', 'asc')).iterate()

g.V().has('person', 'name', 'pras').as('exclude').
    out('created').as('a').values('name')
`
