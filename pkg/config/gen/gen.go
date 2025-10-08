package main

import (
	cfg "github.com/conductorone/baton-metabase-v056/pkg/config"
	"github.com/conductorone/baton-sdk/pkg/config"
)

func main() {
	config.Generate("metabase-v056", cfg.Config)
}
