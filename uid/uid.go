/*
 * Copyright (c) 2022 eightfivefour llc. All rights reserved.
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated
 * documentation files (the "Software"), to deal in the Software without restriction, including without limitation
 * the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to
 * permit persons to whom the Software is furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all copies or substantial portions of the
 * Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE
 * WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
 * COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR
 * OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */

// Package uid is used for creating a mjolnir UID within the game engine.
package uid

import (
	"crypto/sha256"
	"fmt"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// A UID is a unique identifier for a mjolnir entity. UIDs are designed to be compatible with Mongo ObjectIDs so that
// they can be used interchangeably.
type UID struct {
	// The UID string.
	UID string `json:"uid"`
}

// New returns a new random UID.
func New() *UID {
	id, err := uuid.NewRandom()

	if err != nil {
		panic(err)
	}

	hash := sha256.New()
	hash.Write([]byte(id.String()))
	sha := hash.Sum(nil)

	final := sha[:12]

	return &UID{
		UID: fmt.Sprintf("%x", final),
	}
}

// String returns the string representation of the UID.
func (u *UID) String() string {
	return u.UID
}

// BSON returns the BSON representation of the UID.
func (u *UID) BSON() primitive.ObjectID {
	val, err := primitive.ObjectIDFromHex(u.UID)

	if err != nil {
		panic(err)
	}

	return val
}

// FromString returns a UID from a string.
func FromString(s string) (*UID, error) {
	err := validateUID(s)

	if err != nil {
		return nil, err
	}

	return &UID{
		UID: s,
	}, nil
}

// FromBSON returns a UID from a BSON ObjectID.
func FromBSON(id primitive.ObjectID) *UID {
	return &UID{
		UID: id.Hex(),
	}
}

func validateUID(uid string) error {
	_, err := primitive.ObjectIDFromHex(uid)
	return err
}
