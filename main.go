package main

type song struct {
  Title string `json:"title"`
}

var songs  = []song{
  {Title: "Chiki"},
  {Title: "Poin"}
}
