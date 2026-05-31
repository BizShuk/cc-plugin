package cmd

import (
	"testing"

	"github.com/bizshuk/cc-plugin/model"
)

func TestQualifiesForTruth(t *testing.T) {
	c1 := model.Candidate{Kind: "inference", ConfirmedByHuman: true}
	if QualifiesForTruth(c1, 5) {
		t.Error("inference should never qualify")
	}

	c2 := model.Candidate{Kind: "preference", ConfirmedByHuman: true}
	if !QualifiesForTruth(c2, 1) {
		t.Error("human confirmed should qualify")
	}

	c3 := model.Candidate{Kind: "fact", FirstPerson: true}
	if !QualifiesForTruth(c3, 1) {
		t.Error("first person fact should qualify")
	}

	c4 := model.Candidate{Kind: "fact", FirstPerson: false}
	if QualifiesForTruth(c4, 1) {
		t.Error("unconfirmed single source fact should not qualify")
	}
	if !QualifiesForTruth(c4, 2) {
		t.Error("corroborated fact should qualify")
	}
}
