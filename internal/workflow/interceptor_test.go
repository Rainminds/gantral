package workflow

import (
	"testing"

	"github.com/Rainminds/gantral/pkg/models"
	"github.com/stretchr/testify/assert"
)

func TestExtractArtifact(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected *models.CommitmentArtifact
		found    bool
	}{
		{
			name:  "Nil Input",
			input: nil,
			found: false,
		},
		{
			name: "Direct Pointer",
			input: func() interface{} {
				art := &models.CommitmentArtifact{ArtifactID: "art-1", AuthorityState: "APPROVED"}
				return &art
			}(),
			expected: &models.CommitmentArtifact{ArtifactID: "art-1", AuthorityState: "APPROVED"},
			found:    true,
		},
		{
			name: "Interface Wrapper - Map",
			input: func() interface{} {
				m := map[string]interface{}{
					"artifact_id":     "art-2",
					"authority_state": "REJECTED",
					"instance_id":     "inst-1",
				}
				var i interface{} = m
				return &i
			}(),
			expected: &models.CommitmentArtifact{ArtifactID: "art-2", AuthorityState: "REJECTED", InstanceID: "inst-1"},
			found:    true,
		},
		{
			name: "Interface Wrapper - Struct",
			input: func() interface{} {
				art := &models.CommitmentArtifact{ArtifactID: "art-3", AuthorityState: "APPROVED"}
				var i interface{} = art
				return &i
			}(),
			expected: &models.CommitmentArtifact{ArtifactID: "art-3", AuthorityState: "APPROVED"},
			found:    true,
		},
		{
			name: "Invalid Map",
			input: func() interface{} {
				m := map[string]interface{}{"foo": "bar"}
				var i interface{} = m
				return &i
			}(),
			found: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			art, found := extractArtifact(tt.input)
			assert.Equal(t, tt.found, found)
			if found {
				assert.Equal(t, tt.expected.ArtifactID, art.ArtifactID)
				assert.Equal(t, tt.expected.AuthorityState, art.AuthorityState)
				assert.Equal(t, tt.expected.InstanceID, art.InstanceID)
			}
		})
	}
}
