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

package uid

import (
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

func TestNew(t *testing.T) {
	uid := New()

	assert.Len(t, uid.String(), 24)
}

func TestFromString(t *testing.T) {
	uid, err := FromString("123456789012345678901234")

	assert.Nil(t, err)
	assert.Equal(t, "123456789012345678901234", uid.String())
}

func TestFromString_Invalid(t *testing.T) {
	_, err := FromString("12345678901234567890123")

	assert.NotNil(t, err)
}

func TestFromBSON(t *testing.T) {
	bson, err := primitive.ObjectIDFromHex("123456789012345678901234")
	assert.Nil(t, err)

	uid := FromBSON(bson)

	assert.Equal(t, "123456789012345678901234", uid.String())
}

func TestUID_BSON(t *testing.T) {
	uid := New()
	bson := uid.BSON()

	assert.Len(t, bson.Hex(), 24)
}

func TestUID_String(t *testing.T) {
	uid := New()
	assert.Len(t, uid.String(), 24)
}
