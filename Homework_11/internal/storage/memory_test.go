package storage

import (
	"testing"

	"go_course/Homework_5/internal/model"

	"github.com/stretchr/testify/require"
)

func TestMemoryStore_CreateAndGetByUID(t *testing.T) {
	s := NewMemoryStore()

	err := s.Create(model.Task{UID: "u1", Title: "T1", IsDone: false})
	require.NoError(t, err)
	err = s.Create(model.Task{UID: "u2", Title: "T2", IsDone: true})
	require.NoError(t, err)
	err = s.Create(model.Task{UID: "u1", Title: "T3", IsDone: true})
	require.NoError(t, err)

	gotU1, err := s.GetByUID("u1")
	require.NoError(t, err)
	require.Len(t, gotU1, 2)
	require.Equal(t, "u1", gotU1[0].UID)
	require.Equal(t, "u1", gotU1[1].UID)

	gotU2, err := s.GetByUID("u2")
	require.NoError(t, err)
	require.Len(t, gotU2, 1)
	require.Equal(t, "T2", gotU2[0].Title)

	gotEmpty, err := s.GetByUID("nope")
	require.NoError(t, err)
	require.Len(t, gotEmpty, 0)
}
