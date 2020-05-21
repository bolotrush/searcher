package models

import "context"

type Modeler interface {
	WriteIndex(map[string][]string)
	Get(query string) map[string][]string
	Close()
	Listen(ctx context.Context)
}
