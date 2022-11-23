package main

import (
	"reflect"
	"testing"

	"iptv-playlist/internal/store"
)

func TestPlaylist_RecordsByCategories(t *testing.T) {
	tests := []struct {
		name string
		pl   Playlist
		want CatMap
	}{
		{"empty", Playlist{}, CatMap{}},
		{"ok", Playlist{Host: "some-host-name", Records: []store.Record{
			{0, "cat-0", 0, 3, "канал-3", "chan-3", "udp://::@"},
			{5, "cat-1", 0, 0, "канал-4", "chan-4", "udp://::@"},
			{0, "cat-0", 6, 2, "канал-2", "chan-2", "udp://::@"},
			{5, "cat-1", 2, 1, "канал-5", "chan-5", "udp://::@"},
			{0, "cat-0", 5, 1, "канал-1", "chan-1", "udp://::@"},
			{0, "cat-2", 2, 2, "канал-6", "chan-6", "udp://::@"},
			{1, "cat-3", 0, 3, "канал-7", "chan-7", "udp://::@"},
			{0, "cat-0", 0, 0, "канал-0", "chan-0", "udp://::@"},
		}}, CatMap{
			"cat-0": []store.Record{
				{0, "cat-0", 0, 0, "канал-0", "chan-0", "udp://::@"},
				{0, "cat-0", 0, 3, "канал-3", "chan-3", "udp://::@"},
				{0, "cat-0", 5, 1, "канал-1", "chan-1", "udp://::@"},
				{0, "cat-0", 6, 2, "канал-2", "chan-2", "udp://::@"},
			},
			"cat-1": []store.Record{
				{5, "cat-1", 0, 0, "канал-4", "chan-4", "udp://::@"},
				{5, "cat-1", 2, 1, "канал-5", "chan-5", "udp://::@"},
			},
			"cat-2": []store.Record{
				{0, "cat-2", 2, 2, "канал-6", "chan-6", "udp://::@"},
			},
			"cat-3": []store.Record{
				{1, "cat-3", 0, 3, "канал-7", "chan-7", "udp://::@"},
			},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.pl.RecordsByCategories(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RecordsByCategories() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSortInCat(t *testing.T) {
	tests := []struct {
		name string
		recs []store.Record
		want []store.Record
	}{
		{"empty", []store.Record{}, []store.Record{}},
		{"ok",
			[]store.Record{
				{0, "cat-0", 0, 3, "канал-3", "chan-3", "udp://::@"},
				{5, "cat-1", 0, 0, "канал-4", "chan-4", "udp://::@"},
				{0, "cat-0", 6, 2, "канал-2", "chan-2", "udp://::@"},
				{5, "cat-1", 2, 1, "канал-5", "chan-5", "udp://::@"},
				{0, "cat-0", 5, 1, "канал-1", "chan-1", "udp://::@"},
				{0, "cat-2", 2, 2, "канал-6", "chan-6", "udp://::@"},
				{1, "cat-3", 0, 3, "канал-7", "chan-7", "udp://::@"},
				{0, "cat-0", 0, 0, "канал-0", "chan-0", "udp://::@"},
			},
			[]store.Record{
				{0, "cat-0", 0, 0, "канал-0", "chan-0", "udp://::@"},
				{0, "cat-0", 0, 3, "канал-3", "chan-3", "udp://::@"},
				{0, "cat-2", 2, 2, "канал-6", "chan-6", "udp://::@"},
				{0, "cat-0", 5, 1, "канал-1", "chan-1", "udp://::@"},
				{0, "cat-0", 6, 2, "канал-2", "chan-2", "udp://::@"},
				{1, "cat-3", 0, 3, "канал-7", "chan-7", "udp://::@"},
				{5, "cat-1", 0, 0, "канал-4", "chan-4", "udp://::@"},
				{5, "cat-1", 2, 1, "канал-5", "chan-5", "udp://::@"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SortInCat(tt.recs); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SortInCat() = %v, want %v", got, tt.want)
			}
		})
	}
}
