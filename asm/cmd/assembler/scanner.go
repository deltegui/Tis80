package main

type scanner interface {
	Scan() Token
}
