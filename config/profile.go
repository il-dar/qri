package config

import (
	"crypto/rand"
	"encoding/base64"
	"time"

	"github.com/libp2p/go-libp2p-crypto"
	"github.com/libp2p/go-libp2p-peer"
	"github.com/qri-io/doggos"
	"github.com/qri-io/jsonschema"
)

// Profile configures a qri profile
type Profile struct {
	ID       string `json:"id"`
	PrivKey  string `json:"privkey,omitempty"`
	Peername string `json:"peername"`
	// Created timestamp
	Created time.Time `json:"created"`
	// Updated timestamp
	Updated time.Time `json:"updated"`
	// specifies weather this is a user or an organization
	Type string `json:"type"`
	// user's email address
	Email string `json:"email"`
	// user name field. could be first[space]last, but not strictly enforced
	Name string `json:"name"`
	// user-filled description of self
	Description string `json:"description"`
	// url this user wants the world to click
	HomeURL string `json:"homeurl"`
	// color this user likes to use as their theme color
	Color string `json:"color"`
	// Thumb path for user's thumbnail
	Thumb string `json:"thumb"`
	// Profile photo
	Photo string `json:"photo"`
	// Poster photo for users's profile page
	Poster string `json:"poster"`
	// Twitter is a  peer's twitter handle
	Twitter string `json:"twitter"`
}

// DefaultProfile gives a new default profile configuration, generating a new random
// private key, peer.ID, and nickname
func DefaultProfile() *Profile {
	r := rand.Reader
	now := time.Now()
	p := &Profile{
		Created: now,
		Updated: now,
		Type:    "peer",
	}

	// Generate a key pair for this host. We will use it at least
	// to obtain a valid host ID.
	if priv, pub, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r); err == nil {
		if pdata, err := priv.Bytes(); err == nil {
			p.PrivKey = base64.StdEncoding.EncodeToString(pdata)
		}

		// Obtain Peer ID from public key
		if pid, err := peer.IDFromPublicKey(pub); err == nil {
			p.ID = pid.Pretty()
			p.Peername = doggos.DoggoNick(p.ID)
		}
	}

	return p
}

// Validate validates all fields of profile returning all errors found.
func (cfg Profile) Validate() error {
	schema := jsonschema.Must(`{
    "$schema": "http://json-schema.org/draft-06/schema#",
    "title": "Profile",
    "description": "Profile of a qri peer",
    "type": "object",
    "properties": {
      "id": {
        "description": "Unique identifier for a peername",
        "type": "string"
      },
      "privkey": {
        "description": "Private key associated with this peerid",
        "type": "string"
      },
      "peername": {
        "description": "Handle name for this peer on qri",
        "type": "string",
        "not": {
          "enum": [
            "me",
            "status",
            "at",
            "add",
            "history",
            "remove",
            "export",
            "profile",
            "list",
            "peers",
            "connections",
            "save",
            "connect"
          ]
        }
      },
      "created": {
        "description": "Datetime the profile was created",
        "type": "string",
        "format": "date-time"
      },
      "updated": {
        "description": "Datetime the profile was last updated",
        "type": "string",
        "format": "date-time"
      },
      "type": {
        "description": "The type of peer this profile represents",
        "type": "string",
        "enum": [
          "peer",
          "organization"
        ]
      },
      "email": {
        "description": "Email associated with this peer",
        "type": "string",
        "anyOf": [
          {
            "maxLength": 255,
            "format": "email"
          },
          {
            "maxLength": 0
          }
        ]
      },
      "name": {
        "description": "Name of peer",
        "type": "string",
        "maxLength": 255
      },
      "description": {
        "description": "Description or bio of peer",
        "type": "string",
        "maxLength": 255
      },
      "homeUrl": {
        "description": "URL associated with this peer",
        "type": "string",
        "anyOf": [
          {
            "maxLength": 255,
            "format": "uri"
          },
          {
            "maxLength": 0
          }
        ]
      },
      "color": {
        "description": "Color scheme peer prefers viewing qri on webapp",
        "type": "string",
        "anyOf": [
          {
            "enum": [
              "default"
            ]
          },
          {
            "maxLength": 0
          }
        ]
      },
      "thumb": {
        "description": "Location of thumbnail of peer's profile picture, an ipfs hash",
        "type": "string"
      },
      "photo": {
        "description": "Location of peer's profile picture, an ipfs hash",
        "type": "string"
      },
      "poster": {
        "description": "Location of a peer's profile poster, an ipfs hash",
        "type": "string"
      },
      "twitter": {
        "description": "Twitter handle associated with peer",
        "type": "string",
        "maxLength": 15
      }
    },
    "required": [
      "id",
      "created",
      "updated",
      "type",
      "peername",
      "privkey"
    ]
  }`)
	return validate(schema, &cfg)
}