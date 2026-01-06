Run with: go run .

Compares the speed of encoding in TOON format compared to encoding to JSON.

We compare:
* github.com/alpkeskin/gotoon
* github.com/toon-format/toon-go
* encoding/json in standard library

The results:
* github.com/toon-format/toon-go is much faster than github.com/alpkeskin/gotoon
* github.com/toon-format/toon-go is only slightly slower than pretty-printe encoding/json
* encoding/json is much faster when encoding in compact format

Commentary: the readibility of TOON format justify using it in most cases instead of JSON.

