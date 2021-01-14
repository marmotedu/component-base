// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// Package runtime defines some functions used to encode/decode object.
package runtime

// Encoder writes objects to a serialized form.
type Encoder interface {
	// Encode writes an object to a stream. Implementations may return errors if the versions are
	// incompatible, or if no conversion is defined.
	Encode(v interface{}) ([]byte, error)
}

// Decoder attempts to load an object from data.
type Decoder interface {
	Decode(data []byte, v interface{}) error
}

// ClientNegotiator handles turning an HTTP content type into the appropriate encoder.
// Use NewClientNegotiator or NewVersionedClientNegotiator to create this interface from
// a NegotiatedSerializer.
type ClientNegotiator interface {
	Encoder() (Encoder, error)
	Decoder() (Decoder, error)
}
