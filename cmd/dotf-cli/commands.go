package main

type MainCommand interface {
	Run([]string)
	PrintDefaults()
}

