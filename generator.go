package main

type Generator interface {
	Generate(tables map[string]tableInfo) error
}
