package main

import "context"

func main() {
	ProvideRunner().Run(
		context.Background(),
		ProvideServer(),
	)
}
