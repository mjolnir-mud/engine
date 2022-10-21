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
type UID string

// New returns a new random UID.
func New() UID {
	id, err := uuid.NewRandom()

	if err != nil {
		panic(err)
	}

	hash := sha256.New()
	hash.Write([]byte(id.String()))
	sha := hash.Sum(nil)
	uid := UID(fmt.Sprintf("%x", sha[:12]))
	var final UID
	final = uid

	return final
}

// ToBSON returns the BSON representation of the UID.
func (u UID) ToBSON() primitive.ObjectID {
	val, err := primitive.ObjectIDFromHex(string(u))

	if err != nil {
		panic(err)
	}

	return val
}

//func (u *UID) MarshalJSON() ([]byte, error) {
//	return []byte(fmt.Sprintf("\"%s\"", u.uid)), nil
//}
//
//func (u *UID) UnmarshalJSON(b []byte) error {
//	uid, err := FromString(string(b))
//
//	if err != nil {
//		return err
//	}
//
//	u.uid = uid.uid
//
//	return nil
//}
//
//func (u *UID) MarshalBSON() ([]byte, error) {
//	b := u.uid
//
//	return []byte(b), nil
//}
//
//func (u *UID) UnmarshalBSON(b []byte) error {
//	id, err := primitive.ObjectIDFromHex(string(b))
//
//	if err != nil {
//		return err
//	}
//
//	u.uid = id.Hex()
//
//	return nil
//}
//

// FromBSON returns a UID from a BSON ObjectID.
func FromBSON(id primitive.ObjectID) UID {
	var uid UID
	uid = UID(id.Hex())

	return uid
}
