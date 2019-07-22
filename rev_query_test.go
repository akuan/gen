package main

import (
	"testing"
)

func TestLikeQueryColums(t *testing.T) {
	v, e := LikeQueryColums("postgres", "user=postgres host=localhost port=5432  password=share dbname=labwx sslmode=disable", "college")
	t.Logf("LikeQueryColums is %d,%v", len(v), e)
}
