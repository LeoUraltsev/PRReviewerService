package team

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	mocks_team "github.com/LeoUraltsev/PRReviewerService/internal/http/handler/team/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandler_AddingTeam(t *testing.T) {
	m := mocks_team.NewMockSaver(t)
	m.EXPECT().Save(mock.Anything, mock.Anything).Return(nil).Once()

	r := httptest.NewRequest("POST", "/team/add", strings.NewReader(`{
  "team_name": "payments",
  "members": [
    {
      "user_id": "u1",
      "username": "Alice",
      "is_active": true
    },
    {
      "user_id": "u2",
      "username": "Bob",
      "is_active": true
    }
  ]
}`))
	w := httptest.NewRecorder()
	h := NewHandler(m)
	h.AddingTeam(w, r)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
}
